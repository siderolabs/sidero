# sfyra

Integration test for Sidero/Arges.

## Running

It is recommended to run the test suite with a local registry mirror running (refer to Talos Linux development documentation for details).

Test uses CIDRs `172.24.0.0/24` by default.

Sequence of steps:

- build initial bootstrap Talos cluster of two nodes: bootstrap-controlplane and bootstrap-worker
- build management set of VMs (PXE-boot enabled)
- install Cluster API, Sidero and Talos providers
- run the integratoion tests

Build and push Sidero images to the local registry:

```sh
make REGISTRY=127.0.0.1:5005 TAG=v0.6.0 PUSH=true
```

It's possible to run Sfyra manually to avoid tearing down and recreating whole environment each time.
After building Sidero Metal and Sfyra, run:

```sh
# build sfyra
make sfyra
make talos-artifacts # need to run it only once per Talos release change
make clusterctl-release USERNAME=<username> TAG=v0.6.0 PUSH=true
```

There is a `clusterctl` config file overriding images for the bootstrap, controlplane and infrastructure providers to point to the local registry mirror:

```yaml
# ~/.cluster-api/clusterctl.sfyra.yaml
# IMPORTANT: change the paths to point to your local builds
providers:
  - name: "talos"
    url: "file:///home/smira/Documents/cluster-api-bootstrap-provider-talos/_out/bootstrap-talos/v0.6.0/bootstrap-components.yaml"
    type: "BootstrapProvider"
  - name: "talos"
    url: "file:///home/smira/Documents/cluster-api-control-plane-provider-talos/_out/control-plane-talos/v0.5.0/control-plane-components.yaml"
    type: "ControlPlaneProvider"
  - name: "sidero"
    url: "file:///home/smira/Documents/sidero/_out/infrastructure-sidero/v0.6.0/infrastructure-components.yaml"
    type: "InfrastructureProvider"
```

In order to make CABPT and CACPPT local builds, you need to do the following in each of those projects:

```sh
# CABPT:
cd cluster-api-bootstrap-provider-talos/
make REGISTRY=127.0.0.1:5005 PUSH=true PLATFORM=linux/amd64 TAG=v0.6.0
# CACPPT:
cd cluster-api-control-plane-provider-talos/
make REGISTRY=127.0.0.1:5005 PUSH=true PLATFORM=linux/amd64 TAG=v0.5.0
```

> Note: CACPPT is using v0.5.0, while CABPT is using v0.6.0, the major-minor pair should match the last release of the respective provider, but the patch release can be left as zero.

Then launch Sfyra manually with desired flags:

```sh
sudo --preserve-env=HOME _out/sfyra test integration --registry-mirrors docker.io=http://172.24.0.1:5000,gcr.io=http://172.24.0.1:5003,ghcr.io=http://172.24.0.1:5004,127.0.0.1:5005=http://172.24.0.1:5005,registry.k8s.io=http://172.24.0.1:5001 --skip-teardown --clusterctl-config ~/.cluster-api/clusterctl.sfyra.yaml
```

With `--skip-teardown` flag test leaves the bootstrap cluster running so that next iteration of the test can be run without waiting for the bootstrap actions to be finished.
It's possible to run Sfyra tests once again without bringing down the test environment, but make sure that all the clusters are deleted with `kubectl delete clusters --all`.

Flag `--registry-mirror` is optional, but it speeds up provisioning significantly.

Kubernetes config can be pulled with `talosctl -n 172.24.0.2 kubeconfig --force`.

## Useful commands

Watch the cluster status (update `management-cluster` with your cluster name if different):

```sh
clusterctl describe cluster management-cluster --show-conditions=all
# or:
watch -n 2 clusterctl describe cluster management-cluster --show-conditions=all
```

Get current state of servers:

```sh
kubectl get servers
```

Allocation of servers to machines:

```sh
kubectl get serverbindings
```

Current clusters:

```sh
kubectl get clusters
```

When re-running the tests without full teardown, it's useful to clean up existing clusters:

```sh
kubectl delete clusters --all
```

Following the logs of Sidero Metal controllers:

```sh
kubectl logs -n sidero-system deployment/sidero-controller-manager -c manager -f
kubectl logs -n sidero-system deployment/caps-controller-manager -c manager -f
```

Watch PXE VMs logs:

```sh
tail -F ~/.talos/clusters/sfyra/pxe-*.log
```

## Testing Always PXE Boot

By default, QEMU VMs provisioned to emulate metal servers are configured to boot from the disk first, and Sidero uses API call to force PXE boot to run the agent.

Sometimes it's important to test the flow when the servers are configured to boot from the network first always (e.g. if bare metal setup doesn't have IPMI), in that case it's important to force VMs to boot from the network always.
This can be achieved by adding a flag `--default-boot-order=nc` to `sfyra` invocation.
In this case Sidero iPXE server will force VM to boot from disk via iPXE if the server is already provisioned.

## Cleaning up

To destroy Sfyra environment use `talosctl`:

```sh
sudo --preserve-env=HOME talosctl cluster destroy --name=sfyra
```
