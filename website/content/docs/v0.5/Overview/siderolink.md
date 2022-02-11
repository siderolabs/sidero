---
description: ""
weight: 4
title: SideroLink
---

SideroLink provides an overlay Wireguard point-to-point connection from every Talos machine to the Sidero.
Sidero provisions each machine with a unique IPv6 address and Wireguard key for the SideroLink connection.

> Note: SideroLink is only supported with Talos >= 0.14.
>
> SideroLink doesn't provide a way for workload machines to communicate with each other, a connection is only
> point-to-point.

SideroLink connection is both encrypted and authenticated, so Sidero uses that to map data streams coming from the machines
to a specific `ServerBinding`, `MetalMachine`, `Machine` and `Cluster`.

Talos node sends two streams over the SideroLink connection: kernel logs (dmesg) and Talos event stream.
SideroLink is enabled automatically by Sidero when booting Talos.

## Kernel Logs

Kernel logs (`dmesg`) are streamed in real time from the Talos nodes to the `sidero-controller-manager` over SideroLink connection.
Log streaming starts when the kernel passes control to the `init` process, so kernel boot time logs will only be available when control
is passed to the userland.

Logs can be accessed by accessing the logs of the `serverlogs` container of the `sidero-controller-manager` pod:

```bash
$ kubectl -n sidero-system logs deployment/sidero-controller-manager -c serverlogs -f
{"clock":8576583,"cluster":"management-cluster","facility":"user","machine":"management-cluster-cp-ddgsw","metal_machine":"management-cluster-cp-vrff4","msg":"[talos] phase mountState (6/13): 1 tasks(s)\n","namespace":"default","priority":"warning","seq":665,"server_uuid":"6b121f82-24a8-4611-9d23-fa1a5ba564f0","talos-level":"warn","talos-time":"2022-02-11T12:42:02.74807823Z"}
...
```

The format of the message is the following:

```json
{
    "clock": 8576583,
    "cluster": "management-cluster",
    "facility": "user",
    "machine": "management-cluster-cp-ddgsw",
    "metal_machine": "management-cluster-cp-vrff4",
    "msg": "[talos] phase mountState (6/13): 1 tasks(s)\n",
    "namespace": "default",
    "priority": "warning",
    "seq": 665,
    "server_uuid": "6b121f82-24a8-4611-9d23-fa1a5ba564f0",
    "talos-level": "warn",
    "talos-time": "2022-02-11T12:42:02.74807823Z"
}
```

Kernel fields (see [Linux documentation](https://www.kernel.org/doc/Documentation/ABI/testing/dev-kmsg) for details):

- `clock` is the kernel timestamp relative to the boot time
- `facility` of the message
- `msg` is the actual log message
- `seq` is the kernel log sequence
- `priority` is the message priority

Talos-added fields:

- `talos-level` is the translated `priority` into standard logging levels
- `talos-time` is the timestamp of the log message (accuracy of the timestamp depends on time sync)

Sidero-added fields:

- `server_uuid` is the `name` of the matching `Server` and `ServerBinding` resources
- `namespace` is the namespace of the `Cluster`, `MetalMachine` and `Machine`
- `cluster`, `metal_machine` and `machine` are the names of the matching `Cluster`, `MetalMachine` and `Machine` resources

It might be a good idea to send container logs to some log aggregation system and filter the logs for a cluster or a machine.

Quick filtering for a specific server:

```bash
kubectl -n sidero-system logs deployment/sidero-controller-manager -c serverlogs  | jq -R 'fromjson? | select(.server_uuid == "b4e677d9-b59b-4c1c-925a-f9d9ce049d79")'
```

## Talos Events

Talos delivers system events over the SideroLink connection to the `sidero-link-manager` pod.
These events can be accessed with `talosctl events` command.
Events are mostly used to update `ServerBinding`/`MetalMachine` statuses, but they can be also seen in the logs of the `serverevents` container:

```bash
$ kubectl -n sidero-system logs deployment/sidero-controller-manager -c serverevents -f
{"level":"info","ts":1644853714.2700942,"caller":"events-manager/adapter.go:153","msg":"incoming event","component":"sink","node":"[fdae:2859:5bb1:7a03:3ae3:be30:7ec4:4c09]:44530","id":"c857jkm1jjcc7393cbs0","type":"type.googleapis.com/machine.
AddressEvent","server_uuid":"b4e677d9-b59b-4c1c-925a-f9d9ce049d79","cluster":"management-cluster","namespace":"default","metal_machine":"management-cluster-cp-47lll","machine":"management-cluster-cp-7mpsh","hostname":"pxe-2","addresses":"172.25.0.5"}
```

## MetalMachine Conditions

Sidero updates the statuses of `ServerBinding`/`MetalMachine` resources based on the events received from Talos node:

- current addresses of the node
- statuses of machine configuration loading and validation, installation status

See [Resources](../resources/) for details.

## SideroLink State

State of the SideroLink connection is kept in the `ServerBinding` resource:

```yaml
spec:
  siderolink:
    address: fdae:2859:5bb1:7a03:3ae3:be30:7ec4:4c09/64
    publicKey: XIBT49g9xCoBvyb/x36J+ASlQ4qaxXMG20ZgKbBbfE8=
```

Installation-wide SideroLink state is kept in the `siderolink` `Secret` resource:

```bash
$ kubectl get secrets siderolink -o yaml
apiVersion: v1
data:
  installation-id: QUtmZGFmVGJtUGVFcWp0RGMzT1BHSzlGcmlHTzdDQ0JCSU9aRzRSamdtWT0=
  private-key: ME05bHhBd3JwV0hDczhNbm1aR3RDL1ZjK0ZSUFM5UzQwd25IU00wQ3dHOD0=
...
```

Key `installation-id` is used to generate unique SideroLink IPv6 addresses, and `private-key` is the Wireguard key of Sidero.
