package main

import (
	_ "embed"
	"os"
	"os/exec"
)

//go:embed template/paperman@.service
var papermanService []byte

func install() error {
	err := os.WriteFile("/etc/systemd/system/paperman@.service", papermanService, 0644)
	if err != nil {
		return err
	}
	return exec.Command("systemctl", "daemon-reload").Run()
}
