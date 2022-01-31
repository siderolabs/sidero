---
description: ""
weight: 2
title: Servers
---

Servers are the basic resource of bare metal in the Metal Controller Manager.
These are created by PXE booting the servers and allowing them to send a registration request to the management plane.

An example server may look like the following:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: Server
metadata:
  name: 00000000-0000-0000-0000-d05099d333e0
  labels:
    common-label: "true"
    zone: east
    environment: test
spec:
  accepted: false
  configPatches:
    - op: replace
      path: /cluster/network/cni
      value:
        name: custom
        urls:
          - http://192.168.1.199/assets/cilium.yaml
  hardware:
    system:
      manufacturer: Dell Inc.
      productName: PowerEdge R630
      serialNumber: 790H8D2
    compute:
      totalCoreCount: 8
      totalThreadCount: 16
      processorCount: 1
      processors:
        - manufacturer: Intel
          productName: Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz
          speed: 2400
          coreCount: 8
          threadCount: 16
    memory:
      totalSize: 32 GB
      moduleCount: 2
      modules:
        - manufacturer: 002C00B3002C
          productName: 18ASF2G72PDZ-2G3B1
          serialNumber: 12BDC045
          type: LPDDR3
          size: 16384
          speed: 2400
        - manufacturer: 002C00B3002C
          productName: 18ASF2G72PDZ-2G3B1
          serialNumber: 12BDBF5D
          type: LPDDR3
          size: 16384
          speed: 2400
    storage:
      totalSize: 1116 GB
      deviceCount: 1
      devices:
        - productName: PERC H730 Mini
          type: HDD
          name: sda
          deviceName: /dev/sda
          size: 1199101181952
          wwid: naa.61866da055de070028d8e83307cc6df2
    network:
      interfaceCount: 2
      interfaces:
        - index: 1
          name: lo
          flags: up|loopback
          mtu: 65536
          mac: ""
          addresses:
            - 127.0.0.1/8
            - ::1/128
        - index: 2
          name: enp3s0
          flags: up|broadcast|multicast
          mtu: 1500
          mac: "40:8d:5c:86:5a:14"
          addresses:
            - 192.168.2.8/24
            - fe80::dcb3:295c:755b:91bb/64
```

## Installation Disk

An installation disk is required by Talos on bare metal.
This can be specified in a `configPatch`:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: Server
...
spec:
  accepted: false
  configPatches:
    - op: replace
      path: /machine/install/disk
      value: /dev/sda
```

The install disk patch can also be set on the `ServerClass`:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: ServerClass
...
spec:
  configPatches:
    - op: replace
      path: /machine/install/disk
      value: /dev/sda
```

## Server Acceptance

In order for a server to be eligible for consideration, it _must_ be `accepted`.
This is an important separation point which all `Server`s must pass.
Before a `Server` is accepted, no write action will be performed against it.
Thus, it is safe for a computer to be added to a network on which Sidero is operating.
Sidero will never write to or wipe any disk on a computer which is not marked as `accepted`.

This can be tedious for systems in which all attached computers should be considered to be under the control of Sidero.
Thus, you may also choose to automatically accept any machine into Sidero on its discovery.
Please keep in mind that this means that any newly-connected computer **WILL BE WIPED** automatically.
You can enable auto-acceptance by passing the `--auto-accept-servers=true` flag to `sidero-controller-manager`.

Once accepted, a server will be reset (all disks wiped) and then made available to Sidero.

You should never change an accepted `Server` to be _not_ accepted while it is in use.
Because servers which are not accepted will not be modified, if a server which
_was_ accepted is changed to _not_ accepted, the disk will _not_ be wiped upon
its exit.

## IPMI

Sidero can use IPMI information to control `Server` power state, reboot servers and set boot order.

IPMI information will be, by default, setup automatically if possible as part of the acceptance process.
In this design, a "sidero" user will be added to the IPMI user list and a randomly generated password will be issued.
This information is then squirreled away in a Kubernetes secret in the `sidero-system` namespace, with a name format of `<server-uuid>-bmc`.
Users wishing to turn off this feature can pass the `--auto-bmc-setup=false` flag to `sidero-controller-manager`

IMPI connection information can also be set manually in the `Server` spec after initial registration:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: Server
...
spec:
  bmc:
    endpoint: 10.0.0.25
    user: admin
    pass: password
```

If IPMI information is set, server boot order might be set to boot from disk, then network, Sidero will switch servers
to PXE boot once that is required.

Without IPMI info, Sidero can still register servers, wipe them and provision clusters, but Sidero won't be able to reboot servers once they are removed from the cluster.
**If IPMI info is not set, servers should be configured to boot first from network, then from disk.**

Sidero can also fetch IPMI credentials via the `Secret` reference:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: Server
...
spec:
  bmc:
    endpoint: 10.0.0.25
    userFrom:
      secretKeyRef:
        name: ipmi-credentials
        key: username
    passFrom:
      secretKeyRef:
        name: ipmi-credentials
        key: password
```

As the `Server` resource is not namespaced, `Secret` should be created in the `default` namespace.

## Other Settings

### `cordoned`

If `cordoned` is set to `true`, `Server` gets excluded from any `ServerClass` it might match based on qualifiers.
This means that the `Server` will not be allocated automatically.

`Server` might be `cordoned` to temporarily take it out of the `ServerClass` to perform for example hardware maintenance.

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Server
...
spec:
  cordoned: true
```

### `pxeBootAlways`

`Server` might be forced to boot from the network even if the OS is already installed with `pxeBootAlways: true`:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Server
...
spec:
  pxeBootAlways: true
```

### `bootFromDiskMethod`

The method to exit iPXE network boot to force boot from disk can be configured for the `Server`:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Server
...
spec:
  bootFromDiskMethod: ipxe-sanboot
```

Valid values are:

- `ipxe-exit`
- `http-404`
- `ipxe-sanboot`

If not set, the `ServerClass.spec.bootFromDiskMethod` value is used with the fallback to the default boot from disk method  (`SIDERO_CONTROLLER_MANAGER_BOOT_FROM_DISK_METHOD`).
