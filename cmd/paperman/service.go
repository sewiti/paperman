package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sewiti/paperman/internal/server"
)

//go:embed template/paperman@.service
var defaultPapermanSrvc []byte

func install() error {
	if os.Geteuid() != 0 {
		return errors.New("root privilleges required")
	}

	bs, err := os.ReadFile(server.PapermanService)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		bs = defaultPapermanSrvc
	}
	dir := "/usr/lib/systemd/system"
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		dir = "/lib/systemd/system"
	}
	dst := filepath.Join(dir, server.PapermanService)
	err = os.WriteFile(dst, bs, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Installed %s\n", dst)
	return systemControl("daemon-reload", "")
}

func systemControl(action, name string) error {
	switch action {
	case "start", "restart", "stop", "enable", "disable", "daemon-reload":
		// continue
	default:
		return errors.New("invalid action")
	}

	args := []string{action}
	if action != "daemon-reload" {
		service := server.Service(name)
		args = append(args, string(service))
	}

	var cmd *exec.Cmd
	if os.Geteuid() == 0 {
		cmd = exec.Command("systemctl", args...)
	} else {
		cmd = exec.Command("sudo", append([]string{"systemctl"}, args...)...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
