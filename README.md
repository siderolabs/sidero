# Sidero

Kubernetes Bare Metal Lifecycle Management.
Sidero Metal provides lightweight, composable tools that can be used to create bare-metal Talos + Kubernetes clusters.
 Sidero Metal is an open source project from [Sidero Labs](https://www.SideroLabs.com).

## Documentation

Visit the project [site](https://www.sidero.dev).

## Compatibility with Cluster API and Kubernetes Versions

This provider's versions are compatible with the following versions of Cluster API:

|                              | v1alpha3 (v0.3) | v1alpha4 (v0.4) | v1beta1 (v1.0) |
| ---------------------------- | --------------- | --------------- | -------------- |
|                        | v1alpha3 (v0.3) | v1alpha4 (v0.4) | v1beta1 (v1.x) |
| ---------------------- | --------------- | --------------- | -------------- |
| Sidero Provider (v0.3) | ✓               |                 |                |
| Sidero Provider (v0.4) |                 | ✓               |                |
| Sidero Provider (v0.5) |                 |                 | ✓              |

This provider's versions are able to install and manage the following versions of Kubernetes:

|                        | v1.16 | v 1.17 | v1.18 | v1.19 | v1.20 | v1.21 | v1.22 | v1.23 | v1.24 |
| ---------------------- | ----- | ------ | ----- | ----- | ----- | ----- | ----- | ----- | ----- |
| Sidero Provider (v0.3) | ✓     | ✓      | ✓     | ✓     | ✓     | ✓     |       |       |       |
| Sidero Provider (v0.4) |       |        |       | ✓     | ✓     | ✓     | ✓     | ✓     |       |
| Sidero Provider (v0.5) |       |        |       | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     |

This provider's versions are compatible with the following versions of Talos:

|                        | v0.11 | v0.12  | v0.13 | v0.14 | v1.0  | v1.1  |
| ---------------------- | ----- | ------ | ----- | ----- | ----- | ----- |
| Sidero Provider (v0.3) | ✓     | ✓      |       |       |       |       |
| Sidero Provider (v0.4) | ✓     | ✓      | ✓     | ✓     |       |       |
| Sidero Provider (v0.5) |       | ✓ (+)  | ✓ (+) | ✓     | ✓     | ✓     |

> (+): Some Sidero 0.5 features (SideroLink) are only available with Talos v0.14+.

## Support

Join our [Slack](https://slack.dev.talos-systems.io)!
