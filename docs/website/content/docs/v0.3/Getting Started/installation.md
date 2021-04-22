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
* `SIDERO_CONTROLLER_MANAGER_EXTRA_AGENT_KERNEL_ARGS` (empty): specifies additional Linux kernel arguments for the Sidero agent (for example, different console settings)
* `SIDERO_CONTROLLER_MANAGER_AUTO_ACCEPT_SERVERS` (`false`): automatically accept discovered servers, by default `.spec.accepted` should be changed to `true` to accept the server
* `SIDERO_CONTROLLER_MANAGER_INSECURE_WIPE` (`true`): wipe only the first megabyte of each disk on the server, otherwise wipe the full disk
* `SIDERO_CONTROLLER_MANAGER_SERVER_REBOOT_TIMEOUT` (`20m`): timeout for the server reboot (how long it might take for the server to be rebooted before Sidero retries an IPMI reboot operation)
* `SIDERO_METADATA_SERVER_HOST_NETWORK` (`false`): run `sidero-metadta-server` on host network
* `SIDERO_METADATA_SERVER_PORT` (`8080`): port to use for the metadata server
