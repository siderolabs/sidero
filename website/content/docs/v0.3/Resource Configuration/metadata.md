---
description: ""
weight: 4
title: Metadata
---

The Metadata server manages the Machine metadata.
In terms of Talos (the OS on which the Kubernetes cluster is formed), this is the
"[machine config](https://www.talos.dev/docs/v0.11/reference/configuration/)",
which is used during the automated installation.

## Talos Machine Configuration

The configuration of each machine is constructed from a number of sources:

- The Talos bootstrap provider.
- The `Cluster` of which the `Machine` is a member.
- The `ServerClass` which was used to select the `Server` into the `Cluster`.
- Any `Server`-specific patches.

The base template is constructed from the Talos bootstrap provider, using data from the associated `Cluster` manifest.
Then, any configuration patches are applied from the `ServerClass` and `Server`.

Only configuration patches are allowed in the `ServerClass` and `Server` resources.
These patches take the form of an [RFC 6902](https://tools.ietf.org/html/rfc6902) JSON (or YAML) patch.
An example of the use of this patch method can be found in [Patching Guide](../../guides/patching/).

Also note that while a `Server` can be a member of any number of `ServerClass`es, only the `ServerClass` which is used to select the `Server` into the `Cluster` will be used for the generation of the configuration of the `Machine`.
In this way, `Servers` may have a number of different configuration patch sets based on which `Cluster` they are in at any given time.
