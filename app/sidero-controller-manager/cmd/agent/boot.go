// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/siderolabs/go-cmd/pkg/cmd"
	"github.com/siderolabs/go-kmsg"
	"github.com/siderolabs/go-retry/retry"
	"golang.org/x/sys/unix"
)

func setup() error {
	if err := os.MkdirAll("/etc", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/dev", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/proc", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/run", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/sys", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/tmp", 0o777); err != nil {
		return err
	}

	if err := unix.Mount("devtmpfs", "/dev", "devtmpfs", unix.MS_NOSUID, "mode=0755"); err != nil {
		return err
	}

	if err := unix.Mount("proc", "/proc", "proc", unix.MS_NOSUID|unix.MS_NOEXEC|unix.MS_NODEV, ""); err != nil {
		return err
	}

	if err := unix.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
		return err
	}

	if err := unix.Mount("tmpfs", "/run", "tmpfs", 0, ""); err != nil {
		return err
	}

	if err := unix.Mount("tmpfs", "/tmp", "tmpfs", 0, ""); err != nil {
		return err
	}

	if err := kmsg.SetupLogger(nil, "[sidero]", nil); err != nil {
		return err
	}

	// Set the PATH env var.
	if err := os.Setenv("PATH", "/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin"); err != nil {
		return errors.New("error setting PATH")
	}

	return runUdevd()
}

func runUdevd() error {
	if _, err := cmd.Run(
		"/sbin/udevadm",
		"hwdb",
		"--update",
	); err != nil {
		return fmt.Errorf("error running udevadm hwdb --update: %w", err)
	}

	udevdCmd := exec.Command("/sbin/udevd",
		"--resolve-names=never")
	udevdCmd.Stdout = os.Stdout
	udevdCmd.Stderr = os.Stderr

	if err := udevdCmd.Start(); err != nil {
		return fmt.Errorf("error starting udevd: %w", err)
	}

	if err := retry.Constant(time.Minute, retry.WithUnits(100*time.Millisecond)).Retry(func() error {
		if _, err := os.Stat("/run/udev/control"); err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error waiting for udevd to start: %w", err)
	}

	if err := retry.Constant(time.Minute, retry.WithUnits(100*time.Millisecond)).Retry(func() error {
		if _, err := cmd.Run("/sbin/udevadm", "control", "--start-exec-queue"); err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error starting udevd exec queue: %w", err)
	}

	if _, err := cmd.Run(
		"/sbin/udevadm", "trigger", "--type=devices", "--action=add",
	); err != nil {
		return fmt.Errorf("error running udevadm trigger: %w", err)
	}

	if _, err := cmd.Run(
		"/sbin/udevadm", "trigger", "--type=subsystems", "--action=add",
	); err != nil {
		return fmt.Errorf("error running udevadm trigger: %w", err)
	}

	if _, err := cmd.Run(
		"/sbin/udevadm", "settle", "--timeout=50",
	); err != nil {
		return fmt.Errorf("error running udevadm settle: %w", err)
	}

	return nil
}

func shutdown(err error) {
	if err != nil {
		log.Println(err)
	}

	for i := 10; i >= 0; i-- {
		log.Printf("rebooting in %d seconds\n", i)
		time.Sleep(1 * time.Second)
	}

	if unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART) == nil {
		select {}
	}

	os.Exit(1)
}
