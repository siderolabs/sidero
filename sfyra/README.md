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

## Running with Talos HEAD

Build the artifacts in Talos:

    make initramfs kernel talosctl-linux

From Sidero directory run:

    (cd ../talos/; sudo -E ../sidero/_out/sfyra -skip-teardown)

This command doesn't tear down the cluster after the test run, so it can be re-run any time for quick another round of testing.

To destroy Sfyra environment use `talosctl`:

    sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra
    sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra-management
