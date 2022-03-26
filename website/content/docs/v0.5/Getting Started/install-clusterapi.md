---
description: "Install Sidero"
weight: 5
title: "Install Sidero"
---

Sidero is included as a default infrastructure provider in `clusterctl`, so the
installation of both Sidero and the Cluster API (CAPI) components is as simple
as using the `clusterctl` tool.

> Note: Because Cluster API upgrades are _stateless_, it is important to keep all Sidero
> configuration for reuse during upgrades.

Sidero has a number of configuration options which should be supplied at install
time, kept, and reused for upgrades.
These can also be specified in the `clusterctl` configuration file
(`$HOME/.cluster-api/clusterctl.yaml`).
You can reference the `clusterctl`
[docs](https://cluster-api.sigs.k8s.io/clusterctl/configuration.html#clusterctl-configuration-file)
for more information on this.

For our purposes, we will use environment variables for our configuration
options.

```bash
export SIDERO_CONTROLLER_MANAGER_HOST_NETWORK=true
export SIDERO_CONTROLLER_MANAGER_DEPLOYMENT_STRATEGY=Recreate
export SIDERO_CONTROLLER_MANAGER_API_ENDPOINT=192.168.1.150
export SIDERO_CONTROLLER_MANAGER_SIDEROLINK_ENDPOINT=192.168.1.150

clusterctl init -b talos -c talos -i sidero
```

First, we are telling Sidero to use `hostNetwork: true` so that it binds its
ports directly to the host, rather than being available only from inside the
cluster.
There are many ways of exposing the services, but this is the simplest
path for the single-node management cluster.
When you scale the management cluster, you will need to use an alternative
method, such as an external load balancer or something like
[MetalLB](https://metallb.universe.tf).

The `192.168.1.150` IP address is the IP address or DNS hostname as seen from the workload
clusters.
In our case, this should be the main IP address of your Docker
workstation.
