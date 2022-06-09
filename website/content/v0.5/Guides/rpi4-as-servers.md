---
description: "Using Raspberrypi Pi 4 as servers"
weight: 6
title: "Raspberry Pi4 as Servers"
---

This guide will explain on how to use Sidero to manage Raspberrypi-4's as
servers.
This guide goes  hand in hand with the [bootstrapping
guide](../../guides/bootstrapping).

From the bootstrapping guide, reach "Install Sidero" and come back to this
guide.
Once you finish with this guide, you will need to go back to the
bootstrapping guide and continue with "Register the servers".

The rest of this guide goes with the assumption that you've a cluster setup with
Sidero and ready to accept servers.
This guide will explain the changes that needs to be made to be able to accept RPI4 as server.

## RPI4 boot process

To be able to boot talos on the Pi4 via network, we need to undergo a 2-step boot process.
The Pi4 has an EEPROM which contains code to boot up the Pi.
This  EEPROM expects a specific boot folder structure as explained on
[this](https://www.raspberrypi.org/documentation/configuration/boot_folder.md) page.
We will use the EEPROM to boot into UEFI, which we will then use to PXE and iPXE boot into sidero & talos.

## Prerequisites

### Update EEPROM

_NOTE:_ If you've updated the EEPROM with the image that was referenced on [the talos docs](https://www.talos.dev/latest/talos-guides/install/single-board-computers/rpi_4/#updating-the-eeprom),
you can either flash it with the one mentioned below, or visit [the EEPROM config docs](https://www.raspberrypi.org/documentation/hardware/raspberrypi/bcm2711_bootloader_config.md)
and change the boot order of EEPROM to `0xf21`.
Which means try booting from SD first, then try network.

To enable the EEPROM on the Pi to support network booting, we must update it to
the latest version.
Visit the [release](https://github.com/raspberrypi/rpi-eeprom/releases) page and grab the
latest `rpi-boot-eeprom-recovery-*-network.zip` (as of time of writing,
v2021.0v.29-138a1 was used).
Put this on a SD card and plug it into the Pi.
The
Pi's status light will flash rapidly after a few seconds, this indicates that
the EEPROM has been updated.

This operation needs to be done once per Pi.

### Serial number

Power on the Pi without an SD card in it and hook it up to a monitor, you will
be greeted with the boot screen.
On this screen you will find some information
about the Pi.
For this guide, we are only interested in the serial number.
The
first line under the Pi logo will be something like the following:

`board: xxxxxx <serial> <MAC address>`

Write down the 8 character serial.

### talos-systems/pkg

Clone the [talos-systems/pkg](https://github.com/talos-systems/pkgs) repo.
Create a new folder called `raspberrypi4-uefi` and `raspberrypi4-uefi/serials`.
Create a file `raspberrypi4-uefi/pkg.yaml` containing the following:

```yaml
name: raspberrypi4-uefi
variant: alpine
install:
  - unzip
steps:
# {{ if eq .ARCH "aarch64" }} This in fact is YAML comment, but Go templating instruction is evaluated by bldr restricting build to arm64 only
  - sources:
      - url: https://github.com/pftf/RPi4/releases/download/v1.26/RPi4_UEFI_Firmware_v1.26.zip # <-- update version NR accordingly.
        destination: RPi4_UEFI_Firmware.zip
        sha256: d6db87484dd98dfbeb64eef203944623130cec8cb71e553eab21f8917e0285f7
        sha512: 96a71086cdd062b51ef94726ebcbf15482b70c56262555a915499bafc04aff959d122410af37214760eda8534b58232a64f6a8a0a8bb99aba6de0f94c739fe98
    prepare:
      - |
        unzip RPi4_UEFI_Firmware.zip
        rm RPi4_UEFI_Firmware.zip
        mkdir /rpi4
        mv ./* /rpi4
    install:
      - |
        mkdir /tftp
        ls /pkg/serials | while read serial; do mkdir /tftp/$serial && cp -r /rpi4/* /tftp/$serial && cp -r /pkg/serials/$serial/* /tftp/$serial/; done
# {{ else }}
  - install:
      - |
        mkdir -p /tftp
# {{ end }}
finalize:
  - from: /
    to: /
```

## UEFI / RPi4

Now that the EEPROM can network boot, we need to prepare the structure of our
boot folder.
Essentially what the bootloader will do is look for this folder
on the network rather than on the SD card.

Visit the [release page of RPi4](https://github.com/pftf/RPi4/releases) and grab
the latest `RPi4_UEFI_Firmware_v*.zip` (at the time of writing, v1.26 was used).
Extract the zip into a folder, the structure will look like the following:

```bash
.
├── RPI_EFI.fd
├── RPi4_UEFI_Firmware_v1.26.zip
├── Readme.md
├── bcm2711-rpi-4-b.dtb
├── bcm2711-rpi-400.dtb
├── bcm2711-rpi-cm4.dtb
├── config.txt
├── firmware
│   ├── LICENCE.txt
│   ├── Readme.txt
│   ├── brcmfmac43455-sdio.bin
│   ├── brcmfmac43455-sdio.clm_blob
│   └── brcmfmac43455-sdio.txt
├── fixup4.dat
├── overlays
│   └── miniuart-bt.dtbo
└── start4.elf
```

As a one time operation, we need to configure UEFI to do network booting by
default, remove the 3gb mem limit if it's set and optionally set the CPU clock to
max.
Take these files and put them on the SD card and boot the Pi.
You will see the Pi logo, and the option to hit `esc`.

### Remove 3GB mem limit

1. From the home page, visit "Device Manager".
2. Go down to "Raspberry Pi Configuration" and open that menu.
3. Go to "Advanced Configuration".
4. Make sure the option "Limit RAM to 3 GB" is set to `Disabled`.

### Change CPU to Max (optionally)

1. From the home page, visit "Device Manager".
2. Go down to "Raspberry Pi Configuration" and open that menu.
3. Go to "CPU Configuration".
4. Change CPU clock to `Max`.

## Change boot order

1. From the home page, visit "Boot Maintenance Manager".
2. Go to "Boot Options".
3. Go to "Change Boot Order".
4. Make sure that `UEFI PXEv4` is the first boot option.

### Persisting changes

Now that we have made the changes above, we need to persist these changes.
Go back to the home screen and hit `reset` to save the changes to disk.

When you hit `reset`, the settings will be saved to the `RPI_EFI.fd` file on the
SD card.
This is where we will run into a limitation that is explained in the
following issue: [pftf/RPi4#59](https://github.com/pftf/RPi4/issues/59).
What this mean is that we need to create a `RPI_EFI.fd` file for each Pi that we want to use as server.
This is because the MAC address is also stored in the `RPI_EFI.fd` file,
which makes it invalid when you try to use it in a different Pi.

Plug the SD card back into your computer and extract the `RPI_EFI.fd` file from
it and place it into the `raspberrypi4-uefi/serials/<serial>/`.
The dir should look like this:

```bash
raspberrypi4-uefi/
├── pkg.yaml
└── serials
    └─── XXXXXXXX
        └── RPI_EFI.fd
```

## Build the image with the boot folder contents

Now that we have the `RPI_EFI.fd` of our Pi in the correct location, we must now
build a docker image containing the boot folder for the EEPROM.
To do this, run the following command in the pkgs repo:

`make PLATFORM=linux/arm64 USERNAME=$USERNAME PUSH=true TARGETS=raspberrypi4-uefi`

This will build and push the following image:
`ghcr.io/$USERNAME/raspberrypi4-uefi:<tag>`

_If you need to change some other settings like registry etc, have a look in the
Makefile to see the available variables that you can override._

The content of the `/tftp` folder in the image will be the following:

```bash
XXXXXXXX
├── RPI_EFI.fd
├── Readme.md
├── bcm2711-rpi-4-b.dtb
├── bcm2711-rpi-400.dtb
├── bcm2711-rpi-cm4.dtb
├── config.txt
├── firmware
│   ├── LICENCE.txt
│   ├── Readme.txt
│   ├── brcmfmac43455-sdio.bin
│   ├── brcmfmac43455-sdio.clm_blob
│   └── brcmfmac43455-sdio.txt
├── fixup4.dat
├── overlays
│   └── miniuart-bt.dtbo
└── start4.elf
```

## Patch metal controller

To enable the 2 boot process, we need to include this EEPROM boot folder into
the sidero's tftp folder.
To achieve this, we will use an init container using
the image we created above to copy the contents of it into the tftp folder.

Create a file `patch.yaml` with the following contents:

```yaml
spec:
  template:
    spec:
      volumes:
        - name: tftp-folder
          emptyDir: {}
      initContainers:
      - image: ghcr.io/<USER>/raspberrypi4-uefi:v<TAG> # <-- change accordingly.
        imagePullPolicy: Always
        name: tftp-folder-setup
        command:
          - cp
        args:
          - -r
          - /tftp
          - /var/lib/sidero/
        volumeMounts:
          - mountPath: /var/lib/sidero/tftp
            name: tftp-folder
      containers:
      - name: manager
        volumeMounts:
          - mountPath: /var/lib/sidero/tftp
            name: tftp-folder
```

Followed by this command to apply the patch:

```bash
kubectl -n sidero-system patch deployments.apps sidero-controller-manager --patch "$(cat patch.yaml)"
```

## Configure BootFromDiskMethod

By default, Sidero will use iPXE's `exit` command to attempt to force boot from disk.
On Raspberry Pi, this will drop you into the bootloader interface, and you will need to connect a keyboard and manually select the disk to boot from.

The BootFromDiskMethod can be configured on individual [Servers](../../resource-configuration/servers/#bootfromdiskmethod), on [ServerClasses](../../resource-configuration/serverclasses/#bootfromdiskmethod), or as a command-line argument to the Sidero metal controller itself (`--boot-from-disk-method=<value>`).
In order to force the Pi to use the configured bootloader order, the BootFromDiskMethod needs to be set to `ipxe-sanboot`.

## Profit

With the patched metal controller, you should now be able to register the Pi4 to
sidero by just connecting it to the network.
From this point you can continue with the [bootstrapping guide](../../guides/bootstrapping#register-the-servers).
