// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package vm

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	talosnet "github.com/talos-systems/net"
	clientconfig "github.com/talos-systems/talos/pkg/machinery/client/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
	"github.com/talos-systems/talos/pkg/provision"
	"github.com/talos-systems/talos/pkg/provision/providers/qemu"

	"github.com/talos-systems/sidero/sfyra/pkg/constants"
)

// Set is a number of PXE-booted VMs.
type Set struct {
	provisioner provision.Provisioner
	cluster     provision.Cluster
	options     Options
	stateDir    string
	cniDir      string
	bridgeIP    net.IP
}

// Options configure new VM set.
type Options struct {
	Name       string
	Nodes      int
	BootSource net.IP
	CIDR       string

	TalosctlPath string
	CNIBundleURL string

	MemMB  int64
	CPUs   int64
	DiskGB int64

	DefaultBootOrder string
}

// NewSet creates new VM set.
func NewSet(ctx context.Context, options Options) (*Set, error) {
	set := &Set{
		options: options,
	}

	var err error
	set.provisioner, err = qemu.NewProvisioner(ctx)

	if err != nil {
		return nil, err
	}

	return set, nil
}

// Setup the VM set.
func (set *Set) Setup(ctx context.Context) error {
	var err error

	defaultStateDir, err := clientconfig.GetTalosDirectory()
	if err != nil {
		return err
	}

	set.stateDir = filepath.Join(defaultStateDir, "clusters")
	set.cniDir = filepath.Join(defaultStateDir, "cni")

	fmt.Printf("VM set state directory: %s, name: %s\n", set.stateDir, set.options.Name)

	if err = set.findExisting(ctx); err != nil {
		fmt.Printf("VM set not found: %s, creating new one\n", err)

		return set.create(ctx)
	}

	return nil
}

func (set *Set) findExisting(ctx context.Context) error {
	var err error

	set.cluster, err = set.provisioner.Reflect(ctx, set.options.Name, set.stateDir)
	if err != nil {
		return err
	}

	_, cidr, err := net.ParseCIDR(set.options.CIDR)
	if err != nil {
		return err
	}

	set.bridgeIP, err = talosnet.NthIPInNetwork(cidr, 1)
	if err != nil {
		return err
	}

	return nil
}

func (set *Set) create(ctx context.Context) error {
	_, cidr, err := net.ParseCIDR(set.options.CIDR)
	if err != nil {
		return err
	}

	set.bridgeIP, err = talosnet.NthIPInNetwork(cidr, 1)
	if err != nil {
		return err
	}

	ips := make([]net.IP, 1+set.options.Nodes)

	for i := range ips {
		ips[i], err = talosnet.NthIPInNetwork(cidr, i+2)
		if err != nil {
			return err
		}
	}

	request := provision.ClusterRequest{
		Name: set.options.Name,

		Network: provision.NetworkRequest{
			Name:         set.options.Name,
			CIDRs:        []net.IPNet{*cidr},
			GatewayAddrs: []net.IP{set.bridgeIP},
			MTU:          constants.MTU,
			Nameservers:  constants.Nameservers,
			CNI: provision.CNIConfig{
				BinPath:  []string{filepath.Join(set.cniDir, "bin")},
				ConfDir:  filepath.Join(set.cniDir, "conf.d"),
				CacheDir: filepath.Join(set.cniDir, "cache"),

				BundleURL: set.options.CNIBundleURL,
			},
		},

		SelfExecutable: set.options.TalosctlPath,
		StateDirectory: set.stateDir,
	}

	for i := 0; i < set.options.Nodes; i++ {
		request.Nodes = append(request.Nodes,
			provision.NodeRequest{
				Name:     fmt.Sprintf("pxe-%d", i),
				Type:     machine.TypeUnknown,
				IPs:      []net.IP{ips[i+1]},
				Memory:   set.options.MemMB * 1024 * 1024,
				NanoCPUs: set.options.CPUs * 1000 * 1000 * 1000,
				Disks: []*provision.Disk{
					{
						Size: uint64(set.options.DiskGB) * 1024 * 1024 * 1024,
					},
				},
				PXEBooted:           true,
				TFTPServer:          set.options.BootSource.String(),
				IPXEBootFilename:    "undionly.kpxe",
				SkipInjectingConfig: true,
				DefaultBootOrder:    set.options.DefaultBootOrder,
			})
	}

	set.cluster, err = set.provisioner.Create(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

// TearDown the set of VMs.
func (set *Set) TearDown(ctx context.Context) error {
	if set.cluster != nil {
		if err := set.provisioner.Destroy(ctx, set.cluster); err != nil {
			return err
		}

		set.cluster = nil
	}

	return nil
}

// BridgeIP returns the IP of the gateway (bridge).
func (set *Set) BridgeIP() net.IP {
	return set.bridgeIP
}

// Nodes return information about PXE VMs.
func (set *Set) Nodes() []provision.NodeInfo {
	return set.cluster.Info().ExtraNodes
}
