---
description: "A guide for bootstrapping Sidero management plane"
weight: 1
---

# Local Configuration

## Create the Default Environment

We must now create an `Environment` in our bootstrap cluster.
An environment is a CRD that tells the PXE component of Sidero what information to return to nodes that request a PXE boot after completing the registration process above.
Things that can be controlled here are kernel flags and the kernel and init images to use.

To create a default environment that will use the latest published Talos release, issue the following:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: metal.sidero.dev/v1alpha1
kind: Environment
metadata:
  name: default
spec:
  kernel:
    url: "https://github.com/talos-systems/talos/releases/latest/download/vmlinuz-amd64"
    sha512: ""
    args:
      - initrd=initramfs.xz
      - page_poison=1
      - slab_nomerge
      - slub_debug=P
      - pti=on
      - random.trust_cpu=on
      - ima_template=ima-ng
      - ima_appraise=fix
      - ima_hash=sha512
      - console=tty0
      - console=ttyS1,115200n8
      - earlyprintk=ttyS1,115200n8
      - panic=0
      - printk.devkmsg=on
      - talos.platform=metal
      - talos.config=http://$PUBLIC_IP:9091/configdata?uuid=
  initrd:
    url: "https://github.com/talos-systems/talos/releases/latest/download/initramfs-amd64.xz"
    sha512: ""
EOF
```

## Create Server Class

We must now create a server class to wrap our servers we registered.
This is necessary for using the Talos control plane provider for Cluster API.
The qualifiers needed for your server class will differ based on the data provided by your registration flow.
See the [server class docs](/docs/v0.1/configuration/serverclasses) for more info on how these work.

Here is an example of how to apply the server class once you have the proper info:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: metal.sidero.dev/v1alpha1
kind: ServerClass
metadata:
  name: default
spec:
  qualifiers:
    cpu:
      - manufacturer: Intel(R) Corporation
        version: Intel(R) Atom(TM) CPU C3558 @ 2.20GHz
EOF
```

In order to fetch hardware information, you can use

```bash
kubectl get server -o yaml
```

Note that for bare-metal setup, you would need to specify an installation disk.
See the [Installation
Disk](/docs/v0.1/configuration/servers/#installation-disk).

Once created, you should see the servers that make up your server class appear
as "available":

```bash
$ kubectl get serverclass
NAME      AVAILABLE                                  IN USE
default   ["00000000-0000-0000-0000-d05099d33360"]   []
```

## Create Management Plane

We are now ready to template out our management plane.
Using clusterctl, we can create a cluster manifest with:

```bash
clusterctl config cluster management-plane -i sidero > management-plane.yaml
```

Note that there are several variables that should be set in order for the templating to work properly:

- `CONTROL_PLANE_ENDPOINT`: The endpoint used for the Kubernetes API server (e.g. `https://1.2.3.4:6443`).
  This is the equivalent of the `endpoint` you would specify in `talosctl gen config`.
  There are a variety of ways to configure a control plane endpoint.
  Some common ways for an HA setup are to use DNS, a load balancer, or BGP.
  A simpler method is to use the IP of a single node.
  This has the disadvantage of being a single point of failure, but it can be a simple way to get running.
- `CONTROL_PLANE_SERVERCLASS`: The server class to use for control plane nodes.
- `WORKER_SERVERCLASS`: The server class to use for worker nodes.
- `KUBERNETES_VERSION`: The version of Kubernetes to deploy (e.g. `v1.19.4`).
- `CONTROL_PLANE_PORT`: The port used for the Kubernetes API server (port 6443)

For instance:

```bash
export CONTROL_PLANE_SERVERCLASS=master
export WORKER_SERVERCLASS=worker
export KUBERNETES_VERSION=v1.20.1
export CONTROL_PLANE_PORT=6443
export CONTROL_PLANE_ENDPOINT=1.2.3.4
clusterctl config cluster management-plane -i sidero > management-plane.yaml
```

In addition, you can specify the replicas for control-plane & worker nodes in
management-plane.yaml manifest for TalosControlPlane and MachineDeployment
objects.
Also, they can be scaled if needed:

```bash
kubectl get taloscontrolplane
kubectl get machinedeployment
kubectl scale taloscontrolplane management-plane-cp --replicas=3
```

Now that we have the manifest, we can simply apply it:

```bash
kubectl apply -f management-plane.yaml
```

**NOTE: The templated manifest above is meant to act as a starting point.
If customizations are needed to ensure proper setup of your Talos cluster, they should be added before applying.**

Once the management plane is setup, you can fetch the talosconfig by using the cluster label.
Be sure to update the cluster name and issue the following command:

```bash
kubectl get talosconfig \
  -l cluster.x-k8s.io/cluster-name=<CLUSTER NAME> \
  -o yaml -o jsonpath='{.items[0].status.talosConfig}' > management-plane-talosconfig.yaml
```

With the talosconfig in hand, the management plane's kubeconfig can be fetched with `talosctl --talosconfig management-plane-talosconfig.yaml kubeconfig`

## Pivoting

Once we have the kubeconfig for the management cluster, we now have the ability to pivot the cluster from our bootstrap.
Using clusterctl, issue:

```bash
clusterctl init --kubeconfig=/path/to/management-plane/kubeconfig -i sidero -b talos -c talos
```

Followed by:

```bash
clusterctl move --to-kubeconfig=/path/to/management-plane/kubeconfig
```

Upon completion of this command, we can now tear down our bootstrap cluster with `talosctl cluster destroy` and begin using our management plane as our point of creation for all future clusters!
:q!
