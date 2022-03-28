---
description: ""
weight: 1
title: Introduction
---

Sidero ("Iron" in Greek) is a project created by the [Siderolabs](https://www.siderolabs.com/) team.
The goal of this project is to provide lightweight, composable tools that can be used to create bare-metal Talos + Kubernetes clusters.
These tools are built around the Cluster API project.
Sidero is also a subproject of Siderolabs' [Arges](https://github.com/siderolabs/arges) project, which will publish known-good versions of these components (along with others) with each release.

## Overview

Sidero is made currently made up of three components:

- Metal Metadata Server: Provides a Cluster API (CAPI)-aware metadata server
- Metal Controller Manager: Provides custom resources and controllers for managing the lifecycle of metal machines
- Cluster API Provider Sidero (CAPS): A Cluster API infrastructure provider that makes use of the pieces above to spin up Kubernetes clusters

Sidero also needs these co-requisites in order to be useful:

- [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- [Cluster API Control Plane Provider Talos](https://github.com/siderolabs/cluster-api-control-plane-provider-talos)
- [Cluster API Bootstrap Provider Talos](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos)

All componenets mentioned above can be installed using Cluster API's `clusterctl` tool.

Because of the design of Cluster API, there is inherently a "chicken and egg" problem with needing an existing Kubernetes cluster in order to provision the management plane.
Siderolabs and the Cluster API community have created tools to help make this transition easier.
That being said, the management plane cluster does not have to be based on Talos.
If you would, however, like to use Talos as the OS of choice for the Sidero management plane, you can find a number of ways to deploy Talos in the [documentation](https://www.talos.dev).
