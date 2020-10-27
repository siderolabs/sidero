---
description: ""
weight: 4
---

# Concepts

The Metal Controller Manager provides a few custom resources (CRDs) in the management plane Kubernetes cluster that are crucial to understanding the flow of Sidero:

## Environments

These define a desired deployment environment for Talos, including things like which kernel to use, kernel args to pass, and the initrd to use.
Sidero allows you to define a default environment, as well as other environments that may be specific to a subset of nodes.
Users can override the environment at the ServerClass or Server level, if you have requirements for different kernels or kernel parameters.

See the [Environments](/docs/v0.1/Configuration/environments.md) section of our Configuration docs for examples and more detail.

## Servers

These represent physical machines as resources in the management plane.
These servers are created when the physical machine PXE boots and completes a "discovery" process in which it registers with the management plane and provides SMBIOS information such as the CPU manufacturer and version, and memory information.

See the [Servers](/docs/v0.1/Configuration/servers.md) section of our Configuration docs for examples and more detail.

## ServerClasses

ServerClasses are a grouping of the Servers mentioned above, grouped to create classes of servers based on Memory, CPU or other attributes.
These can be used to compose a bank of Servers that are eligible for provisioning.

See the [ServerClasses](/docs/v0.1/Configuration/serverclasses.md) section of our Configuration docs for examples and more detail.

## Metadata

The metadata server may be familiar to you if you have used cloud environments previously.
Using Talos machine configurations created by the Talos Cluster API bootstrap provider, along with patches specified by editing Server/ServerClass resources, metadata is returned to servers who query the metadata server at boot time.

See the [Metadata](/docs/v0.1/Configuration/metadata.md) section of our Configuration docs for examples and more detail.
