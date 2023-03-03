// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bootstrap

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"time"

	talosnet "github.com/siderolabs/net"
	taloscluster "github.com/siderolabs/talos/pkg/cluster"
	"github.com/siderolabs/talos/pkg/cluster/check"
	clientconfig "github.com/siderolabs/talos/pkg/machinery/client/config"
	"github.com/siderolabs/talos/pkg/machinery/config/bundle"
	"github.com/siderolabs/talos/pkg/machinery/config/generate"
	"github.com/siderolabs/talos/pkg/machinery/config/machine"
	"github.com/siderolabs/talos/pkg/provision"
	"github.com/siderolabs/talos/pkg/provision/access"
	"github.com/siderolabs/talos/pkg/provision/providers/qemu"

	"github.com/siderolabs/sidero/sfyra/pkg/constants"
)

// Cluster sets up initial Talos cluster.
type Cluster struct {
	options Options

	provisioner provision.Provisioner
	cluster     provision.Cluster
	access      *access.Adapter

	bridgeIP       netip.Addr
	controlplaneIP netip.Addr
	workerIP       netip.Addr

	stateDir   string
	cniDir     string
	configPath string
}

// Options for the bootstrap cluster.
type Options struct {
	Name string
	CIDR string

	Vmlinuz, Initramfs string
	InstallerImage     string
	CNIBundleURL       string

	TalosctlPath string

	RegistryMirrors []string

	BootstrapMemMB  int64
	BootstrapCPUs   int64
	BootstrapDiskGB int64

	VMNodes int

	VMMemMB  int64
	VMCPUs   int64
	VMDiskGB int64

	VMDefaultBootOrder string
}

// NewCluster creates new bootstrap Talos cluster.
func NewCluster(ctx context.Context, options Options) (*Cluster, error) {
	cluster := &Cluster{
		options: options,
	}

	var err error
	cluster.provisioner, err = qemu.NewProvisioner(ctx)

	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// Setup the bootstrap cluster.
func (cluster *Cluster) Setup(ctx context.Context) error {
	var err error

	defaultStateDir, err := clientconfig.GetTalosDirectory()
	if err != nil {
		return err
	}

	cluster.stateDir = filepath.Join(defaultStateDir, "clusters")
	cluster.cniDir = filepath.Join(defaultStateDir, "cni")

	fmt.Printf("bootstrap cluster state directory: %s, name: %s\n", cluster.stateDir, cluster.options.Name)

	if err = cluster.findExisting(ctx); err != nil {
		fmt.Printf("bootstrap cluster not found: %s, creating new one\n", err)

		err = cluster.create(ctx)
	}

	if err != nil {
		return err
	}

	checkCtx, checkCtxCancel := context.WithTimeout(ctx, 10*time.Minute)
	defer checkCtxCancel()

	return check.Wait(checkCtx, cluster.access, check.DefaultClusterChecks(), check.StderrReporter())
}

func (cluster *Cluster) findExisting(ctx context.Context) error {
	var err error

	cluster.cluster, err = cluster.provisioner.Reflect(ctx, cluster.options.Name, cluster.stateDir)
	if err != nil {
		return err
	}

	config, err := clientconfig.Open(cluster.configPath)
	if err != nil {
		return err
	}

	cidr, err := netip.ParsePrefix(cluster.options.CIDR)
	if err != nil {
		return err
	}

	cluster.bridgeIP, err = talosnet.NthIPInNetwork(cidr, 1)
	if err != nil {
		return err
	}

	cluster.controlplaneIP, err = talosnet.NthIPInNetwork(cidr, 2)
	if err != nil {
		return err
	}

	cluster.workerIP, err = talosnet.NthIPInNetwork(cidr, 3)
	if err != nil {
		return err
	}

	cluster.access = access.NewAdapter(cluster.cluster, provision.WithTalosConfig(config))

	return nil
}

func (cluster *Cluster) create(ctx context.Context) error {
	cidr, err := netip.ParsePrefix(cluster.options.CIDR)
	if err != nil {
		return err
	}

	cluster.bridgeIP, err = talosnet.NthIPInNetwork(cidr, 1)
	if err != nil {
		return err
	}

	cluster.controlplaneIP, err = talosnet.NthIPInNetwork(cidr, 2)
	if err != nil {
		return err
	}

	cluster.workerIP, err = talosnet.NthIPInNetwork(cidr, 3)
	if err != nil {
		return err
	}

	request := provision.ClusterRequest{
		Name: cluster.options.Name,

		Network: provision.NetworkRequest{
			Name:         cluster.options.Name,
			CIDRs:        []netip.Prefix{cidr},
			GatewayAddrs: []netip.Addr{cluster.bridgeIP},
			MTU:          constants.MTU,
			Nameservers:  constants.Nameservers,
			CNI: provision.CNIConfig{
				BinPath:  []string{filepath.Join(cluster.cniDir, "bin")},
				ConfDir:  filepath.Join(cluster.cniDir, "conf.d"),
				CacheDir: filepath.Join(cluster.cniDir, "cache"),

				BundleURL: cluster.options.CNIBundleURL,
			},
		},

		KernelPath:    cluster.options.Vmlinuz,
		InitramfsPath: cluster.options.Initramfs,

		SelfExecutable: cluster.options.TalosctlPath,
		StateDirectory: cluster.stateDir,
	}

	defaultInternalLB, _ := cluster.provisioner.GetLoadBalancers(request.Network)

	genOptions := cluster.provisioner.GenOptions(request.Network)

	for _, registryMirror := range cluster.options.RegistryMirrors {
		parts := strings.SplitN(registryMirror, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("unexpected registry mirror format: %q", registryMirror)
		}

		genOptions = append(genOptions, generate.WithRegistryMirror(parts[0], parts[1]))
	}

	controlplaneEndpoint := cluster.controlplaneIP.String()

	configBundle, err := bundle.NewBundle(bundle.WithInputOptions(
		&bundle.InputOptions{
			ClusterName: cluster.options.Name,
			Endpoint:    fmt.Sprintf("https://%s", net.JoinHostPort(defaultInternalLB, "6443")),
			GenOptions: append(
				genOptions,
				generate.WithEndpointList([]string{controlplaneEndpoint}),
				generate.WithInstallImage(cluster.options.InstallerImage),
				generate.WithDNSDomain("cluster.local"),
			),
		}))
	if err != nil {
		return err
	}

	request.Nodes = append(request.Nodes,
		provision.NodeRequest{
			Name:     constants.BootstrapControlPlane,
			Type:     machine.TypeControlPlane,
			IPs:      []netip.Addr{cluster.controlplaneIP},
			Memory:   cluster.options.BootstrapMemMB * 1024 * 1024,
			NanoCPUs: cluster.options.BootstrapCPUs * 1000 * 1000 * 1000,
			Disks: []*provision.Disk{
				{
					Size: uint64(cluster.options.BootstrapDiskGB) * 1024 * 1024 * 1024,
				},
			},
			Config: configBundle.ControlPlane(),
		},
		provision.NodeRequest{
			Name:     constants.BootstrapWorker,
			Type:     machine.TypeWorker,
			IPs:      []netip.Addr{cluster.workerIP},
			Memory:   cluster.options.BootstrapMemMB * 1024 * 1024,
			NanoCPUs: cluster.options.BootstrapCPUs * 1000 * 1000 * 1000,
			Disks: []*provision.Disk{
				{
					Size: uint64(cluster.options.BootstrapDiskGB) * 1024 * 1024 * 1024,
				},
			},
			Config: configBundle.Worker(),
		},
	)

	vmIPs := make([]netip.Addr, cluster.options.VMNodes)

	for i := range vmIPs {
		vmIPs[i], err = talosnet.NthIPInNetwork(cidr, i+4)
		if err != nil {
			return err
		}
	}

	for i := 0; i < cluster.options.VMNodes; i++ {
		request.Nodes = append(request.Nodes,
			provision.NodeRequest{
				Name:     fmt.Sprintf("pxe-%d", i),
				Type:     machine.TypeUnknown,
				IPs:      []netip.Addr{vmIPs[i]},
				Memory:   cluster.options.VMMemMB * 1024 * 1024,
				NanoCPUs: cluster.options.VMCPUs * 1000 * 1000 * 1000,
				Disks: []*provision.Disk{
					{
						Size: uint64(cluster.options.VMDiskGB) * 1024 * 1024 * 1024,
					},
				},
				PXEBooted: true,
				// TFTPServer:          set.options.BootSource.String(),
				// IPXEBootFilename:    "undionly.kpxe",
				SkipInjectingConfig: true,
				DefaultBootOrder:    cluster.options.VMDefaultBootOrder,
			})
	}

	cluster.cluster, err = cluster.provisioner.Create(ctx, request,
		provision.WithBootlader(true),
		// TODO: UEFI doesn't work correctly on PXE timeout, as it drops to UEFI shell
		// provision.WithUEFI(true),
		provision.WithTalosConfig(configBundle.TalosConfig()),
	)
	if err != nil {
		return err
	}

	cluster.access = access.NewAdapter(cluster.cluster, provision.WithTalosConfig(configBundle.TalosConfig()))

	c, err := clientconfig.Open(cluster.configPath)
	if err != nil {
		return err
	}

	c.Merge(configBundle.TalosConfig())

	if err = c.Save(cluster.configPath); err != nil {
		return err
	}

	if err = cluster.access.Bootstrap(ctx, os.Stderr); err != nil {
		return err
	}

	return nil
}

// TearDown the bootstrap cluster.
func (cluster *Cluster) TearDown(ctx context.Context) error {
	if cluster.cluster != nil {
		if err := cluster.provisioner.Destroy(ctx, cluster.cluster); err != nil {
			return err
		}

		cluster.cluster = nil
	}

	return nil
}

// KubernetesClient returns k8s client access adapter.
func (cluster *Cluster) KubernetesClient() taloscluster.K8sProvider {
	return &cluster.access.KubernetesClient
}

// SideroComponentsIP returns the IP of the master node.
func (cluster *Cluster) SideroComponentsIP() netip.Addr {
	return cluster.workerIP
}

// BridgeIP returns the IP of the gateway (bridge).
func (cluster *Cluster) BridgeIP() netip.Addr {
	return cluster.bridgeIP
}

// Name returns cluster name.
func (cluster *Cluster) Name() string {
	return cluster.cluster.Info().ClusterName
}

// Nodes return information about PXE VMs.
func (cluster *Cluster) Nodes() []provision.NodeInfo {
	return cluster.cluster.Info().ExtraNodes
}
