---
description: "A guide for bootstrapping Sidero management plane"
weight: 11
title: "Optional: Management cluster"
---

Having the Sidero cluster running inside a Docker container is not the most
robust place for it, but it did make for an expedient start.

It might be a good idea to move it to a more robust setup, such as a dedicated bare-metal server,
or a virtual machine.
It also makes sense to establish regular backups of `etcd` in the management cluster to make sure the state of the cluster can be restored in case of a disaster.
