---
description: ""
weight: 1
title: Environments
---

Environments are a custom resource provided by the Metal Controller Manager.
An environment is a codified description of what should be returned by the PXE server when a physical server attempts to PXE boot.

Especially important in the environment types are the kernel args.
From here, one can tweak the IP to the metadata server as well as various other kernel options that [Talos](https://www.talos.dev/docs/v0.8/introduction/getting-started/#kernel-parameters) and/or the Linux kernel supports.

Environments can be supplied to a given server either at the Server or the ServerClass level.
The hierarchy from most to least respected is:

- `.spec.environmentRef` provided at `Server` level
- `.spec.environmentRef` provided at `ServerClass` level
- `"default"` `Environment` created by administrator

A sample environment definition looks like this:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Environment
metadata:
  name: default
spec:
  kernel:
    url: "https://github.com/siderolabs/talos/releases/download/v0.8.1/vmlinuz-amd64"
    sha512: ""
    args:
      - init_on_alloc=1
      - init_on_free=1
      - slab_nomerge
      - pti=on
      - consoleblank=0
      - random.trust_cpu=on
      - ima_template=ima-ng
      - ima_appraise=fix
      - ima_hash=sha512
      - console=tty0
      - console=ttyS1,115200n8
      - earlyprintk=ttyS1,115200n8
      - panic=0
      - printk.devkmsg=on
      - talos.platform=metal
      - talos.config=http://$PUBLIC_IP:9091/configdata?uuid=
  initrd:
    url: "https://github.com/siderolabs/talos/releases/download/v0.8.1/initramfs-amd64.xz"
    sha512: ""
```

Example of overriding `"default"` `Environment` at the `Server` level:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
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
apiVersion: metal.sidero.dev/v1alpha1
kind: ServerClass
...
spec:
  environmentRef:
    namespace: default
    name: boot
  ...
```
