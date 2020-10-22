---
description: ""
weight: 1
---

# Introduction

Sidero ("Iron" in Greek) is a project created by the [Talos Systems](https://www.talos-systems.com/) team.
The goal of this project is to provide lightweight, composable tools that can be used to create bare-metal Talos + Kubernetes clusters.
These tools are built around the Cluster API project.
Sidero is also a subproject of Talos Systems' [Arges](https://github.com/talos-systems/arges) project, which will publish known-good versions of these components (along with others) with each release.

## Overview

Sidero is made currently made up of three components:

- Metal Metadata Server: Provides a Cluster API (CAPI)-aware metadata server
- Metal Controller Manager: Provides custom resources and controllers for managing the lifecycle of metal machines
- Cluster API Provider Sidero (CAPS): A Cluster API infrastructure provider that makes use of the pieces above to spin up Kubernetes clusters

Sidero also needs these co-requisites in order to be useful:

- [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- [Cluster API Control Plane Provider Talos](https://github.com/talos-systems/cluster-api-control-plane-provider-talos)
- [Cluster API Bootstrap Provider Talos](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos)

These components and Sidero are all installed using Cluster API's `clusterctl` tool.
Since Sidero is built on top of Cluster API, an existing Kuberntes cluster is required for a "management plane".
This cluster does not have to be based on Talos.
However, if you would like to use Talos as the OS of choice for the Sidero management plane, you can find a number of ways to deploy Talos in the [documentation](https://www.talos.dev/docs/v0.6/en/guides/getting-started/intro).
