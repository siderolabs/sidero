// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	talosnet "github.com/talos-systems/net"
	taloscluster "github.com/talos-systems/talos/pkg/cluster"
	"github.com/talos-systems/talos/pkg/cluster/check"
	clientconfig "github.com/talos-systems/talos/pkg/machinery/client/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/bundle"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
	"github.com/talos-systems/talos/pkg/provision"
	"github.com/talos-systems/talos/pkg/provision/access"
	"github.com/talos-systems/talos/pkg/provision/providers/qemu"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	"github.com/talos-systems/sidero/sfyra/pkg/constants"
)

// Cluster sets up initial Talos cluster.
type Cluster struct {
	options Options

	provisioner provision.Provisioner
	cluster     provision.Cluster
	access      *access.Adapter

	bridgeIP net.IP
	masterIP net.IP

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

	MemMB  int64
	CPUs   int64
	DiskGB int64
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

	cluster.configPath, err = clientconfig.GetDefaultPath()
	if err != nil {
		return err
	}

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

	if err = check.Wait(checkCtx, cluster.access, check.DefaultClusterChecks(), check.StderrReporter()); err != nil {
		return err
	}

	return cluster.untaint(ctx)
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

	_, cidr, err := net.ParseCIDR(cluster.options.CIDR)
	if err != nil {
		return err
	}

	cluster.bridgeIP, err = talosnet.NthIPInNetwork(cidr, 1)
	if err != nil {
		return err
	}

	cluster.masterIP, err = talosnet.NthIPInNetwork(cidr, 2)
	if err != nil {
		return err
	}

	cluster.access = access.NewAdapter(cluster.cluster, provision.WithTalosConfig(config))

	return nil
}

func (cluster *Cluster) create(ctx context.Context) error {
	_, cidr, err := net.ParseCIDR(cluster.options.CIDR)
	if err != nil {
		return err
	}

	cluster.bridgeIP, err = talosnet.NthIPInNetwork(cidr, 1)
	if err != nil {
		return err
	}

	cluster.masterIP, err = talosnet.NthIPInNetwork(cidr, 2)
	if err != nil {
		return err
	}

	request := provision.ClusterRequest{
		Name: cluster.options.Name,

		Network: provision.NetworkRequest{
			Name:         cluster.options.Name,
			CIDRs:        []net.IPNet{*cidr},
			GatewayAddrs: []net.IP{cluster.bridgeIP},
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

	masterEndpoint := cluster.masterIP.String()

	configBundle, err := bundle.NewConfigBundle(bundle.WithInputOptions(
		&bundle.InputOptions{
			ClusterName: cluster.options.Name,
			Endpoint:    fmt.Sprintf("https://%s:6443", defaultInternalLB),
			GenOptions: append(
				genOptions,
				generate.WithEndpointList([]string{masterEndpoint}),
				generate.WithInstallImage(cluster.options.InstallerImage),
				generate.WithDNSDomain("cluster.local"),
			),
		}))
	if err != nil {
		return err
	}

	request.Nodes = append(request.Nodes,
		provision.NodeRequest{
			Name:     constants.BootstrapMaster,
			Type:     machine.TypeControlPlane,
			IPs:      []net.IP{cluster.masterIP},
			Memory:   cluster.options.MemMB * 1024 * 1024,
			NanoCPUs: cluster.options.CPUs * 1000 * 1000 * 1000,
			Disks: []*provision.Disk{
				{
					Size: uint64(cluster.options.DiskGB) * 1024 * 1024 * 1024,
				},
			},
			Config: configBundle.ControlPlane(),
		})

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

func (cluster *Cluster) untaint(ctx context.Context) error {
	clientset, err := cluster.access.K8sClient(ctx)
	if err != nil {
		return err
	}

	n, err := clientset.CoreV1().Nodes().Get(ctx, constants.BootstrapMaster, metav1.GetOptions{})
	if err != nil {
		return err
	}

	oldData, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal unmodified node %q into JSON: %w", n.Name, err)
	}

	n.Spec.Taints = []corev1.Taint{}

	newData, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal modified node %q into JSON: %w", n.Name, err)
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, corev1.Node{})
	if err != nil {
		return fmt.Errorf("failed to create two way merge patch: %w", err)
	}

	if _, err := clientset.CoreV1().Nodes().Patch(ctx, n.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{}); err != nil {
		return fmt.Errorf("error patching node %q: %w", n.Name, err)
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
func (cluster *Cluster) SideroComponentsIP() net.IP {
	return cluster.masterIP
}

// BridgeIP returns the IP of the gateway (bridge).
func (cluster *Cluster) BridgeIP() net.IP {
	return cluster.bridgeIP
}

// Name returns cluster name.
func (cluster *Cluster) Name() string {
	return cluster.cluster.Info().ClusterName
}
