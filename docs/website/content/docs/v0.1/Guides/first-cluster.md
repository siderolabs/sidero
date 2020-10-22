---
description: "A guide for creating your first cluster with the Sidero management plane"
weight: 2
---

# Creating Your First Cluster

## Introduction

This guide will detail the steps needed to provision your first bare metal Talos cluster after completing the bootstrap and pivot steps detailed in the previous guide.
There will be two main steps in this guide: reconfiguring the Sidero components now that they have been pivoted and the actual cluster creation.

## Reconfigure Sidero

### Patch Services

In this guide, we will convert the metadata service to a NodePort service and the other services to use host networking.
This is also necessary because some protocols like TFTP don't allow for port configuration.
Along with some nodeSelectors and a scale up of the metal controller manager deployment, creating the services this way allows for the creation of DNS names that point to all management plane nodes and provide an HA experience if desired.
It should also be noted, however, that there are many options for acheiving this functionality.
Users can look into projects like MetalLB or KubeRouter with BGP and ECMP if they desire something else.

Metal Controller Manager:

```bash
## Use host networking
kubectl patch deploy -n sidero-system sidero-controller-manager --type='json' -p='[{"op": "add", "path": "/spec/template/spec/hostNetwork", "value": true}]'
```

Metadata Server:

```bash
# Convert metadata server service to nodeport
kubectl patch service -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "replace", "path": "/spec/type", "value": "NodePort"}]'

## Set a known nodeport for metadata server
kubectl patch service -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "replace", "path": "/spec/ports", "value": [{"port": 80, "protocol": "TCP", "targetPort": "http", "nodePort": 30005}]}]'
```

#### Update Environment

The metadata server's information needs to be updated in the default environment.
Edit the environment with `kubectl edit environment default` and update the `talos.config` kernel arg with the IP of one of the management plane nodes (or the DNS entry you created) and the nodeport we specified above (30005).

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
  filename "ipxe.efi";
}

host talos-mgmt-0 {
   fixed-address 192.168.254.2;
   hardware ethernet d0:50:99:d3:33:60;
}
```

## Register the Servers

At this point, any servers on the same network as Sidero should PXE boot using the Sidero PXE service.
To register a server with Sidero, simply turn it on and Sidero will do the rest.
Once the registration is complete, you should see the servers registered with `kubectl get servers`:

```bash
$ kubectl get servers
NAME                                   AGE
00000000-0000-0000-0000-d05099d33360   3m14s
```

## Create the Cluster

The cluster creation process should be identical to what was detailed in the previous guide.
Note that, for this example, the same "default" serverclass that we used in the previous guide is used again.
Using clusterctl, we can create a cluster manifest with:

```bash
clusterctl config cluster workload-cluster -i sidero > workload-cluster.yaml
```

Note that there are several variables that should be set in order for the templating to work properly:

- CONTROL_PLANE_ENDPOINT
- CONTROL_PLANE_SERVERCLASS
- WORKER_SERVERCLASS
- KUBERNETES_VERSION

Now that we have the manifest, we can simply apply it:

```bash
kubectl apply -f workload-cluster.yaml
```

**NOTE: The templated manifest above is meant to act as a starting point. If customizations are needed to ensure proper setup of your Talos cluster, they should be added before applying.**

Once the workload cluster is setup, you can fetch the talosconfig with a command like:

```bash
kubectl get talosconfig -o yaml workload-cluster-cp-xxx -o jsonpath='{.status.talosConfig}' > workload-cluster-talosconfig.yaml
```

Then the workload cluster's kubeconfig can be fetched with `talosctl --talosconfig workload-cluster-talosconfig.yaml kubeconfig /desired/path`.
