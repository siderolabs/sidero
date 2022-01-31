---
description: ""
weight: 4
title: Metadata
---

The Sidero controller manager manages the Machine metadata.
In terms of Talos (the OS on which the Kubernetes cluster is formed), this is the
"[machine config](https://www.talos.dev/docs/v0.13/reference/configuration/)",
which is used during the automated installation.

## Talos Machine Configuration

The configuration of each machine is constructed from a number of sources:

- The `TalosControlPlane` custom resource for control plane nodes.
- The `TalosConfigTemplate` custom resource.
- The `ServerClass` which was used to select the `Server` into the `Cluster`.
- Any `Server`-specific patches.

An example usage of setting a virtual IP for the control plane nodes and adding extra `node-labels` to nodes is shown below:

> Note: because of the way JSON patches work the interface setting also needs to be set in `TalosControlPlane` when defining a Virtual IP.
This experience is not ideal, but will be addressed in a future release.

*TalosControlPlane* custom resource:

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
kind: TalosControlPlane
metadata:
  name: workload-cluster
  namespace: default
spec:
  controlPlaneConfig:
    controlplane:
      configPatches:
      - op: add
        path: /machine/network
        value:
          interfaces:
          - interface: eth0
            dhcp: true
            vip:
              ip: 172.16.200.52
      generateType: controlplane
      talosVersion: v0.13
    init:
      configPatches:
      - op: add
        path: /machine/network
        value:
          interfaces:
          - interface: eth0
            dhcp: true
            vip:
              ip: 172.16.200.52
      generateType: init
      talosVersion: v0.13
  infrastructureTemplate:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: MetalMachineTemplate
    name: workload-cluster
  replicas: 3
  version: v1.23.0
```

*TalosConfigTemplate* custom resource:

```yaml
---
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: TalosConfigTemplate
metadata:
  name: workload-cluster
  namespace: default
spec:
  template:
    spec:
      generateType: join
      talosVersion: v0.13
      configPatches:
      - op: add
        path: /machine/kubelet
        value:
          extraArgs:
            node-labels:
              talos.dev/part-of: cluster/workload-cluster
```

and finally in the control plane `ServerClass` custom resource we augment the network information for other interfaces:

```yaml
---
apiVersion: metal.sidero.dev/v1alpha2
kind: ServerClass
metadata:
  name: cp.small.x86
spec:
  configPatches:
  - op: replace
    path: /machine/install/disk
    value: /dev/nvme0n1
  - op: add
    path: /machine/install/extraKernelArgs
    value:
    - console=tty0
    - console=ttyS1,115200n8
  - op: add
    path: /machine/network/interfaces/-
    value:
      interface: eth1
      dhcp: true
  qualifiers:
    - system:
        manufacturer: Supermicro
      compute:
        processors:
          - productName: Intel(R) Xeon(R) E-2124G CPU @ 3.40GHz
      memory:
        totalSize: 8 GB
  selector:
    matchLabels:
      metal.sidero.dev/serverclass: cp.small.x86
```

the workload `ServerClass` defines the complete networking config

```yaml
---
apiVersion: metal.sidero.dev/v1alpha2
kind: ServerClass
metadata:
  name: general.medium.x86
spec:
  configPatches:
  - op: replace
    path: /machine/install/disk
    value: /dev/nvme1n1
  - op: add
    path: /machine/install/extraKernelArgs
    value:
    - console=tty0
    - console=ttyS1,115200n8
  - op: add
    path: /machine/network
    value:
      interfaces:
      - interface: eth0
        dhcp: true
      - interface: eth1
        dhcp: true
  qualifiers:
    - system:
        manufacturer: Supermicro
      compute:
        processors:
          - productName: Intel(R) Xeon(R) E-2136 CPU @ 3.30GHz
      memory:
        totalSize: 16 GB
  selector:
    matchLabels:
      metal.sidero.dev/serverclass: general.medium.x86
```

The base template is constructed from the Talos bootstrap provider, using data from the associated `TalosControlPlane` and `TalosConfigTemplate` manifest.
Then, any configuration patches are applied from the `ServerClass` and `Server`.

These patches take the form of an [RFC 6902](https://tools.ietf.org/html/rfc6902) JSON (or YAML) patch.
An example of the use of this patch method can be found in [Patching Guide](../../guides/patching/).

Also note that while a `Server` can be a member of any number of `ServerClass`es, only the `ServerClass` which is used to select the `Server` into the `Cluster` will be used for the generation of the configuration of the `Machine`.
In this way, `Servers` may have a number of different configuration patch sets based on which `Cluster` they are in at any given time.
