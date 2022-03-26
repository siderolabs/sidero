---
description: "Running Sidero on Raspberry Pi 4 to provision bare-metal servers."
title: Sidero on Raspberry Pi 4
weight: 7
---

Sidero doesn't require a lot of computing resources, so SBCs are a perfect fit to run
the Sidero management cluster.
In this guide, we are going to install Talos on Raspberry Pi4, deploy Sidero and other CAPI components.

## Prerequisites

Please see Talos documentation for additional information on [installing Talos on Raspberry Pi4](https://www.talos.dev/docs/v0.13/single-board-computers/rpi_4/).

Download the `clusterctl` CLI  from [CAPI releases](https://github.com/kubernetes-sigs/cluster-api/releases).
The minimum required version is 0.4.3.

## Installing Talos

Prepare the SD card with the Talos RPi4 image, and boot the RPi4.
Talos should drop into maintenance mode printing the acquired IP address.
Record the IP address as the environment variable `SIDERO_ENDPOINT`:

```bash
export SIDERO_ENDPOINT=192.168.x.x
```

> Note: it makes sense to transform DHCP lease for RPi4 into a static reservation so that RPi4 always has the same IP address.

Generate Talos machine configuration for a single-node cluster:

```bash
talosctl gen config --config-patch='[{"op": "add", "path": "/cluster/allowSchedulingOnMasters", "value": true},{"op": "replace", "path": "/machine/install/disk", "value": "/dev/mmcblk0"}]' rpi4-sidero https://${SIDERO_ENDPOINT}:6443/
```

Submit the generated configuration to Talos:

```bash
talosctl apply-config --insecure -n ${SIDERO_ENDPOINT} -f controlplane.yaml
```

Merge client configuration `talosconfig` into default `~/.talos/config` location:

```bash
talosctl config merge talosconfig
```

Update default endpoint and nodes:

```bash
talosctl config endpoints ${SIDERO_ENDPOINT}
talosctl config nodes ${SIDERO_ENDPOINT}
```

You can verify that Talos has booted by running:

```bash
$ talosctl version
talosctl version
Client:
    Tag:         v0.10.3
    SHA:         21018f28
    Built:
    Go version:  go1.16.3
    OS/Arch:     linux/amd64

Server:
    NODE:        192.168.0.31
    Tag:         v0.10.3
    SHA:         8f90c6a8
    Built:
    Go version:  go1.16.3
    OS/Arch:     linux/arm64
```

Bootstrap the etcd cluster:

```bash
talosctl bootstrap
```

At this point, Kubernetes is bootstrapping, and it should be available once all the images are fetched.

Fetch the `kubeconfig` from the cluster with:

```bash
talosctl kubeconfig
```

You can watch the bootstrap progress by running:

```bash
talosctl dmesg -f
```

Once Talos prints `[talos] boot sequence: done`, Kubernetes should be up:

```bash
kubectl get nodes
```

## Installing Sidero

Install Sidero with host network mode, exposing the endpoints on the node's address:

```bash
SIDERO_CONTROLLER_MANAGER_HOST_NETWORK=true SIDERO_CONTROLLER_MANAGER_DEPLOYMENT_STRATEGY=Recreate SIDERO_CONTROLLER_MANAGER_API_ENDPOINT=${SIDERO_IP} clusterctl init -i sidero -b talos -c talos
```

Watch the progress of installation with:

```bash
watch -n 2 kubectl get pods -A
```

Once images are downloaded, all pods should be in running state:

```bash
$ kubectl get pods -A
NAMESPACE             NAME                                         READY   STATUS    RESTARTS   AGE
cabpt-system          cabpt-controller-manager-6458494888-d7lnm    1/1     Running   0          29m
cacppt-system         cacppt-controller-manager-f98854db8-qgkf9    1/1     Running   0          29m
capi-system           capi-controller-manager-58f797cb65-8dwpz     2/2     Running   0          30m
capi-webhook-system   cabpt-controller-manager-85fd964c9c-ldzb6    1/1     Running   0          29m
capi-webhook-system   cacppt-controller-manager-75c479b7f-5hw89    1/1     Running   0          29m
capi-webhook-system   capi-controller-manager-7d596cc4cb-kjrfk     2/2     Running   0          30m
capi-webhook-system   caps-controller-manager-79664cf677-zqbvw     1/1     Running   0          29m
cert-manager          cert-manager-86cb5dcfdd-v86wr                1/1     Running   0          31m
cert-manager          cert-manager-cainjector-84cf775b89-swk25     1/1     Running   0          31m
cert-manager          cert-manager-webhook-7f9f4f8dcb-29xm4        1/1     Running   0          31m
kube-system           coredns-fcc4c97fb-wkxkg                      1/1     Running   0          35m
kube-system           coredns-fcc4c97fb-xzqzj                      1/1     Running   0          35m
kube-system           kube-apiserver-talos-192-168-0-31            1/1     Running   0          33m
kube-system           kube-controller-manager-talos-192-168-0-31   1/1     Running   0          33m
kube-system           kube-flannel-qmlw6                           1/1     Running   0          34m
kube-system           kube-proxy-j24hg                             1/1     Running   0          34m
kube-system           kube-scheduler-talos-192-168-0-31            1/1     Running   0          33m
```

Verify Sidero installation and network setup with:

```bash
$ curl -I http://${SIDERO_ENDPOINT}:8081/tftp/ipxe.efi
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 1020416
Content-Type: application/octet-stream
Last-Modified: Thu, 03 Jun 2021 15:40:58 GMT
Date: Thu, 03 Jun 2021 15:41:51 GMT
```

Now Sidero is installed, and it is ready to be used.
Configure your DHCP server to PXE boot your bare metal servers from `$SIDERO_ENDPOINT` (see [Bootstrapping guide](../bootstrapping/) on DHCP configuration).

## Backup and Recovery

SD cards are not very reliable, so make sure you are taking regular [etcd backups](https://www.talos.dev/docs/v0.13/guides/disaster-recovery/#backup),
so that you can [recover](https://www.talos.dev/docs/v0.13/guides/disaster-recovery/#recovery) your Sidero installation in case of data loss.
