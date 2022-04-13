---
title: System Requirements
---

## System Requirements

Most of the time, Sidero does very little, so it needs very few resources.
However, since it is in charge of any number of workload clusters, it **should**
be built with redundancy.
It is also common, if the cluster is single-purpose,
to combine the controlplane and worker node roles.
Virtual machines are also
perfectly well-suited for this role.

Minimum suggested dimensions:

- Node count: 3
- Node RAM: 4GB
- Node CPU: ARM64 or x86-64 class
- Node storage: 32GB storage on system disk
