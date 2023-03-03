---
description: ""
weight: 15
title: What's New
---

## New API Version for `metal.sidero.dev` Resources

Resources under `metal.sidero.dev` (`Server`, `ServerClass`, `Environment`) now have a new version `v1alpha2`.
Old version `v1alpha1` is still supported, but it is recommended to update templates to use the new resource version.

### `Server` Changes

Hardware information was restructured and extended when compared with `v1alpha1`:

* `.spec.systemInformation` -> `.spec.hardware.system`
* `.spec.cpu` -> `.spec.hardware.compute.processors[]`

### `ServerClass` Changes

* `.spec.qualifiers.systemInformation` -> `.spec.qualifiers.system`
* `.spec.qualifiers.cpu` -> `.spec.qualifiers.hardware.compute.processors[]`

## Metadata Server

Sidero Metadata Server no longer depends on the version of Talos machinery library it is built with.
Sidero should be able to process machine config for future versions of Talos.

## Sidero Agent

Sidero Agent now runs DHCP client in the userland, on the link which was used to PXE boot the machine.
This allows to run Sidero Agent on the machine with several autoconfigured network interfaces, when one of them is used for the management network.

## DHCP Proxy

Sidero Controller Manager now includes DHCP proxy which augments DHCP response with additional PXE boot options.
When enabled, DHCP server in the environment only handles IP allocation and network configuration, while DHCP proxy
provides PXE boot information automatically based on the architecture and boot method.
