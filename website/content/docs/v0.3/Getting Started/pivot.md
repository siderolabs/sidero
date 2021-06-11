---
description: "A guide for bootstrapping Sidero management plane"
weight: 11
---

# Optional: Pivot management cluster

Having the Sidero cluster running inside a Docker container is not the most
robust place for it, but it did make for an expedient start.

Conveniently, you can create a Kubernetes cluster in Sidero and then _pivot_ the
management plane over to it.

Start by creating a workload cluster as you have already done.
In this example, this new cluster is called `management`.

After the new cluster is available, install Sidero onto it as we did before,
making sure to set all the environment variables or configuration parameters for
the _new_ management cluster first.

```bash
export SIDERO_CONTROLLER_MANAGER_API_ENDPOINT=sidero.mydomain.com

clusterctl init \
  --kubeconfig-context=management
  -i sidero -b talos -c talos
```

Now, you can move the database from `sidero-demo` to `management`:

```bash
clusterctl move \
  --kubeconfig-context=sidero-demo \
  --to-kubeconfig-context=management
```

## Delete the old Docker Management Cluster

If you created your `sidero-demo` cluster using Docker as described in this
tutorial, you can now remove it:

```bash
talosctl cluster destroy --name sidero-demo
```
