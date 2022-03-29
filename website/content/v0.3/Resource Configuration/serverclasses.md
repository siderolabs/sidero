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

There are currently two keys: `cpu`, `systemInformation`.
Each of these keys accepts a list of entries.
The top level keys are a "logical `AND`", while the lists under each key are a "logical `OR`".
Qualifiers that are not specified are not evaluated.

An example:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
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
    cpu:
      - manufacturer: "Intel(R) Corporation"
        version: "Intel(R) Atom(TM) CPU C3558 @ 2.20GHz"
      - manufacturer: Advanced Micro Devices, Inc.
        version: AMD Ryzen 7 2700X Eight-Core Processor
    systemInformation:
      - manufacturer: Dell Inc.
```

Servers would only be added to the above class if they:

- had _EITHER_ CPU info
- _AND_ the label key/value in `matchLabels`
- _AND_ match the `matchExpressions`

Additionally, Sidero automatically creates and maintains a server class called `"any"` that includes all (accepted) servers.
Attempts to add qualifiers to it will be reverted.

[label-selector-docs]: https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/label-selector/

## `configPatches`

Server configs of servers matching a server class can be updated by using the `configPatches` section of the custom resource.
See [patching](/docs/v0.3/guides/patching) for more information on how this works.

An example of settings the default install disk for all servers matching a server class:

```yaml
apiVersion: metal.sidero.dev/v1alpha1
kind: ServerClass
...
spec:
  configPatches:
    - op: replace
      path: /machine/install/disk
      value: /dev/sda
```
