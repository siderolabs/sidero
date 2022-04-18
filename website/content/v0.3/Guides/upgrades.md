---
description: "A guide describing upgrades"
title: "Upgrading"
weight: 5
---

Upgrading a running workload cluster or management plane is the same process as describe in the Talos documentation.

To upgrade the Talos OS, see [here](https://www.talos.dev/latest/talos-guides/upgrading-talos/).

In order to upgrade Kubernetes itself, see [here](https://www.talos.dev/latest/kubernetes-guides/upgrading-kubernetes/).

## Upgrading Talos 0.8 -> 0.9

It is important, however, to take special consideration for upgrades of the Talos v0.8.x series to v0.9.x.
Because of the move from self-hosted control plane to static pods, some certificate information has changed that needs to be manually updated.
The steps are as follows:

- Upgrade a single control plane node to the v0.9.x series using the upgrade instructions above.
upgrade

- After upgrade, carry out a `talosctl convert-k8s` to move from the self-hosted control plane to static pods.

- Targeting the upgraded node, issue `talosctl read -n <node-ip> /system/state/config.yaml` and copy out the `cluster.aggregatorCA` and `cluster.serviceAccount` sections.

- In the management cluster, issue `kubectl edit secret <cluster-name>-talos`.

- While in editing view, copy the `data.certs` field and decode it with `echo '<certs-content>' | base64 -d`

> Note: It may also be a good idea to copy the secret in its entirety as a backup.
> This can be done with a simple `kubectl get secret <cluster-name>-talos -o yaml`.

- Copying the output above to a text editor, update the aggregator and service account sections with the certs and keys copied previously and save it.
The resulting file should look like:

```yaml
admin:
  crt: xxx
  key: xxx
etcd:
  crt: xxx
  key: xxx
k8s:
  crt: xxx
  key: xxx
k8saggregator:
  crt: xxx
  key: xxx
k8sserviceaccount:
  key: xxx
os:
  crt: xxx
  key: xxx
```

- Re-encode the data with `cat <saved-file> | base64 | tr -d '\n'`

- With the secret still open for editing, update the `data.certs` field to contain the new base64 data.

- Edit the cluster's TalosControlPlane resource with `kubectl edit tcp <name-of-control-plane>`.
Update the `spec.controlPlaneConfig.[controlplane,init].talosVersion` fields to be `v0.9`.

- Edit any TalosConfigTemplate resources and update `spec.template.spec.talosVersion` to be the same value.

- At this point, any new controlplane or worker machines should receive the newer machine config format and join the cluster successfully.
You can also proceed to upgrade existing nodes.
