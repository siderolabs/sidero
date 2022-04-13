---
description: "A guide for decommissioning servers"
weight: 1
title: "Decommissioning Servers"
---

This guide will detail the process for removing a server from Sidero.
The process is fairly simple with a few pieces of information.

- For the given server, take note of any serverclasses that are configured to match the server.

- Take note of any clusters that make use of aforementioned serverclasses.

- For each matching cluster, edit the cluster resource with `kubectl edit cluster` and set `.spec.paused` to `true`.
  Doing this ensures that no new machines will get created for these servers during the decommissioning process.

- If you want to mark a server to be not allocated after it's accepted into the cluster, set the `.spec.cordoned` field to `true`.
  This will prevent the server from being allocated to any new clusters (still allowing it to be wiped).

- If the server is already part of a cluster (`kubectl get serverbindings -o wide` should provide this info), you can now delete the machine that corresponds with this server via `kubectl delete machine <machine_name>`.

- With the machine deleted, Sidero will reboot the machine and wipe its disks.

- Once the disk wiping is complete and the server is turned off, you can finally delete the server from Sidero with `kubectl delete server <server_name>` and repurpose the server for something else.

- Finally, unpause any clusters that were edited in step 3 by setting `.spec.paused` to `false`.
