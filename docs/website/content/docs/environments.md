---
description: ""
sidebar: "docs"
prev: "/docs/concepts/"
next: "/docs/servers/"
---

# Environments

Environments are a custom resource provided by the Metal Controller Manager.
An environment is a codified description of what should be returned by the PXE server when a physical server attempts to PXE boot.

Especially important in the environment types are the kernel args.
From here, one can tweak the IP to the metadata server as well as various other kernel options that [Talos](https://www.talos.dev/docs/v0.6/en/guides/metal/overview#kernel-parameters) and/or the Linux kernel supports.

A sample environment definition looks like this:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Environment
metadata:
  name: default
spec:
  kernel:
    url: "https://github.com/talos-systems/talos/releases/download/v0.6.0-beta.0/vmlinuz"
    sha512: ""
    args:
      - initrd=initramfs.xz
      - page_poison=1
      - slab_nomerge
      - slub_debug=P
      - pti=on
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
    url: "https://github.com/talos-systems/talos/releases/download/v0.6.0-beta.0/initramfs.xz"
    sha512: ""
```
