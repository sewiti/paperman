package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/sewiti/paperman/internal/server"
)

func listInstances(w io.Writer, runningInstances []string, path string) error {
	servers, err := server.ReadAll(path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = fmt.Fprintln(w, "No servers found")
			return err
		}
		return err
	}

	maxNameLen := 4    // "NAME" header
	maxVersionLen := 7 // "VERSION" header
	maxPortLen := 4    // "PORT" header
	for _, srv := range servers {
		if l := len(srv.Name); l > maxNameLen {
			maxNameLen = l
		}
		if l := len(srv.Version); l > maxVersionLen {
			maxVersionLen = l
		}
		l := int(math.Ceil(math.Log10(float64(srv.Port))))
		if l > maxPortLen {
			maxPortLen = l
		}
	}
	srvFmt := fmt.Sprintf("%%%ds\t%%%dv\t%%%ds\t%%7s\t%%7v\n", -maxNameLen, maxPortLen, -maxVersionLen)

	// Header
	_, err = fmt.Fprintf(w, srvFmt, "NAME", "PORT", "VERSION", "RUNNING", "BACKUPS")
	if err != nil {
		return err
	}

	// Info
	for _, srv := range servers {
		running := "no"
		if srv.IsRunning(runningInstances) {
			running = "yes"
		}
		_, err = fmt.Fprintf(w, srvFmt, srv.Name, srv.Port, srv.Version, running, srv.Backups)
		if err != nil {
			return err
		}
	}
	return nil
}
