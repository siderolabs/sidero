// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package dhcp

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"github.com/insomniacslk/dhcp/iana"
	"github.com/siderolabs/gen/xslices"
)

// ServeDHCP starts the DHCP proxy server.
func ServeDHCP(logger logr.Logger, apiEndpoint string, apiPort int) error {
	server, err := server4.NewServer(
		"",
		nil,
		handlePacket(logger, apiEndpoint, apiPort),
	)
	if err != nil {
		logger.Error(err, "error on DHCP4 proxy startup")

		return err
	}

	return server.Serve()
}

func handlePacket(logger logr.Logger, apiEndpoint string, apiPort int) func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		if err := isBootDHCP(m); err != nil {
			logger.Info("ignoring packet", "source", m.ClientHWAddr, "reason", err)

			return
		}

		fwtype, err := validateDHCP(m)
		if err != nil {
			logger.Info("invalid packet", "source", m.ClientHWAddr, "reason", err)

			return
		}

		resp, err := offerDHCP(m, apiEndpoint, apiPort, fwtype)
		if err != nil {
			logger.Error(err, "failed to construct ProxyDHCP offer", "source", m.ClientHWAddr)

			return
		}

		logger.Info("offering boot response", "source", m.ClientHWAddr, "server", resp.TFTPServerName(), "boot_filename", resp.BootFileNameOption())

		_, err = conn.WriteTo(resp.ToBytes(), peer)
		if err != nil {
			logger.Error(err, "failure sending response", "source", m.ClientHWAddr)
		}
	}
}

func isBootDHCP(pkt *dhcpv4.DHCPv4) error {
	if pkt.MessageType() != dhcpv4.MessageTypeDiscover {
		return fmt.Errorf("packet is %s, not %s", pkt.MessageType(), dhcpv4.MessageTypeDiscover)
	}

	if pkt.Options[93] == nil {
		return errors.New("not a PXE boot request (missing option 93)")
	}

	return nil
}

func validateDHCP(m *dhcpv4.DHCPv4) (fwtype Firmware, err error) {
	arches := m.ClientArch()

	for _, arch := range arches {
		switch arch { //nolint:exhaustive
		case iana.INTEL_X86PC:
			fwtype = FirmwareX86PC
		case iana.EFI_IA32, iana.EFI_X86_64, iana.EFI_BC:
			fwtype = FirmwareX86EFI
		case iana.EFI_ARM64:
			fwtype = FirmwareARMEFI
		case iana.EFI_X86_HTTP, iana.EFI_X86_64_HTTP:
			fwtype = FirmwareX86HTTP
		case iana.EFI_ARM64_HTTP:
			fwtype = FirmwareARMHTTP
		}
	}

	if fwtype == FirmwareUnsupported {
		return 0, fmt.Errorf("unsupported client arch: %v", xslices.Map(arches, func(a iana.Arch) string { return a.String() }))
	}

	// Now, identify special sub-breeds of client firmware based on
	// the user-class option. Note these only change the "firmware
	// type", not the architecture we're reporting to Booters. We need
	// to identify these as part of making the internal chainloading
	// logic work properly.
	if userClasses := m.UserClass(); len(userClasses) > 0 {
		// If the client has had iPXE burned into its ROM (or is a VM
		// that uses iPXE as the PXE "ROM"), special handling is
		// needed because in this mode the client is using iPXE native
		// drivers and chainloading to a UNDI stack won't work.
		if userClasses[0] == "iPXE" && fwtype == FirmwareX86PC {
			fwtype = FirmwareX86Ipxe
		}
	}

	guid := m.GetOneOption(dhcpv4.OptionClientMachineIdentifier)
	switch len(guid) {
	case 0:
		// A missing GUID is invalid according to the spec, however
		// there are PXE ROMs in the wild that omit the GUID and still
		// expect to boot. The only thing we do with the GUID is
		// mirror it back to the client if it's there, so we might as
		// well accept these buggy ROMs.
	case 17:
		if guid[0] != 0 {
			return 0, errors.New("malformed client GUID (option 97), leading byte must be zero")
		}
	default:
		return 0, errors.New("malformed client GUID (option 97), wrong size")
	}

	return fwtype, nil
}

func offerDHCP(req *dhcpv4.DHCPv4, apiEndpoint string, apiPort int, fwtype Firmware) (*dhcpv4.DHCPv4, error) {
	serverIPs, err := net.LookupIP(apiEndpoint)
	if err != nil {
		return nil, err
	}

	if len(serverIPs) == 0 {
		return nil, fmt.Errorf("no IPs found for %s", apiEndpoint)
	}

	// pick up the first address
	serverIP := serverIPs[0]

	modifiers := []dhcpv4.Modifier{
		dhcpv4.WithServerIP(serverIP),
		dhcpv4.WithOptionCopied(req, dhcpv4.OptionClientMachineIdentifier),
		dhcpv4.WithOptionCopied(req, dhcpv4.OptionClassIdentifier),
	}

	resp, err := dhcpv4.NewReplyFromRequest(req,
		modifiers...,
	)
	if err != nil {
		return nil, err
	}

	if resp.GetOneOption(dhcpv4.OptionClassIdentifier) == nil {
		resp.UpdateOption(dhcpv4.OptClassIdentifier("PXEClient"))
	}

	switch fwtype {
	case FirmwareX86PC:
		// This is completely standard PXE: just load a file from TFTP.
		resp.UpdateOption(dhcpv4.OptTFTPServerName(serverIP.String()))
		resp.UpdateOption(dhcpv4.OptBootFileName("undionly.kpxe"))
	case FirmwareX86Ipxe:
		// Almost standard PXE, but the boot filename needs to be a URL.
		resp.UpdateOption(dhcpv4.OptBootFileName(fmt.Sprintf("tftp://%s/undionly.kpxe", serverIP)))
	case FirmwareX86EFI:
		// This is completely standard PXE: just load a file from TFTP.
		resp.UpdateOption(dhcpv4.OptTFTPServerName(serverIP.String()))
		resp.UpdateOption(dhcpv4.OptBootFileName("snp.efi"))
	case FirmwareARMEFI:
		// This is completely standard PXE: just load a file from TFTP.
		resp.UpdateOption(dhcpv4.OptTFTPServerName(serverIP.String()))
		resp.UpdateOption(dhcpv4.OptBootFileName("snp-arm64.efi"))
	case FirmwareX86HTTP:
		// This is completely standard HTTP-boot: just load a file from HTTP.
		resp.UpdateOption(dhcpv4.OptBootFileName(fmt.Sprintf("http://%s/tftp/snp.efi", net.JoinHostPort(serverIP.String(), strconv.Itoa(apiPort)))))
	case FirmwareARMHTTP:
		// This is completely standard HTTP-boot: just load a file from HTTP.
		resp.UpdateOption(dhcpv4.OptBootFileName(fmt.Sprintf("http://%s/tftp/snp-arm64.efi", net.JoinHostPort(serverIP.String(), strconv.Itoa(apiPort)))))
	case FirmwareUnsupported:
		fallthrough
	default:
		return nil, fmt.Errorf("unsupported firmware type %d", fwtype)
	}

	return resp, nil
}

// Firmware describes a kind of firmware attempting to boot.
//
// This should only be used for selecting the right bootloader,
// kernel selection should key off the more generic architecture.
type Firmware int

// The bootloaders that we know how to handle.
const (
	FirmwareUnsupported Firmware = iota // Unsupported
	FirmwareX86PC                       // "Classic" x86 BIOS with PXE/UNDI support
	FirmwareX86EFI                      // EFI x86
	FirmwareARMEFI                      // EFI ARM64
	FirmwareX86Ipxe                     // "Classic" x86 BIOS running iPXE (no UNDI support)
	FirmwareX86HTTP                     // HTTP Boot X86
	FirmwareARMHTTP                     // ARM64 HTTP Boot
)
