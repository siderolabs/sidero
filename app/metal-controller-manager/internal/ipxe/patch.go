// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipxe

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PatchBinaries patches iPXE binaries on the fly with the new embedded script.
//
// This relies on special build in `pkgs/ipxe` where a placeholder iPXE script is embedded.
// EFI iPXE binaries are uncompressed, so these are patched directly.
// BIOS amd64 undionly.pxe is compressed, so we instead patch uncompressed version and compress it back using zbin.
// (zbin is built with iPXE).
func PatchBinaries(script []byte) error {
	if err := patchScript("/var/lib/sidero/ipxe/amd64/ipxe.efi", "/var/lib/sidero/tftp/ipxe.efi", script); err != nil {
		return err
	}

	if err := patchScript("/var/lib/sidero/ipxe/arm64/ipxe.efi", "/var/lib/sidero/tftp/ipxe-arm64.efi", script); err != nil {
		return err
	}

	if err := patchScript("/var/lib/sidero/ipxe/amd64/kpxe/undionly.kpxe.bin", "/var/lib/sidero/ipxe/amd64/kpxe/undionly.kpxe.bin.patched", script); err != nil {
		return err
	}

	if err := compressKPXE("/var/lib/sidero/ipxe/amd64/kpxe/undionly.kpxe.bin.patched", "/var/lib/sidero/ipxe/amd64/kpxe/undionly.kpxe.zinfo", "/var/lib/sidero/tftp/undionly.kpxe"); err != nil {
		return err
	}

	if err := compressKPXE("/var/lib/sidero/ipxe/amd64/kpxe/undionly.kpxe.bin.patched", "/var/lib/sidero/ipxe/amd64/kpxe/undionly.kpxe.zinfo", "/var/lib/sidero/tftp/undionly.kpxe.0"); err != nil {
		return err
	}

	return nil
}

var (
	placeholderStart = []byte("# *PLACEHOLDER START*")
	placeholderEnd   = []byte("# *PLACEHOLDER END*")
)

func patchScript(source, destination string, script []byte) error {
	contents, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	start := bytes.Index(contents, placeholderStart)
	if start == -1 {
		return fmt.Errorf("placeholder start not found in %q", source)
	}

	end := bytes.Index(contents, placeholderEnd)
	if end == -1 {
		return fmt.Errorf("placeholder end not found in %q", source)
	}

	if end < start {
		return fmt.Errorf("placeholder end before start")
	}

	end += len(placeholderEnd)

	length := end - start

	if len(script) > length {
		return fmt.Errorf("script size %d is larger than placeholder space %d", len(script), length)
	}

	script = append(script, bytes.Repeat([]byte{'\n'}, length-len(script))...)

	copy(contents[start:end], script)

	if err = os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}

	return os.WriteFile(destination, contents, 0o644)
}

// compressPXE is equivalent to: ./util/zbin bin/undionly.kpxe.bin bin/undionly.kpxe.zinfo > bin/undionly.kpxe.zbin.
func compressKPXE(binFile, infoFile, outFile string) error {
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}

	defer out.Close()

	cmd := exec.Command("/bin/zbin", binFile, infoFile)
	cmd.Stdout = out

	return cmd.Run()
}
