---
description: "A guide for bootstrapping Sidero management plane"
weight: 7
title: "Import Workload Machines"
---

At this point, any servers on the same network as Sidero should network boot from Sidero.
To register a server with Sidero, simply turn it on and Sidero will do the rest.
Once the registration is complete, you should see the servers registered with `kubectl get servers`:

```bash
$ kubectl get servers -o wide
NAME                                   HOSTNAME        ACCEPTED   ALLOCATED   CLEAN
00000000-0000-0000-0000-d05099d33360   192.168.1.201   false      false       false
```

## Accept the Servers

Note in the output above that the newly registered servers are not `accepted`.
In order for a server to be eligible for consideration, it _must_ be marked as `accepted`.
Before a `Server` is accepted, no write action will be performed against it.
This default is for safety (don't accidentally delete something just because it
was plugged in) and security (make sure you know the machine before it is given
credentials to communicate).

> Note: if you are running in a safe environment, you can configure Sidero to
> automatically accept new machines.

For more information on server acceptance, see the [server docs](../../resource-configuration/servers/#server-acceptance).

## Create ServerClasses

By default, Sidero comes with a single ServerClass `any` which matches any
(accepted) server.
This is sufficient for this demo, but you may wish to have
more flexibility by defining your own ServerClasses.

ServerClasses allow you to group machines which are sufficiently similar to
allow for unnamed allocation.
This is analogous to cloud providers using such classes as `m3.large` or
`c2.small`, but the names are free-form and only need to make sense to you.

For more information on ServerClasses, see the [ServerClass
docs](../../resource-configuration/serverclasses/).

## Hardware differences

In baremetal systems, there are commonly certain small features and
configurations which are unique to the hardware.
In many cases, such small variations may not require special configurations, but
others do.

If hardware-specific differences do mandate configuration changes, we need a way
to keep those changes local to the hardware specification so that at the higher
level, a Server is just a Server (or a server in a ServerClass is just a Server
like all the others in that Class).

The most common variations seem to be the installation disk and the console
serial port.

Some machines have NVMe drives, which show up as something like `/dev/nvme0n1`.
Others may be SATA or SCSI, which show up as something like `/dev/sda`.
Some machines use `/dev/ttyS0` for the serial console; others `/dev/ttyS1`.

Configuration patches can be applied to either Servers or ServerClasses, and
those patches will be applied to the final machine configuration for those
nodes without having to know anything about those nodes at the allocation level.

For examples of install disk patching, see the [Installation Disk
doc](../../resource-configuration/servers/#installation-disk).

For more information about patching in general, see the [Patching
Guide](../../guides/patching).
