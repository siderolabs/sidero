---
title: "Getting Started"
weight: 20
---


This tutorial will walk you through a complete Sidero setup and the formation,
scaling, and destruction of a workload cluster.

To complete this tutorial, you will need a few things:

- ISC DHCP server.
  While any DHCP server will do, we will be presenting the
  configuration syntax for ISC DHCP.
  This is the standard DHCP server available on most Linux distributions (NOT
  dnsmasq) as well as on the Ubiquiti EdgeRouter line of products.
- Machine or Virtual Machine on which to run Sidero itself.
  The requirements for this machine are very low, it can be x86 or arm64
  and it should have at least 4GB of RAM.
- Machines on which to run Kubernetes clusters.
  These have the same minimum specifications as the Sidero machine.
- Workstation on which `talosctl`, `kubectl`, and `clusterctl` can be run.

## Steps

1. Prerequisite: CLI tools
1. Prerequisite: DHCP server
1. Prerequisite: Kubernetes
1. Install Sidero
1. Expose services
1. Import workload machines
1. Create a workload cluster
1. Scale the workload cluster
1. Destroy the workload cluster
1. Optional: Pivot management cluster

## Useful Terms

**ClusterAPI** or **CAPI** is the common system for managing Kubernetes clusters
in a declarative fashion.

**Management Cluster** is the cluster on which Sidero itself runs.
It is generally a special-purpose Kubernetes cluster whose sole responsibility
is maintaining the CRD database of Sidero and providing the services necessary
to manage your workload Kubernetes clusters.

**Sidero** is the ClusterAPI-powered system which manages baremetal
infrastructure for Kubernetes.

**Talos** is the Kubernetes-focused Linux operating system built by the same
people who bring to you Sidero.
It is a very small, entirely API-driven OS which is meant to provide a reliable
and self-maintaining base on which Kubernetes clusters may run.
More information about Talos can be found at
[https://talos.dev](https://talos.dev).

**Workload Cluster** is a cluster, managed by Sidero, on which your Kubernetes
workloads may be run.
The workload clusters are where you run your own applications and infrastructure.
Sidero creates them from your available resources, maintains them over time as
your needs and resources change, and removes them whenever it is told to do so.
