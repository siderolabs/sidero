---
description: "A guide for bootstrapping Sidero management plane"
weight: 1
---

# Bootstrapping

## Introduction

Imagine a scenario in which you have shown up to a datacenter with only a laptop and your task is to transition a rack of bare metal machines into an HA management plane and multiple Kubernetes clusters created by that management plane.
In this guide, we will go through how to create a bootstrap cluster using a Docker-based Talos cluster, provision the management plane, and pivot over to it.
Guides around post-pivoting setup and subsequent cluster creation should also be found in the "Guides" section of the sidebar.

Because of the design of Cluster API, there is inherently a "chicken and egg" problem with needing a Kubernetes cluster in order to provision the management plane.
Talos Systems and the Cluster API community have created tools to help make this transition easier.

## Prerequisites

There are two big prerequisites involved with bootstrapping Sidero: routing and DHCP setup.

From the routing side, the laptop from which you are bootstrapping _must_ be accessible by the bare metal machines that we will be booting.
In the datacenter scenario described above, the easiest way to achieve this is probably to hook the laptop onto the server rack's subnet by plugging it into the top-of-rack switch.
This is needed for TFTP, PXE booting, and for the ability to register machines with the bootstrap plane.

DHCP configuration is needed to tell the metal servers what their "next server" is when PXE booting.
The configuration of this is different for each environment and each DHCP server, thus it's impossible to give an easy guide.
However, here is an example of the configuration for an Ubiquti EdgeRouter that uses vyatta-dhcpd as the DHCP service:

This block shows the subnet setup, as well as the extra "subnet-parameters" that tell the DHCP server to include the ipxe-metal.conf file.

```bash
$ show service dhcp-server shared-network-name MetalDHCP

 authoritative enable
 subnet 192.168.254.0/24 {
     default-router 192.168.254.1
     dns-server 192.168.1.200
     lease 86400
     start 192.168.254.2 {
         stop 192.168.254.252
     }
     subnet-parameters "include &quot;/etc/dhcp/ipxe-metal.conf&quot;;"
 }
```

Here is the ipxe-metal.conf file.

```bash
$ cat /etc/dhcp/ipxe-metal.conf

allow bootp;
allow booting;

next-server 192.168.1.150;
if exists user-class and option user-class = "iPXE" {
  filename "http://192.168.1.150:8081/boot.ipxe";
} else {
  filename "ipxe.efi";
}

host talos-mgmt-0 {
    fixed-address 192.168.254.2;
    hardware ethernet d0:50:99:d3:33:60;
}
```

Notice that it sets a static address for the management node that I'll be booting, in addition to providing the "next server" info.
This "next server" IP address will match references to `PUBLIC_IP` found below in this guide.

## Create a Local Cluster

The `talosctl` CLI tool has built-in support for spinning up Talos in docker containers.
Let's use this to our advantage as an easy Kubernetes cluster to start from.

Set an environment variable called `PUBLIC_IP` which is the "public" IP of your machine.
Note that "public" is a bit of a misnomer.
We're really looking for the IP of your machine, not the IP of the node on the docker bridge (ex: `192.168.1.150`).

```bash
export PUBLIC_IP="192.168.1.150"
```

We can now create our Docker cluster.
Issue the following to create a single-node cluster:

```bash
talosctl cluster create \
  -p 69:69/udp,8081:8081/tcp,9091:9091/tcp,50100:50100/tcp \
  --workers 0 \
  --endpoint $PUBLIC_IP
```

Note that there are several ports mentioned in the command above.
These allow us to access the services that will get deployed on this node.

Once the cluster create command is complete, issue `talosctl kubeconfig /desired/path` to fetch the kubeconfig for this cluster.
You should then set your `KUBECONFIG` environment variable to the path of this file.

## Untaint Control Plane

Because this is a single node cluster, we need to remove the "NoSchedule" taint on the node to make sure non-controlplane components can be scheduled.

```bash
kubectl taint node talos-default-master-1 node-role.kubernetes.io/master:NoSchedule-
```

## Install Sidero

As of Cluster API version 0.3.9, Sidero is included as a default infrastructure provider in clusterctl.

To install Sidero and the other Talos providers, simply issue:

```bash
clusterctl init -b talos -c talos -i sidero
```

## Patch Components

We will now want to ensure that the Sidero services that got created are publicly accessible across our subnet.
This will allow the metal machines to speak to these services later.

### Patch the Metadata Server

Update the metadata server component with the following patches:

```bash
## Update args to use 9091 for port
kubectl patch deploy -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args", "value": ["--port=9091"]}]'

## Tweak container port to match
kubectl patch deploy -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/ports", "value": [{"containerPort": 9091,"name": "http"}]}]'

## Use host networking
kubectl patch deploy -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "add", "path": "/spec/template/spec/hostNetwork", "value": true}]'
```

### Patch the Metal Controller Manager

```bash
## Update args to specify the api endpoint to use for registration
kubectl patch deploy -n sidero-system sidero-controller-manager --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/1/args", "value": ["--api-endpoint='$PUBLIC_IP'","--metrics-addr=127.0.0.1:8080","--enable-leader-election"]}]'

## Use host networking
kubectl patch deploy -n sidero-system sidero-controller-manager --type='json' -p='[{"op": "add", "path": "/spec/template/spec/hostNetwork", "value": true}]'
```

## Register the Servers

At this point, any servers on the same network as Sidero should PXE boot using the Sidero PXE service.
To register a server with Sidero, simply turn it on and Sidero will do the rest.
Once the registration is complete, you should see the servers registered with `kubectl get servers`:

```bash
$ kubectl get servers -o wide
NAME                                   HOSTNAME        ACCEPTED   ALLOCATED   CLEAN
00000000-0000-0000-0000-d05099d33360   192.168.254.2   false      false       false
```

## Accept the Servers

Note in the output above that the newly registered servers are not `accepted`.
In order for a server to be eligible for consideration, it _must_ be marked as `accepted`.
Before a `Server` is accepted, no write action will be performed against it.
Servers can be accepted by issuing a patch command like:

```bash
kubectl patch server 00000000-0000-0000-0000-d05099d33360 --type='json' -p='[{"op": "replace", "path": "/spec/accepted", "value": true}]'
```

For more information on server acceptance, see the [server docs](/docs/v0.1/configuration/servers).

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

Once created, you should see the servers that make up your server class appear as "available":

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

- CONTROL_PLANE_ENDPOINT
- CONTROL_PLANE_SERVERCLASS
- WORKER_SERVERCLASS
- KUBERNETES_VERSION

Now that we have the manifest, we can simply apply it:

```bash
kubectl apply -f management-plane.yaml
```

**NOTE: The templated manifest above is meant to act as a starting point. If customizations are needed to ensure proper setup of your Talos cluster, they should be added before applying.**

Once the management plane is setup, you can fetch the talosconfig with a command like:

```bash
kubectl get talosconfig -o yaml management-plane-cp-xxx -o jsonpath='{.status.talosConfig}' > management-plane-talosconfig.yaml
```

Then the management plane's kubeconfig can be fetched with `talosctl --talosconfig management-plane-talosconfig.yaml kubeconfig /desired/path`

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
