# Sidero

<!-- textlint-disable -->
> [!CAUTION]
> Sidero Labs is no longer actively developing Sidero Metal.
>
> For an alternative, please see [Omni](https://github.com/siderolabs/omni.git)
> and the [Bare-Metal Infrastructure Provider](https://omni.siderolabs.com/tutorials/setting-up-the-bare-metal-infrastructure-provider).
>
> Unless you have an existing support contract covering Sidero Metal, all support will be provided by the community (including questions in our Slack workspace).
<!-- textlint-enable -->

Kubernetes Bare Metal Lifecycle Management.
Sidero Metal provides lightweight, composable tools that can be used to create bare-metal Talos + Kubernetes clusters.
Sidero Metal is an open-source project from [Sidero Labs](https://www.SideroLabs.com).

## Documentation

Visit the project [site](https://www.sidero.dev).

## Compatibility with Cluster API and Kubernetes Versions

This provider's versions are compatible with the following versions of Cluster API:

|                        | v1alpha3 (v0.3) | v1alpha4 (v0.4) | v1beta1 (v1.x) |
| ---------------------- | --------------- | --------------- | -------------- |
| Sidero Provider (v0.5) |                 |                 | ✓              |
| Sidero Provider (v0.6) |                 |                 | ✓              |

This provider's versions are able to install and manage the following versions of Kubernetes:

|                        | v1.19 | v1.20 | v1.21 | v1.22 | v1.23 | v1.24 | v1.25 | v1.26 | v1.27 | v1.28 | v1.29 | v1.30 | v1.31 | v1.32 |
| ---------------------- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- |
| Sidero Provider (v0.5) | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     |       |       |       |       |       |
| Sidero Provider (v0.6) |       |       |       |       |       | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     |

This provider's versions are compatible with the following versions of Talos:

|                        | v0.12  | v0.13 | v0.14 | v1.0  | v1.1  | v1.2  | v1.3  | v1.4  | v1.5  | v1.6  | v1.7  | v1.8  | v1.9  |
| ---------------------- | ------ | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- |
| Sidero Provider (v0.5) | ✓ (+)  | ✓ (+) | ✓     | ✓     | ✓     | ✓     | ✓     |       |       |       |       |       |       |
| Sidero Provider (v0.6) |        |       |       | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     |

## Support

Join our [Slack](https://slack.dev.talos-systems.io)!
