---
description: ""
weight: 2
---

# Installation

As of Cluster API version 0.3.9, Sidero is included as a default infrastructure provider in `clusterctl`.

To install Sidero and the other Talos providers, simply issue:

```bash
clusterctl init -b talos -c talos -i sidero
```

Sidero supports several variables to configure the installation, these variables can be set either as environment
variables or as variables in the `clusterctl` configuration:

* `SIDERO_CONTROLLER_MANAGER_HOST_NETWORK` (`false`): run `sidero-controller-manager` on host network
* `SIDERO_CONTROLLER_MANAGER_API_ENDPOINT` (empty): specifies the IP address controller manager can be reached on, defaults to the node IP
* `SIDERO_CONTROLLER_MANAGER_API_PORT` (8081): specifies the port controller manager can be reached on
* `SIDERO_CONTROLLER_MANAGER_EXTRA_AGENT_KERNEL_ARGS` (empty): specifies additional Linux kernel arguments for the Sidero agent (for example, different console settings)
* `SIDERO_CONTROLLER_MANAGER_AUTO_ACCEPT_SERVERS` (`false`): automatically accept discovered servers, by default `.spec.accepted` should be changed to `true` to accept the server
* `SIDERO_CONTROLLER_MANAGER_AUTO_BMC_SETUP` (`true`): automatically attempt to configure the BMC with a `sidero` user that will be used for all IPMI tasks.
* `SIDERO_CONTROLLER_MANAGER_INSECURE_WIPE` (`true`): wipe only the first megabyte of each disk on the server, otherwise wipe the full disk
* `SIDERO_CONTROLLER_MANAGER_SERVER_REBOOT_TIMEOUT` (`20m`): timeout for the server reboot (how long it might take for the server to be rebooted before Sidero retries an IPMI reboot operation)

Sidero provides two endpoints which should be made available to the infrastructure:

* TCP port 8081 which provides combined iPXE, metadata and gRPC service (external endpoint should be passed to Sidero as `SIDERO_CONTROLLER_MANAGER_API_ENDPOINT` and  `SIDERO_CONTROLLER_MANAGER_API_PORT`)
* UDP port 69 for the TFTP service (DHCP server should point the nodes to PXE boot from that IP)

These endpoints could be exposed to the infrastructure using different strategies:

* running `sidero-controller-manager` on the host network.
* using Kubernetes load balancers (e.g. MetalLB), ingress controllers, etc.
