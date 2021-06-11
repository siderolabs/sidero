---
description: "A guide for bootstrapping Sidero management plane"
weight: 6
---

# Expose Sidero Services

> If you built your cluster as specified in the [Prerequisite: Kubernetes] section in this tutorial, your services are already exposed and you can skip this section.

There are two external Services which Sidero serves and which much be made
reachable by the servers which it will be driving.

For most servers, TFTP (port 69/udp) will be needed.
This is used for PXE booting, both BIOS and UEFI.
Being a primitive UDP protocl, many load balancers do not support TFTP.
Instead, solutions such as [MetalLB](https://metallb.universe.tf) may be used to expose TFTP over a known IP address.
For servers which support UEFI HTTP Network Boot, TFTP need not be used.

The kernel, initrd, and all configuration assets are served from the HTTP service
(port 8081/tcp).
It is needed for all servers, but since it is HTTP-based, it
can be easily proxied, load balanced, or run through an ingress controller.

The main thing to keep in mind is that the services **MUST** match the IP or
hostname specified by the `SIDERO_CONTROLLER_MANAGER_API_ENDPOINT` environment
variable (or configuration parameter) when you installed Sidero.

It is a good idea to verify that the services are exposed as you think they
should be.

```bash
$ curl -I http://192.168.1.150:8081/tftp/ipxe.efi
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 1020416
Content-Type: application/octet-stream
```
