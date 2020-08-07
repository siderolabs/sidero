---
description: ""
sidebar: "docs"
prev: "/docs/"
next: "/docs/architecture/"
---

# Installation

## Configuring `clusterctl`

You will need to add Sidero to `~/.cluster-api/clusterctl.yaml`:

```yaml
providers:
  - name: "sidero"
    url: "file:///home/andrewrynhard/workspace/code/github.com/talos-systems/sidero/_out/infrastructure-sidero/v0.1.0-alpha.0-61f6621-dirty/components.yaml"
    type: "InfrastructureProvider"
```

## Installing Sidero

```bash
clusterctl init -b talos -c talos -i sidero
```
