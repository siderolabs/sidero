---
description: "A guide for bootstrapping Sidero management plane"
weight: 9
title: "Scale the Workload Cluster"
---

If you have more machines available, you can scale both the controlplane
(`TalosControlPlane`) and the workers (`MachineDeployment`) for any cluster
after it has been deployed.
This is done just like normal Kubernetes `Deployments`.

```bash
kubectl scale taloscontrolplane cluster-0-cp --replicas=3
```
