---
description: ""
sidebar: "docs"
prev: "/docs/installation/"
next: "/docs/patching/"
---

# Bootstrapping

## Create a Local Cluster

```bash
talosctl cluster create \
  -p 69:69/udp,8081:8081/tcp,9091:9091/tcp,50100:50100/tcp \
  --workers 0 \
  --endpoint $PUBLIC_IP
```

## Untaint the Control Plane Node

```bash
kubectl taint node talos-default-master-1 node-role.kubernetes.io/master:NoSchedule-
```

## Patch the Metadata Server

```bash
## Update args to use 9091 for port
kubectl patch deploy -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args", "value": ["--port=9091"]}]'

## Tweak container port to match
kubectl patch deploy -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/ports", "value": [{"containerPort": 9091,"name": "http"}]}]'

## Use host networking
kubectl patch deploy -n sidero-system sidero-metadata-server --type='json' -p='[{"op": "add", "path": "/spec/template/spec/hostNetwork", "value": true}]'
```

## Patch the Metal Controller Manager

```bash
## Update args to specify the api endpoint to use for registration
kubectl patch deploy -n sidero-system sidero-controller-manager --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/1/args", "value": ["--api-endpoint=192.168.1.203","--metrics-addr=127.0.0.1:8080","--enable-leader-election"]}]'

## Use host networking
kubectl patch deploy -n sidero-system sidero-controller-manager --type='json' -p='[{"op": "add", "path": "/spec/template/spec/hostNetwork", "value": true}]'
```

## Create the Default Environment

```bash
cat <<EOF | kubectl apply -f -
apiVersion: metal.sidero.dev/v1alpha1
kind: Environment
metadata:
  name: default
spec:
  kernel:
    url: "https://github.com/talos-systems/talos/releases/latest/download/vmlinuz"
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
    url: "https://github.com/talos-systems/talos/releases/latest/download/initramfs.xz"
    sha512: ""
EOF
```

## Register the Servers

At this point, any servers on the same networkd as Sidero should PXE boot using the Sidero PXE service.
To register a server with Sidero, simply turn it on and Sidero will do the rest.

- Create server class at least for control plane nodes

```yaml
apiVersion: cluster.x-k8s.io/v1alpha3
kind: Cluster
metadata:
  name: management-plane
  namespace: test
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
        - 10.97.0.0/16
    services:
      cidrBlocks:
        - 10.98.0.0/16
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    kind: MetalCluster
    name: management-plane
    namespace: test
  controlPlaneRef:
    kind: TalosControlPlane
    apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
    name: management-plane-cp
    namespace: test
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: MetalCluster
metadata:
  name: management-plane
  namespace: test
spec:
  controlPlaneEndpoint:
    host: management-plane.rsmitty.cloud
    port: 6443
---
apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
kind: MetalMachineTemplate
metadata:
  name: management-plane-cp
  namespace: test
spec:
  template:
    spec:
      serverClassRef:
        apiVersion: metal.sidero.dev/v1alpha1
        kind: ServerClass
---
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
kind: TalosControlPlane
metadata:
  name: management-plane-cp
  namespace: test
spec:
  version: v1.18.1
  replicas: 1
  infrastructureTemplate:
    kind: MetalMachineTemplate
    apiVersion: infrastructure.cluster.x-k8s.io/v1alpha3
    name: management-plane-cp
    namespace: test
  controlPlaneConfig:
    init:
      generateType: init
    controlplane:
      generateType: controlplane
```
