---
description: ""
sidebar: "docs"
prev: "/docs/bootstrapping/"
---

# Patching

```yaml
- op: add
  path: /machine/network/interfaces
  value:
    - interface: eth1
      ignore: true
    - interface: eth2
      ignore: true
    - interface: eth3
      ignore: true
- op: replace
  path: /machine/install
  value:
    disk: /dev/sda
    image: docker.io/autonomy/installer:v0.6.0-beta.0
    bootloader: true
    wipe: false
    force: false
- op: replace
  path: /cluster/network/cni
  value:
    name: "custom"
    urls:
      - "http://192.168.1.199/assets/cilium.yaml"
- op: replace
  path: /cluster/controlPlane/endpoint
  value: https://management-plane.rsmitty.cloud:6443
```
