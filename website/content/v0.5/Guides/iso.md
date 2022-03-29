---
description: "A guide for bootstrapping Sidero management plane using the ISO image"
weight: 1
title: "Building A Management Plane with ISO Image"
---

This guide will provide some very basic detail about how you can also build a Sidero management plane using the Talos ISO image instead of following the Docker-based process that we detail in our Getting Started tutorials.

Using the ISO is a perfectly valid way to build a Talos cluster, but this approach is not recommended for Sidero as it avoids the "pivot" step detailed [here](../../getting-started/pivot).
Skipping this step means that the management plane does not become "self-hosted", in that it cannot be upgraded and scaled using the Sidero processes we follow for workload clusters.
For folks who are willing to take care of their management plane in other ways, however, this approach will work fine.

The rough outline of this process is very short and sweet, as it relies on other documentation:

- For each management plane node, boot the ISO and install Talos using the "apply-config" process mentioned in our Talos [Getting Started](https://www.talos.dev/docs/v0.13/introduction/getting-started/) docs.
  These docs go into heavy detail on using the ISO, so they will not be recreated here.

- With a Kubernetes cluster now in hand (and with access to it via `talosctl` and `kubectl`), you can simply pickup the Getting Started tutorial at the "Install Sidero" section [here](../../getting-started/install-clusterapi).
  Keep in mind, however, that you will be unable to do the "pivoting" section of the tutorial, so just skip that step when you reach the end of the tutorial.

> Note: It may also be of interest to view the prerequisite guides on [CLI](../../getting-started/prereq-cli-tools) and [DHCP](../../getting-started/prereq-dhcp) setup, as they will still apply to this method.

- For long-term maintenance of a management plane created in this way, refer to the Talos documentation for upgrading [Kubernetes](https://www.talos.dev/docs/v0.13/guides/upgrading-kubernetes/) and [Talos](https://www.talos.dev/docs/v0.13/guides/upgrading-talos/) itself.
