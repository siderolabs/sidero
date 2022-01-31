---
description: "A guide describing patching"
weight: 3
title: "Patching"
---

Server resources can be updated by using the `configPatches` section of the custom resource.
Any field of the [Talos machine config](https://www.talos.dev/docs/v0.13/reference/configuration/)
can be overridden on a per-machine basis using this method.
The format of these patches is based on [JSON 6902](http://jsonpatch.com/) that you may be used to in tools like kustomize.

Any patches specified in the server resource are processed by the Sidero controller before it returns a Talos machine config for a given server at boot time.

A set of patches may look like this:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: Server
metadata:
  name: 00000000-0000-0000-0000-d05099d33360
spec:
  configPatches:
    - op: replace
      path: /machine/install
      value:
        disk: /dev/sda
    - op: replace
      path: /cluster/network/cni
      value:
        name: "custom"
        urls:
          - "http://192.168.1.199/assets/cilium.yaml"
```

## Testing Configuration Patches

While developing config patches it is usually convenient to test generated config with patches
before actual server is provisioned with the config.

This can be achieved by querying the metadata server endpoint directly:

```sh
$ curl http://$PUBLIC_IP:8081/configdata?uuid=$SERVER_UUID
version: v1alpha1
...
```

Replace `$PUBLIC_IP` with the Sidero IP address and `$SERVER_UUID` with the name of the `Server` to test
against.

If metadata endpoint returns an error on applying JSON patches, make sure config subtree being patched exists in the config.
If it doesn't exist, create it with the `op: add` above the `op: replace` patch.

## Combining Patches from Multiple Sources

Config patches might be combined from multiple sources (`Server`, `ServerClass`, `TalosControlPlane`, `TalosConfigTemplate`), which is explained in details
in [Metadata](../../resource-configuration/metadata/) section.
