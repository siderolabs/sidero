---
description: "A guide for creating your first cluster with the Sidero management plane"
weight: 2
title: "Creating Your First Cluster"
---

## Introduction

This guide will detail the steps needed to provision your first bare metal Talos cluster after completing the bootstrap and pivot steps detailed in the previous guide.
There will be two main steps in this guide: reconfiguring the Sidero components now that they have been pivoted and the actual cluster creation.

## Reconfigure Sidero

### Patch Services

In this guide, we will convert the services to use host networking.
This is also necessary because some protocols like TFTP don't allow for port configuration.
Along with some nodeSelectors and a scale up of the metal controller manager deployment, creating the services this way allows for the creation of DNS names that point to all management plane nodes and provide an HA experience if desired.
It should also be noted, however, that there are many options for achieving this functionality.
Users can look into projects like MetalLB or KubeRouter with BGP and ECMP if they desire something else.

Metal Controller Manager:

```bash
## Use host networking
kubectl patch deploy -n sidero-system sidero-controller-manager --type='json' -p='[{"op": "add", "path": "/spec/template/spec/hostNetwork", "value": true}]'
```

#### Update Environment

<!-- textlint-disable -->

Sidero by default appends `talos.config` kernel argument with based on the flags `--api-endpoint` and `--api-port` to the `sidero-controller-manager`:
`talos.config=http://$API_ENDPOINT:$API_PORT/configdata?uuid=`.

<!-- textlint-enable -->

If this default value doesn't apply, edit the environment with `kubectl edit environment default` and add the `talos.config` kernel arg with the IP of one of the management plane nodes (or the DNS entry you created).

### Update DHCP

The DHCP options configured in the previous guide should now be updated to point to your new management plane IP or to the DNS name if it was created.

A revised ipxe-metal.conf file looks like:

```bash
allow bootp;
allow booting;

next-server 192.168.254.2;
if exists user-class and option user-class = "iPXE" {
  filename "http://192.168.254.2:8081/boot.ipxe";
} else {
  if substring (option vendor-class-identifier, 15, 5) = "00000" {
    # BIOS
    if substring (option vendor-class-identifier, 0, 10) = "HTTPClient" {
      option vendor-class-identifier "HTTPClient";
      filename "http://192.168.254.2:8081/tftp/undionly.kpxe";
    } else {
      filename "undionly.kpxe";
    }
  } else {
    # UEFI
    if substring (option vendor-class-identifier, 0, 10) = "HTTPClient" {
      option vendor-class-identifier "HTTPClient";
      filename "http://192.168.254.2:8081/tftp/snp.efi";
    } else {
      filename "snp.efi";
    }
  }
}

host talos-mgmt-0 {
   fixed-address 192.168.254.2;
   hardware ethernet d0:50:99:d3:33:60;
}
```

There are multiple ways to boot the via iPXE:

- if the node has built-in iPXE, direct URL to the iPXE script can be used: `http://192.168.254.2:8081/boot.ipxe`.
- depending on the boot mode (BIOS or UEFI), either `snp.efi` or `undionly.kpxe` can be used (these images contain embedded iPXE scripts).
- iPXE binaries can be delivered either over TFTP or HTTP (HTTP support depends on node firmware).

## Register the Servers

At this point, any servers on the same network as Sidero should PXE boot using the Sidero PXE service.
To register a server with Sidero, simply turn it on and Sidero will do the rest.
Once the registration is complete, you should see the servers registered with `kubectl get servers`:

```bash
$ kubectl get servers -o wide
NAME                                   HOSTNAME        ACCEPTED   ALLOCATED   CLEAN
00000000-0000-0000-0000-d05099d33360   192.168.254.2   false      false       false
```

## Accept the Servers

Note in the output above that the newly registered servers are not `accepted`.
In order for a server to be eligible for consideration, it _must_ be marked as `accepted`.
Before a `Server` is accepted, no write action will be performed against it.
Servers can be accepted by issuing a patch command like:

```bash
kubectl patch server 00000000-0000-0000-0000-d05099d33360 --type='json' -p='[{"op": "replace", "path": "/spec/accepted", "value": true}]'
```

For more information on server acceptance, see the [server docs](../../resource-configuration/servers).

## Create the Cluster

The cluster creation process should be identical to what was detailed in the previous guide.
Using clusterctl, we can create a cluster manifest with:

```bash
clusterctl generate cluster workload-cluster -i sidero > workload-cluster.yaml
```

Note that there are several variables that should be set in order for the templating to work properly:

- `CONTROL_PLANE_ENDPOINT` and `CONTROL_PLANE_PORT`: The endpoint (IP address or hostname) and the port used for the Kubernetes API server
  (e.g. for `https://1.2.3.4:6443`: `CONTROL_PLANE_ENDPOINT=1.2.3.4` and `CONTROL_PLANE_PORT=6443`).
  This is the equivalent of the `endpoint` you would specify in `talosctl gen config`.
  There are a variety of ways to configure a control plane endpoint.
  Some common ways for an HA setup are to use DNS, a load balancer, or BGP.
  A simpler method is to use the IP of a single node.
  This has the disadvantage of being a single point of failure, but it can be a simple way to get running.
- `CONTROL_PLANE_SERVERCLASS`: The server class to use for control plane nodes.
- `WORKER_SERVERCLASS`: The server class to use for worker nodes.
- `KUBERNETES_VERSION`: The version of Kubernetes to deploy (e.g. `v1.19.4`).
- `TALOS_VERSION`: This should correspond to the minor version of Talos that you will be deploying (e.g. `v0.10`).
  This value is used in determining the fields present in the machine configuration that gets generated for Talos nodes.
  Note that the default is currently `v0.13`.

Now that we have the manifest, we can simply apply it:

```bash
kubectl apply -f workload-cluster.yaml
```

**NOTE: The templated manifest above is meant to act as a starting point.**
**If customizations are needed to ensure proper setup of your Talos cluster, they should be added before applying.**

Once the workload cluster is setup, you can fetch the talosconfig with a command like:

```bash
kubectl get talosconfig -o yaml workload-cluster-cp-xxx -o jsonpath='{.status.talosConfig}' > workload-cluster-talosconfig.yaml
```

Then the workload cluster's kubeconfig can be fetched with `talosctl --talosconfig workload-cluster-talosconfig.yaml kubeconfig /desired/path`.
