package main

import (
	"errors"
	"fmt"
	"os/exec"
)

func systemControl(action, name string) error {
	switch action {
	case "start", "restart", "stop", "enable", "disable":
		// continue
	default:
		return errors.New("invalid action")
	}
	service := fmt.Sprintf("paperman@%s.service", name)
	cmd := exec.Command("systemctl", action, service)
	return cmd.Run()
}
