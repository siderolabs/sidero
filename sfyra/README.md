# sfyra

Integration test for Sidero/Arges.

## Running

Build the test binary and Sidero, push images:

    make USERNAME=<username> PUSH=true

Run the test (this will trigger `make release`):

    make run-sfyra USERNAME=<username>

Test uses CIDRs `172.24.0.0/24`, `172.25.0.0/24` by default.

Sequence of steps:

* build initial bootstrap Talos cluster of one node
* build management set of VMs (PXE-boot enabled)
* install Cluster API, Sidero and Talos providers
* run the unit-tests

With `-skip-teardown` flag test leaves the bootstrap cluster running so that next iteration of the test
can be run without waiting for the boostrap actions to be finished.

## Running manually

Download the Talos artifacts with `make talos-artifacts`.
Build Sfyra with `make sfyra`.

Run full Sfyra integration test with latest Sidero release (unless overridden in `~/.cluster-api/clusterctl.yaml`):

    sudo -E _out/sfyra test integration

One can also run parts of the test flow:

* setup Talos bootstrap cluster (single node): `sudo -E _out/sfyra bootstrap cluster`
* install and patch CAPI and providers: `_out/sfyra bootstrap capi`
* launch a set of VMs ready for PXE booting: `sudo -E _out/sfyra bootstrap servers`

See each command help on how to customize the operations.

## Running with Talos HEAD

Build the artifacts in Talos:

    make initramfs kernel talosctl-linux

From Sidero directory run:

    sudo -E _out/sfyra test integration --skip-teardown --bootstrap-initramfs=../talos/_out/initramfs-amd64.xz --bootstrap-vmlinuz=../talos/_out/vmlinuz-amd64 --talosctl-path=../talos/_out/talosctl-linux-amd64

This command doesn't tear down the cluster after the test run, so it can be re-run any time for quick another round of testing.

## Cleaning up

To destroy Sfyra environment use `talosctl`:

    sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra
    sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra-management
