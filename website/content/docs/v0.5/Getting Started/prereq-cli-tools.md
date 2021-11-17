---
description: "Prerequisite: CLI tools"
weight: 2
title: "Prerequisite: CLI tools"
---

You will need three CLI tools installed on your workstation in order to interact
with Sidero:

- `kubectl`
- `clusterctl`
- `talosctl`

## Install `kubectl`

Since `kubectl` is the standard Kubernetes control tool, many distributions
already exist for it.
Feel free to check your own package manager to see if it is available natively.

Otherwise, you may install it directly from the main distribution point.
The main article for this can be found
[here](https://kubernetes.io/docs/tasks/tools/#kubectl).

```bash
sudo curl -Lo /usr/local/bin/kubectl \
  "https://dl.k8s.io/release/$(\
  curl -L -s https://dl.k8s.io/release/stable.txt\
  )/bin/linux/amd64/kubectl"
sudo chmod +x /usr/local/bin/kubectl
```

## Install `clusterctl`

The `clusterctl` tool is the standard control tool for ClusterAPI (CAPI).
It is less common, so it is also less likely to be in package managers.

The main article for installing `clusterctl` can be found
[here](https://cluster-api.sigs.k8s.io/user/quick-start.html#install-clusterctl).

```bash
sudo curl -Lo /usr/local/bin/clusterctl \
  "https://github.com/kubernetes-sigs/cluster-api/releases/download/v0.4.4/clusterctl-$(uname -s | tr '[:upper:]' '[:lower:]')-amd64" \
sudo chmod +x /usr/local/bin/clusterctl
```

> Note: This version of Sidero is only compatible with CAPI v1alpha4,
> so versions of `clusterctl` above v0.4.x will not work.

## Install `talosctl`

The `talosctl` tool is used to interact with the Talos (our Kubernetes-focused
operating system) API.
The latest version can be found on our
[Releases](https://github.com/talos-systems/talos/releases) page.

```bash
sudo curl -Lo /usr/local/bin/talosctl \
 "https://github.com/talos-systems/talos/releases/latest/download/talosctl-$(uname -s | tr '[:upper:]' '[:lower:]')-amd64"
chmod +x /usr/local/bin/talosctl
```
