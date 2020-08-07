---
description: ""
sidebar: "docs"
prev: "/docs/concepts/"
next: "/docs/servers/"
---

# Environments

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
