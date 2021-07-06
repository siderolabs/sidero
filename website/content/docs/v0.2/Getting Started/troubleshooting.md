---
description: "Troubleshooting"
weight: 99
title: "Troubleshooting"
---

The first thing to do in troubleshooting problems with the Sidero installation
and operation is to figure out _where_ in the process that failure is occurring.

Keep in mind the general flow of the pieces.
For instance:

1. A server is configured by its BIOS/CMOS to attempt a network boot using the PXE firmware on
its network card(s).
1. That firmware requests network and PXE boot configuration via DHCP.
1. DHCP points the firmware to the Sidero TFTP or HTTP server (depending on the firmware type).
1. The second stage boot, iPXE, is loaded and makes an HTTP request to the
    Sidero metadata server for its configuration, which contains the URLs for
    the kernel and initrd images.
1. The kernel and initrd images are downloaded by iPXE and boot into the Sidero
    agent software (if the machine is not yet known and assigned by Sidero).
1. The agent software reports to the Sidero metadata server via HTTP the hardware information of the machine.
1. A (usually human or external API) operator verifies and accepts the new
    machine into Sidero.
1. The agent software reboots and wipes the newly-accepted machine, then powers
    off the machine to wait for allocation into a cluster.
1. The machine is allocated by Sidero into a Kubernetes Cluster.
1. Sidero tells the machine, via IPMI, to boot into the OS installer
     (following all the same network boot steps above).
1. The machine downloads its configuration from the Sidero metadata server via
     HTTP.
1. The machine applies its configuration, installs a bootloader, and reboots.
1. The machine, upon reboot from its local disk, joins the Kubernetes cluster
     and continues until Sidero tells it to leave the cluster.
1. Sidero tells the machine to leave the cluster and reboots it into network
     boot mode, via IPMI.
1. The machine netboots into wipe mode, wherein its disks are again wiped to
     come back to the "clean" state.
1. The machine again shuts down and waits to be needed.

## Device firmware (PXE boot)

The worst place to fail is also, unfortunately, the most common.
This is the firmware phase, where the network card's built-in firmware attempts
to initiate the PXE boot process.
This is the worst place because the firmware is completely opaque, with very
little logging, and what logging _does_ appear frequently is wiped from the
console faster than you can read it.

If you fail here, the problem will most likely be with your DHCP configuration,
though it _could_ also be in the Sidero TFTP service configuration.

## Validate Sidero TFTP service

The easiest to validate is to use a `tftp` client to validate that the Sidero
TFTP service is available at the IP you are advertising via DHCP.

```bash
  $ atftp 172.16.199.50
  tftp> get ipxe.efi
```

TFTP is an old, slow protocol with very little feedback or checking.
Your only real way of telling if this fails is by timeout.
Over a local network, this `get` command should take a few seconds.
If it takes longer than 30 seconds, it is probably not working.

Success is also not usually indicated:
you just get a prompt returned, and the file should show up in your current
directory.

If you are failing to connect to TFTP, the problem is most likely with your
Sidero Service exposure:
how are you exposing the TFTP service in your management cluster to the outside
world?
This normally involves either setting host networking on the Deployment or
installing and using something like MetalLB.
