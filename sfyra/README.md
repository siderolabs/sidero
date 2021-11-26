# sfyra

Integration test for Sidero/Arges.

## Running

It is recommended to run the test suite with a local registry mirror running.

To run a local mirror run `hack/start-registry-proxies.sh`.
This will create a local registry along with mirrors for `registry-1.docker.io`, `k8s.gcr.io`, `https://quay.io`, `https://gcr.io` and `https://ghcr.io`

Build the test binary and Sidero, push images:

> If you have the local registry mirrors running add `REGISTRY=127.0.0.1:5010` to all the make commands (This will speed up development a lot).

```sh
make USERNAME=<username> TAG=v0.4.0 PUSH=true
```

Run the test (this will trigger `make release`):

```sh
make run-sfyra USERNAME=<username> TAG=v0.4.0
```

Test uses CIDRs `172.24.0.0/24`, `172.25.0.0/24` by default.

Sequence of steps:

- build initial bootstrap Talos cluster of one node
- build management set of VMs (PXE-boot enabled)
- install Cluster API, Sidero and Talos providers
- run the unit-tests

It's also possible to run Sfyra manually to avoid tearing down and recreating whole environment each time.
After `make USERNAME=<username> TAG=v0.4.0 PUSH=true` run:

```sh
# build sfyra
make sfyra
make talos-artifacts # need to run it only once per Talos release change
make clusterctl-release USERNAME=<username> TAG=v0.4.0 PUSH=true
```

Then launch Sfyra manually with desired flags:

```sh
sudo -E _out/sfyra test integration --registry-mirror docker.io=http://172.24.0.1:5000,k8s.gcr.io=http://172.24.0.1:5001,quay.io=http://172.24.0.1:5002,gcr.io=http://172.24.0.1:5003,ghcr.io=http://172.24.0.1:5004,127.0.0.1:5005=http://172.24.0.1:5005 --skip-teardown --clusterctl-config ~/.cluster-api/clusterctl.sfyra.yaml
```

Alternatively, you may use `run-sfyra` target with `SFYRA_EXTRA_FLAGS` and `REGISTRY_MIRROR_FLAGS` environment variables:

```sh
export USERNAME=<username>
export TAG=v0.4.0
export REGISTRY_MIRROR_FLAGS="--registry-mirror docker.io=http://172.24.0.1:5000,k8s.gcr.io=http://172.24.0.1:5001,quay.io=http://172.24.0.1:5002,gcr.io=http://172.24.0.1:5003,ghcr.io=http://172.24.0.1:5004,127.0.0.1:5005=http://172.24.0.1:5005"
export SFYRA_EXTRA_FLAGS="--skip-teardown"
make run-sfyra
```

With `--skip-teardown` flag test leaves the bootstrap cluster running so that next iteration of the test can be run without waiting for the bootstrap actions to be finished.
It's possible to run Sfyra tests once again without bringing down the test environment, but make sure that all the clusters are deleted with `kubectl delete clusters --all`.

Flag `--registry-mirror` is optional, but it speeds up provisioning significantly.

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

## Running with Talos HEAD as a bootstrap cluster

Build the artifacts in Talos:

```sh
make initramfs kernel talosctl-linux
```

From Sidero directory run:

```sh
sudo -E _out/sfyra test integration --skip-teardown --bootstrap-initramfs=../talos/_out/initramfs-amd64.xz --bootstrap-vmlinuz=../talos/_out/vmlinuz-amd64 --talosctl-path=../talos/_out/talosctl-linux-amd64
```

This command doesn't tear down the cluster after the test run, so it can be re-run any time for quick another round of testing.

## Running with Talos HEAD as workload clusters

Build the artifacts in Talos:

```sh
make initramfs kernel
```

Make sure that the `initramfs` has a matching installer image pushed to the registry.

Start _some_ webserver in Talos directory, e.g.

```sh
$ python -m http.server
Serving HTTP on 0.0.0.0 port 8000 (http://0.0.0.0:8000/) ...
```

From Sidero directory run:

```sh
sudo -E _out/sfyra test integration --skip-teardown --talos-initrd-url http://172.24.0.1:8000/_out/initramfs-amd64.xz --talos-kernel-url http://172.24.0.1:8000/_out/vmlinuz-amd64
```

## Cleaning up

To destroy Sfyra environment use `talosctl`:

```sh
sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra
sudo -E talosctl cluster destroy --provisioner=qemu --name=sfyra-management
```

## Manually registering a server for testing

### Registering

```bash
grpcurl \
    -proto app/sidero-controller-manager/internal/api/api.proto \
    -plaintext \
    -d '{"hostname":"fake","cpu":{"manufacturer":"QEMU","version":"pc-q35-5.2"},"system_information":{"uuid": "a9cf15ab-d96b-4544-b6ab-baebb262213b","family":"Unknown","manufacturer":"QEMU","productName":"Standard PC (Q35 + ICH9, 2009)","serialNumber":"Unknown","skuNumber":"Unknown","version":"pc-q35-5.2"}}' \
    172.24.0.2:8081 \
    api.Agent/CreateServer
```

### Marking server as cleaned

```bash
grpcurl \
    -proto app/sidero-controller-manager/internal/api/api.proto \
    -plaintext \
    -d '{"uuid": "a9cf15ab-d96b-4544-b6ab-baebb262213b"}' \
    172.24.0.2:8081 \
    api.Agent/MarkServerAsWiped
```
