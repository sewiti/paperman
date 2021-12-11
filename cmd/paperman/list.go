package main

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/sewiti/paperman/internal/server"
	"github.com/sewiti/paperman/pkg/tmux"
)

func listServers(w io.Writer, sessions []tmux.Session, path string) error {
	servers, err := server.ReadAll(path)
	if err != nil {
		return err
	}
	if len(servers) == 0 {
		_, err = fmt.Fprintln(w, "No servers found")
		return err
	}

	maxNameLen := 4    // "Name" header
	maxVersionLen := 7 // "Version" header
	maxPortLen := 4    // "Port" header
	for _, srv := range servers {
		if l := len(srv.Name); l > maxNameLen {
			maxNameLen = l
		}
		if l := len(srv.Version); l > maxVersionLen {
			maxVersionLen = l
		}
		l := int(math.Log10(float64(srv.Port)))
		if l > maxPortLen {
			maxPortLen = l
		}
	}
	srvFmt := fmt.Sprintf("%%-%ds\t%%%ds\t%%-%dv", maxNameLen, maxVersionLen, maxPortLen)

	// Header
	_, err = fmt.Fprintf(w, srvFmt, "Name", "Version", "Port\n")
	if err != nil {
		return err
	}

	// Info
	for _, srv := range servers {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, srvFmt, srv.Name, srv.Version, srv.Port)
		if srv.IsRunning(sessions) {
			fmt.Fprint(&sb, "\tRunning")
		}
		_, err = fmt.Fprintln(w, sb.String())
		if err != nil {
			return err
		}
	}
	return nil
}
