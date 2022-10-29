---
description: "Prerequisite: DHCP Service"
weight: 4
title: "Prerequisite: DHCP service"
---

In order to network boot Talos, we need to set up our DHCP server to supply the
network boot parameters to our servers.
For maximum flexibility, Sidero makes use of iPXE to be able to reference
artifacts via HTTP.
Some modern servers support direct UEFI HTTP boot, but most existing servers
still rely on the old, slow TFTP-based PXE boot first.
Therefore, we need to tell our DHCP server to find the iPXE binary on a TFTP
server.

Conveniently, Sidero comes with a TFTP server which will serve the appropriate
files.
We need only set up our DHCP server to point to it.

The tricky bit is that at different phases, we need to serve different assets,
but they all use the same DHCP metadata key.

In fact, for each architecture, we have as many as four different client types:

- Legacy BIOS-based PXE boot (undionly.kpxe via TFTP)
- UEFI-based PXE boot (ipxe.efi via TFTP)
- UEFI HTTP boot (ipxe.efi via HTTP URL)
- iPXE (boot.ipxe via HTTP URL)

## Common client types

If you are lucky and all of the machines in a given DHCP zone can use the same
network boot client mechanism, your DHCP server only needs to provide two
options:

- `Server-Name` (option 66) with the IP of the Sidero TFTP service
- `Bootfile-Name` (option 67) with the appropriate value for the boot client type:
  - Legacy BIOS PXE boot: `undionly.kpxe`
  - UEFI-based PXE boot: `ipxe.efi`
  - UEFI HTTP boot: `http://sidero-server-url/tftp/ipxe.efi`
  - iPXE boot: `http://sidero-server-url/boot.ipxe`

In the ISC DHCP server, these options look like:

```text
next-server 172.16.199.50;
filename "ipxe.efi";
```

## Multiple client types

Any given server will usually use only one of those, but if you have a mix of
machines, you may need a combination of them.
In this case, you would need a way to provide different images for different
client or machine types.

Both ISC DHCP server and dnsmasq provide ways to supply such conditional responses.
In this tutorial, we are working with ISC DHCP.

For modularity, we are breaking the conditional statements into a separate file
and using the `include` statement to load them into the main `dhcpd.conf` file.

In our example below, `172.16.199.50` is the IP address of our Sidero service.

`ipxe-metal.conf`:

```text
allow bootp;
allow booting;

# IP address for PXE-based TFTP methods
next-server 172.16.199.50;

# Configuration for iPXE clients
class "ipxeclient" {
  match if exists user-class and (option user-class = "iPXE");
  filename "http://172.16.199.50/boot.ipxe";
}

# Configuration for legacy BIOS-based PXE boot
class "biosclients" {
  match if not exists user-class and substring (option vendor-class-identifier, 15, 5) = "00000";
  filename "undionly.kpxe";
}

# Configuration for UEFI-based PXE boot
class "pxeclients" {
  match if not exists user-class and substring (option vendor-class-identifier, 0, 9) = "PXEClient";
  filename "ipxe.efi";
}

# Configuration for UEFI-based HTTP boot
class "httpclients" {
  match if not exists user-class and substring (option vendor-class-identifier, 0, 10) = "HTTPClient";
  option vendor-class-identifier "HTTPClient";
  filename "http://172.16.199.50/tftp/ipxe.efi";
}
```

Once this file is created, we can include it from our main `dhcpd.conf` inside a
`subnet` section.

```text
shared-network sidero {
  subnet 172.16.199.0 netmask 255.255.255.0 {
    option domain-name-servers 8.8.8.8, 1.1.1.1;
    option routers 172.16.199.1;
    include "/config/ipxe-metal.conf";
  }
}
```

Since we use a number of Ubiquiti EdgeRouter devices especially in our home test
networks, it is worth mentioning the curious syntax gymnastics we must go
through there.
Essentially, the quotes around the path need to be entered as HTML entities:
`&quot;`.

Ubiquiti EdgeRouter configuration statement:

```text
set service dhcp-server shared-network-name sidero \
  subnet 172.16.199.1 \
  subnet-parameters "include &quot;/config/ipxe-metal.conf&quot;;"
```

Also note the fact that there are two semicolons at the end of the line.
The first is part of the HTML-encoded **"** (`&quot;`) and the second is the actual terminating semicolon.

## Troubleshooting

Getting the netboot environment is tricky and debugging it is difficult.
Once running, it will generally stay running;
the problem is nearly always one of a missing or incorrect configuration, since
the process involves several different components.

We are working toward integrating as much as possible into Sidero, to provide as
much intelligence and automation as can be had, but until then, you will likely
need to figure out how to begin hunting down problems.

See the Sidero [Troubleshooting](../troubleshooting) guide for more assistance.
