---
description: ""
weight: 1
title: Environments
---

Environments are a custom resource provided by the Metal Controller Manager.
An environment is a codified description of what should be returned by the PXE server when a physical server attempts to PXE boot.

Especially important in the environment types are the kernel args.
From here, one can tweak the IP to the metadata server as well as various other kernel options that [Talos](https://www.talos.dev/docs/v0.13/reference/kernel/#commandline-parameters) and/or the Linux kernel supports.

Environments can be supplied to a given server either at the Server or the ServerClass level.
The hierarchy from most to least respected is:

- `.spec.environmentRef` provided at `Server` level
- `.spec.environmentRef` provided at `ServerClass` level
- `"default"` `Environment` created automatically and modified by an administrator

A sample environment definition looks like this:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: Environment
metadata:
  name: default
spec:
  kernel:
    url: "https://github.com/talos-systems/talos/releases/download/v0.14.0/vmlinuz-amd64"
    sha512: ""
    args:
      - console=tty0
      - console=ttyS1,115200n8
      - consoleblank=0
      - earlyprintk=ttyS1,115200n8
      - ima_appraise=fix
      - ima_hash=sha512
      - ima_template=ima-ng
      - init_on_alloc=1
      - initrd=initramfs.xz
      - nvme_core.io_timeout=4294967295
      - printk.devkmsg=on
      - pti=on
      - random.trust_cpu=on
      - slab_nomerge=
      - talos.platform=metal
  initrd:
    url: "https://github.com/talos-systems/talos/releases/download/v0.14.0/initramfs-amd64.xz"
    sha512: ""
```

Example of overriding `"default"` `Environment` at the `Server` level:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: Server
...
spec:
  environmentRef:
    namespace: default
    name: boot
  ...
```

Example of overriding `"default"` `Environment` at the `ServerClass` level:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: ServerClass
...
spec:
  environmentRef:
    namespace: default
    name: boot
  ...
```
