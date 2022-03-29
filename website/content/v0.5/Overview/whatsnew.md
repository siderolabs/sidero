---
description: ""
weight: 15
title: What's New
---

### Cluster API v1.x (v1beta1)

This release of Sidero brings compatibility with CAPI v1.x (v1beta1).

### Cluster Template

Sidero ships with new cluster template without `init` nodes.
This template is only compatible with Talos >= 0.14 (it requires SideroLink feature which was introduced in Talos 0.14).

On upgrade, Sidero supports clusters running Talos < 0.14 if they were created before the upgrade.
Use [legacy template](https://github.com/talos-systems/sidero/blob/release-0.4/templates/cluster-template.yaml) to deploy clusters with Talos < 0.14.

### New `MetalMachines` Conditions

New set of conditions is now available which can simplify cluster troubleshooting:

- `TalosConfigLoaded` is set to false when the config load has failed.
- `TalosConfigValidated` is set to false when the config validation
fails on the node.
- `TalosInstalled` is set to true/false when talos installer finishes.

Requires Talos >= v0.14.

### Machine Addresses

Sidero now populates `MetalMachine` addresses with the ones discovered from Siderolink server events.
Which is then propagated to CAPI `Machine` resources.

Requires Talos >= v0.14.

### SideroLink

Sidero now connects to all servers using SideroLink (available only with Talos >= 0.14).
This enables streaming of kernel logs and events back to Sidero.

All server logs can now be viewed by getting logs of one of the container of the `sidero-controller-manager`:

```bash
kubectl logs -f -n sidero-system deployment/sidero-controller-manager -c serverlogs
```

Events:

```bash
kubectl logs -f -n sidero-system deployment/sidero-controller-manager -c serverevents
```

### iPXE Boot From Disk Method

iPXE boot from disk method can now be set not only on the global level, but also in the Server and ServerClass specs.

### IPMI PXE Method

IPMI PXE method (UEFI, BIOS) can now be configured with `SIDERO_CONTROLLER_MANAGER_IPMI_PXE_METHOD` while installing Sidero.

### Retry PXE Boot

Sidero server controller now keeps track of Talos installation progress.
Now the node will be PXE booted until Talos installation succeeds.
