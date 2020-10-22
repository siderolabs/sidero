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
Once accepted, a server will be reset (all disks wiped) and then made available to Sidero.
