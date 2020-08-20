---
description: ""
sidebar: "docs"
prev: "/docs/architecture/"
next: "/docs/environments/"
---

# Concepts

The Metal Controller Manager provides a few custom resources (CRDs) in the management plane Kubernetes cluster that are crucial to understanding the flow of Sidero:

## Environments

These define a desired deployment environment for Talos.
This includes things like which kernel to use, kernel args to pass, and the initrd to use.

## Servers

These represent physical machines as resources in the management plane.
These servers are created when the physical machine PXE boots and completes a "discovery" process in which it registers with the management plane and provides SMBIOS information.

## ServerClasses

ServerClasses are a grouping of the Servers mentioned above.
These can be used to compose a bank of Servers that are eligible for provisioning.

## Metadata

The metadata server may be familiar to you if you have used cloud environments previously.
Given Talos machine configurations created by the Sidero Cluster API provider, along with patches specified by editing server resources, metadata is returned to servers who query the metadata server with their UUID as input.
