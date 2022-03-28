---
description: ""
weight: 1
title: Introduction
---

Sidero ("Iron" in Greek) is a project created by the [Sidero Labs](https://www.SideroLabs.com/) team.
Sidero Metal provides lightweight, composable tools that can be used to create bare-metal [Talos Linux](https://www.talos.dev) + Kubernetes clusters.
These tools are built around the Cluster API project.

Because of the design of Cluster API, there is inherently a "chicken and egg" problem: you need an existing Kubernetes cluster in order to provision the management plane, that can then provision more clusters.
The initial management plane cluster that runs the Sidero Metal provider does not need to be based on Talos Linux - although it is recommended for security and stability reasons.
The [Getting Started](../../getting-started/) guide will walk you through installing Sidero Metal either on an existing cluster, or by quickly creating a docker based cluster used to bootstrap the process.

## Overview

Sidero Metal is currently made up of two components:

- Metal Controller Manager: Provides custom resources and controllers for managing the lifecycle of metal machines, iPXE server, metadata service, and gRPC API service
- Cluster API Provider Sidero (CAPS): A Cluster API infrastructure provider that makes use of the pieces above to spin up Kubernetes clusters

Sidero Metal also needs these co-requisites in order to be useful:

- [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- [Cluster API Control Plane Provider Talos](https://github.com/siderolabs/cluster-api-control-plane-provider-talos)
- [Cluster API Bootstrap Provider Talos](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos)

All components mentioned above can be installed using Cluster API's `clusterctl` tool.
See the [Getting Started](../../getting-started/) for more details.
