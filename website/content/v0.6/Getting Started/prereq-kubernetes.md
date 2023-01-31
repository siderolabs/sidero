---
description: "Prerequisite: Kubernetes"
weight: 3
title: "Prerequisite: Kubernetes"
---

In order to run Sidero, you first need a Kubernetes "cluster".
There is nothing special about this cluster.
It can be, for example:

- a Kubernetes cluster you already have
- a single-node cluster running in Docker on your laptop
- a cluster running inside a virtual machine stack such as VMWare
- a Talos Kubernetes cluster running on a spare machine

Two important things are needed in this cluster:

- Kubernetes `v1.19` or later
- Ability to expose TCP and UDP Services to the workload cluster machines

For the purposes of this tutorial, we will create this cluster in Docker on a
workstation, perhaps a laptop.

If you already have a suitable Kubernetes cluster, feel free to skip this step.

## Create a Local Management Cluster

The `talosctl` CLI tool has built-in support for spinning up Talos in docker containers.
Let's use this to our advantage as an easy Kubernetes cluster to start from.

Issue the following to create a single-node Docker-based Kubernetes cluster:

```bash
export HOST_IP="192.168.1.150"

talosctl cluster create \
  --name sidero-demo \
  -p 69:69/udp,8081:8081/tcp,51821:51821/udp \
  --workers 0 \
  --config-patch '[{"op": "add", "path": "/cluster/allowSchedulingOnControlPlanes", "value": true}]' \
  --endpoint $HOST_IP
```

The `192.168.1.150` IP address should be changed to the IP address of your Docker
host.
This is _not_ the Docker bridge IP but the standard IP address of the
workstation.

Note that there are three ports mentioned in the command above.
The first (69) is
for TFTP.
The second (8081) is for the web server (which serves netboot
artifacts and configuration).
The third (51821) is for the SideroLink Wireguard network.

Exposing them here allows us to access the services that will get deployed on this node.
In turn, we will be running our Sidero services with `hostNetwork: true`,
so the Docker host will forward these to the Docker container,
which will in turn be running in the same namespace as the Sidero Kubernetes components.
A full separate management cluster will likely approach this differently,
with a load balancer or a means of sharing an IP address across multiple nodes (such as with MetalLB).

Finally, the `--config-patch` is optional,
but since we are running a single-node cluster in this Tutorial,
adding this will allow Sidero to run on the controlplane.
Otherwise, you would need to add worker nodes to this management plane cluster to be
able to run the Sidero components on it.

## Access the cluster

Once the cluster create command is complete, you can retrieve the kubeconfig for it using the Talos API:

```bash
talosctl kubeconfig
```

> Note: by default, Talos will merge the kubeconfig for this cluster into your
> standard kubeconfig under the context name matching the cluster name your
> created above.
> If this name conflicts, it will be given a `-1`, a `-2` or so
> on, so it is generally safe to run.
> However, if you would prefer to not modify your standard kubeconfig, you can
> supply a directory name as the third parameter, which will cause a new
> kubeconfig to be created there instead.
> Remember that if you choose to not use the standard location, your should set
> your `KUBECONFIG` environment variable or pass the `--kubeconfig` option to
> tell the `kubectl` client the name of the `kubeconfig` file.
