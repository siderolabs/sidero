# sfyra

Integration test for Sidero/Arges.

## Running

Build the test binary and Sidero, push images:

```sh
make USERNAME=<username> TAG=v0.1.0 PUSH=true
```

Run the test (this will trigger `make release`):

```sh
make run-sfyra USERNAME=<username> TAG=v0.1.0
```

Test uses CIDRs `172.24.0.0/24`, `172.25.0.0/24` by default.

Sequence of steps:

- build initial bootstrap Talos cluster of one node
- build management set of VMs (PXE-boot enabled)
- install Cluster API, Sidero and Talos providers
- run the unit-tests

It's also possible to run Sfyra manually to avoid tearing down and recreating whole environment each time.
After `make USERNAME=<username> TAG=v0.1.0 PUSH=true` run:

```sh
make talos-artifacts # need to run it only once per Talos release change
make clusterctl-release USERNAME=<username> TAG=v0.1.0 PUSH=true
```

Then launch Sfyra manually with desired flags:

```sh
sudo -E _out/sfyra test integration --registry-mirror docker.io=http://172.24.0.1:5000,k8s.gcr.io=http://172.24.0.1:5001,quay.io=http://172.24.0.1:5002,gcr.io=http://172.24.0.1:5003,ghcr.io=http://172.24.0.1:5004,127.0.0.1:5005=http://172.24.0.1:5005 --skip-teardown --clusterctl-config ~/.cluster-api/clusterctl.sfyra.yaml
```

Alternatively, you may use `run-sfyra` target with `SFYRA_EXTRA_FLAGS` and `REGISTRY_MIRROR_FLAGS` environment variables:

```sh
export USERNAME=<username>
export TAG=v0.1.0
export REGISTRY_MIRROR_FLAGS="--registry-mirror docker.io=http://172.24.0.1:5000,k8s.gcr.io=http://172.24.0.1:5001,quay.io=http://172.24.0.1:5002,gcr.io=http://172.24.0.1:5003,ghcr.io=http://172.24.0.1:5004,127.0.0.1:5005=http://172.24.0.1:5005"
export SFYRA_EXTRA_FLAGS="--skip-teardown"
make run-sfyra
```

With `--skip-teardown` flag test leaves the bootstrap cluster running so that next iteration of the test can be run without waiting for the boostrap actions to be finished.
It's possible to run Sfyra tests once again without bringing down the test environment, but make sure that all the clusters are deleted with `kubectl delete clusters --all`.

Flag `--registry-mirror` is optional, but it speeds up provisioning significantly.
See Talos guides on setting up registry pull-through caches, or just run `hack/start-registry-proxies.sh`.

Kubernetes config can be pulled with `talosctl -n 172.24.0.2 kubeconfig --force`.

When `sfyra` is not running, loadbalancer for `management-cluster` control plane is also down, it can be restarted for manual testing with `_out/sfyra loadbalancer create --kubeconfig=$HOME/.kube/config --load-balancer-port 10000`.

One can also run parts of the test flow:

- setup Talos bootstrap cluster (single node): `sudo -E _out/sfyra bootstrap cluster`
- install and patch CAPI and providers: `_out/sfyra bootstrap capi`
- launch a set of VMs ready for PXE booting: `sudo -E _out/sfyra bootstrap servers`

See each command help on how to customize the operations.

## Testing Always PXE Boot

By default, QEMU VMs provisioned to emulate metal servers are configured to boot from the disk first, and Sidero uses API call to force PXE boot to run the agent.

Sometimes it's important to test the flow when the servers are configured to boot from the network first always (e.g. if bare metal setup doesn't have IPMI), in that case it's important to force VMs to boot from the network always.
This can be achieved by adding a flag `--default-boot-order=nc` to `sfyra` invocation.
In this case Sidero iPXE server will force VM to boot from disk via iPXE if the server is already provisioned.

> Note: due to the dependency on new `talosctl`, this feature will be available once Talos in Sfyra is updated to version >= 0.11.

## Running with Talos HEAD

Build the artifacts in Talos:

```sh
make initramfs kernel talosctl-linux
```

From Sidero directory run:

```sh
sudo -E _out/sfyra test integration --skip-teardown --bootstrap-initramfs=../talos/_out/initramfs-amd64.xz --bootstrap-vmlinuz=../talos/_out/vmlinuz-amd64 --talosctl-path=../talos/_out/talosctl-linux-amd64
```

This command doesn't tear down the cluster after the test run, so it can be re-run any time for quick another round of testing.

## Cleaning up

To destroy Sfyra environment use `talosctl`:

```sh
sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra
sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra-management
```
