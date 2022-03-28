---
description: ""
weight: 50
title: Resources
---

Sidero, the Talos bootstrap/controlplane providers, and Cluster API each provide several custom resources (CRDs) to Kubernetes.
These CRDs are crucial to understanding the connections between each provider and in troubleshooting problems.
It may also help to look at the [cluster template](https://github.com/siderolabs/sidero/blob/master/templates/cluster-template.yaml) to get an idea of the relationships between these.

---

## Cluster API (CAPI)

It's worth defining the most basic resources that CAPI provides first, as they are related to several subsequent resources below.

### `Cluster`

`Cluster` is the highest level CAPI resource.
It allows users to specify things like network layout of the cluster, as well as contains references to the infrastructure and control plane resources that will be used to create the cluster.

### `Machines`

`Machine` represents an infrastructure component hosting a Kubernetes node.
Allows for specification of things like Kubernetes version, as well as contains reference to the infrastructure resource that relates to this machine.

### `MachineDeployments`

`MachineDeployments` are similar to a `Deployment` and their relationship to `Pods` in Kubernetes primitives.
A `MachineDeployment` allows for specification of a number of Machine replicas with a given specification.

---

## Cluster API Bootstrap Provider Talos (CABPT)

### `TalosConfigs`

The `TalosConfig` resource allows a user to specify the type (init, controlplane, join) for a given machine.
The bootstrap provider will then generate a Talos machine configuration for that machine.
This resource also provides the ability to pass a full, pre-generated machine configuration.
Finally, users have the ability to pass `configPatches`, which are applied to edit a generate machine configuration with user-defined settings.
The `TalosConfig` corresponds to the `bootstrap` sections of Machines, `MachineDeployments`, and the `controlPlaneConfig` section of `TalosControlPlanes`.

### `TalosConfigTemplates`

`TalosConfigTemplates` are similar to the `TalosConfig` above, but used when specifying a bootstrap reference in a `MachineDeployment`.

---

## Cluster API Control Plane Provider Talos (CACPPT)

### `TalosControlPlanes`

The control plane provider presents a single CRD, the `TalosControlPlane`.
This resource is similar to `MachineDeployments`, but is targeted exclusively for the Kubernetes control plane nodes.
The `TalosControlPlane` allows for specification of the number of replicas, version of Kubernetes for the control plane nodes, references to the infrastructure resource to use (`infrastructureTemplate` section), as well as the configuration of the bootstrap data via the `controlPlaneConfig` section.
This resource is referred to by the CAPI Cluster resource via the `controlPlaneRef` section.

---

## Sidero

### Cluster API Provider Sidero (CAPS)

#### `MetalClusters`

A `MetalCluster` is Sidero's view of the cluster resource.
This resource allows users to define the control plane endpoint that corresponds to the Kubernetes API server.
This resource corresponds to the `infrastructureRef` section of Cluster API's `Cluster` resource.

#### `MetalMachines`

A `MetalMachine` is Sidero's view of a machine.
Allows for reference of a single server or a server class from which a physical server will be picked to bootstrap.

`MetalMachine` provides a set of statuses describing the state (available with SideroLink, requires Talos >= 0.14):

```yaml
status:
  addresses:
    - address: 172.25.0.5
        type: InternalIP
    - address: pxe-2
        type: Hostname
  conditions:
    - lastTransitionTime: "2022-02-11T14:20:42Z"
      message: 'Get ... connection refused'
      reason: ProviderUpdateFailed
      severity: Warning
      status: "False"
      type: ProviderSet
    - lastTransitionTime: "2022-02-11T12:48:35Z"
      status: "True"
      type: TalosConfigLoaded
    - lastTransitionTime: "2022-02-11T12:48:35Z"
      status: "True"
      type: TalosConfigValidated
    - lastTransitionTime: "2022-02-11T12:48:35Z"
      status: "True"
      type: TalosInstalled
```

Statuses:

- `addresses` lists the current IP addresses and hostname of the node, `addresses` are updated when the node addresses are changed
- `conditions`:
  - `ProviderSet`: captures the moment infrastrucutre provider ID is set in the `Node` specification; depends on workload cluster control plane availability
  - `TalosConfigLoaded`: Talos successfully loaded machine configuration from Sidero; if this condition indicates a failure, check `sidero-controller-manager` logs
  - `TalosConfigValidated`: Talos successfully validated machine configuration; a failure in this condition indicates that the machine config is malformed
  - `TalosInstalled`: Talos was successfully installed to disk

#### `MetalMachineTemplates`

A `MetalMachineTemplate` is similar to a `MetalMachine` above, but serves as a template that is reused for resources like `MachineDeployments` or `TalosControlPlanes` that allocate multiple `Machines` at once.

#### `ServerBindings`

`ServerBindings` represent a one-to-one mapping between a Server resource and a `MetalMachine` resource.
A `ServerBinding` is used internally to keep track of servers that are allocated to a Kubernetes cluster and used to make decisions on cleaning and returning servers to a `ServerClass` upon deallocation.

### Metal Controller Manager

#### `Environments`

These define a desired deployment environment for Talos, including things like which kernel to use, kernel args to pass, and the initrd to use.
Sidero allows you to define a default environment, as well as other environments that may be specific to a subset of nodes.
Users can override the environment at the `ServerClass` or `Server` level, if you have requirements for different kernels or kernel parameters.

See the [Environments](../../resource-configuration/environments/) section of our Configuration docs for examples and more detail.

#### `Servers`

These represent physical machines as resources in the management plane.
These `Servers` are created when the physical machine PXE boots and completes a "discovery" process in which it registers with the management plane and provides SMBIOS information such as the CPU manufacturer and version, and memory information.

See the [Servers](../../resource-configuration/servers/) section of our Configuration docs for examples and more detail.

#### `ServerClasses`

`ServerClasses` are a grouping of the `Servers` mentioned above, grouped to create classes of servers based on Memory, CPU or other attributes.
These can be used to compose a bank of `Servers` that are eligible for provisioning.

See the [ServerClasses](../../resource-configuration/serverclasses/) section of our Configuration docs for examples and more detail.

### Sidero Controller Manager

While the controller does not present unique CRDs within Kubernetes, it's important to understand the metadata resources that are returned to physical servers during the boot process.

#### Metadata

The Sidero controller manager server may be familiar to you if you have used cloud environments previously.
Using Talos machine configurations created by the Talos Cluster API bootstrap provider, along with patches specified by editing `Server`/`ServerClass` resources or `TalosConfig`/`TalosControlPlane` resources, metadata is returned to servers who query the controller manager at boot time.

See the [Metadata](../../resource-configuration/metadata/) section of our Configuration docs for examples and more detail.
