---
description: "A guide describing patching"
weight: 3
---

# Patching

Server resources can be updated by using the `configPatches` section of the custom resource.
Any field of the Talos machine config can be overridden on a per-machine basis using this method.
The format of these patches is based on [JSON 6902](http://jsonpatch.com/) that you may be used to in tools like kustomize.

Any patches specified in the server resource are processed by the Metal Metadata Server before it returns a Talos machine config for a given server at boot time.

A set of patches may look like this:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: Server
metadata:
  name: 00000000-0000-0000-0000-d05099d33360
spec:
configPatches:
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
```
