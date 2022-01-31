---
description: ""
weight: 3
title: Server Classes
---

Server classes are a way to group distinct server resources.
The `qualifiers` and `selector` keys allow the administrator to specify criteria upon which to group these servers.
If both of these keys are missing, the server class matches all servers that it is watching.
If both of these keys define requirements, these requirements are combined (logical `AND`).

## `selector`

`selector` groups server resources by their labels.
The [Kubernetes documentation][label-selector-docs] has more information on how to use this field.

## `qualifiers`

A list of hardware criteria, where each entry in the list is interpreted as a logical `OR`.
All criteria inside each entry is interpreted as a logical `AND`.
Qualifiers that are not specified are not evaluated.

An example:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: ServerClass
metadata:
  name: serverclass-sample
spec:
  selector:
    matchLabels:
      common-label: "true"
    matchExpressions:
      - key: zone
        operator: In
        values:
          - central
          - east
      - key: environment
        operator: NotIn
        values:
          - prod
  qualifiers:
    hardware:
      - system:
          manufacturer: Dell Inc.
        compute:
          processors:
            - manufacturer: Advanced Micro Devices, Inc.
              productName: AMD Ryzen 7 2700X Eight-Core Processor
      - compute:
          processors:
            - manufacturer: "Intel(R) Corporation"
              productName: "Intel(R) Atom(TM) CPU C3558 @ 2.20GHz"
        memory:
          totalSize: "8 GB"
```

Servers would only be added to the above class if they:

- have the label `common-label` with value `true`
- _AND_ match the `matchExpressions`
- _AND_ match either 1 of the following criteria:
  - has a system manufactured by `Dell Inc.` _AND_ has at least 1 processor that is an `AMD Ryzen 7 2700X Eight-Core Processor`
  - has at least 1 processor that is an `Intel(R) Atom(TM) CPU C3558 @ 2.20GHz` _AND_ has exactly 8 GB of total memory

Additionally, Sidero automatically creates and maintains a server class called `"any"` that includes all (accepted) servers.
Attempts to add qualifiers to it will be reverted.

[label-selector-docs]: https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/label-selector/

## `configPatches`

Server configs of servers matching a server class can be updated by using the `configPatches` section of the custom resource.
See [patching](../guides/patching) for more information on how this works.

An example of settings the default install disk for all servers matching a server class:

```yaml
apiVersion: metal.sidero.dev/v1alpha2
kind: ServerClass
...
spec:
  configPatches:
    - op: replace
      path: /machine/install/disk
      value: /dev/sda
```

## Other Settings

### `environmentRef`

Servers from a `ServerClass` can be set to use the specific `Environment` by linking the `Environment` from the `ServerClass`:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: ServerClass
...
spec:
  environmentRef:
    name: production-env
```

### `bootFromDiskMethod`

The method to exit iPXE network boot to force boot from disk can be configured for all `Server` resources belonging to the `ServerClass`:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: ServerClass
...
spec:
  bootFromDiskMethod: ipxe-sanboot
```

Valid values are:

- `ipxe-exit`
- `http-404`
- `ipxe-sanboot`

If not set, the default boot from disk method is used (`SIDERO_CONTROLLER_MANAGER_BOOT_FROM_DISK_METHOD`).
