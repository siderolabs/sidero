---
description: ""
weight: 2
---

# Servers

Servers are the basic resource of bare metal in the Metal Controller Manager.
These are created by PXE booting the servers and allowing them to send a registration request to the management plane.

An example server may look like the following:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Server
metadata:
  name: 00000000-0000-0000-0000-d05099d333e0
spec:
  accepted: false
  configPatches:
    - op: replace
      path: /cluster/network/cni
      value:
        name: custom
        urls:
          - http://192.168.1.199/assets/cilium.yaml
  cpu:
    manufacturer: Intel(R) Corporation
    version: Intel(R) Atom(TM) CPU C3558 @ 2.20GHz
  system:
    family: Unknown
    manufacturer: Unknown
    productName: Unknown
    serialNumber: Unknown
    skuNumber: Unknown
    version: Unknown
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
You can enable auto-acceptance by pasing the `--auto-accept-servers=true` flag to `sidero-controller-manager`.

Once accepted, a server will be reset (all disks wiped) and then made available to Sidero.

You should never change an accepted `Server` to be _not_ accepted while it is in use.
Because servers which are not accepted will not be modified, if a server which
_was_ accepted is changed to _not_ accepted, the disk will _not_ be wiped upon
its exit.

