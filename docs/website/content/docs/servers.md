---
description: ""
sidebar: "docs"
prev: "/docs/environments/"
next: "/docs/serverclasses/"
---

# Servers

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Server
metadata:
  name: 00000000-0000-0000-0000-d05099d333e0
spec:
  configPatches:
    - op: add
      path: /machine/network/interfaces
      value:
        - ignore: true
          interface: eth1
        - ignore: true
          interface: eth2
        - ignore: true
          interface: eth3
    - op: replace
      path: /machine/install
      value:
        bootloader: true
        disk: /dev/sda
        force: false
        image: docker.io/autonomy/installer:latest
        wipe: false
    - op: replace
      path: /cluster/network/cni
      value:
        name: custom
        urls:
          - http://192.168.1.199/assets/cilium.yaml
    - op: replace
      path: /cluster/controlPlane/endpoint
      value: https://management-plane.rsmitty.cloud:6443
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
