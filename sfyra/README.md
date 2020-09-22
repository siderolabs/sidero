# sfyra

Integration test for Sidero/Arges.

## Running

Build the test binary:

    make sfyra

Run the test:

    make run-sfyra

Registry mirrors could be dropped if not being used.
Test uses CIDR `172.24.0.0/24` by default.

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

From Sfyra directory run:

    (cd ../talos/; sudo -E ../sfyra/_out/integration-test -skip-teardown)

This command doesn't tear down the cluster after the test run, so it can be re-run any time for quick another round of testing.

To destroy Sfyra environment use `talosctl`:

    sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra
    sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra-management

## Running with Sidero HEAD

Build Sidero and push to registry under your username:

    make USERNAME=smira PUSH=true
    make release USERNAME=smira PUSH=true

Create/update `clusterctl` config file to install Sidero for that version:

    $ cat ~/.cluster-api/clusterctl.yaml
    providers:
    - name: "sidero"
        url: "file:///home/smira/Documents/sidero/_out/infrastructure-sidero/v0.1.0-alpha.1-12-g8f9ba14-dirty/infrastructure-components.yaml"
        type: "InfrastructureProvider"

Update the path to match your directory layout.
