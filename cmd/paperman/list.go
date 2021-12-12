package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sewiti/paperman/internal/server"
	"github.com/sewiti/paperman/internal/systemd"
	"github.com/sewiti/paperman/pkg/printab"
	"github.com/sewiti/paperman/pkg/screen"
)

func listInstances(w io.Writer, running []screen.Screen, path string) error {
	servers, err := server.ReadAll(path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = fmt.Fprintln(w, "No servers found")
			return err
		}
		return err
	}
	enabled, err := systemd.ListWanted("multi-user.target")
	if err != nil {
		return err
	}

	const layout = "l\tl\tr\tr\tr\tr\tr\tr\n"
	rows := [][]interface{}{
		{"NAME", "VERSION", "PORT", "RUNNING", "ENABLED", "MEMORY", "BACKUPS", "LAST BACKUP"}, // Header
	}

	for _, srv := range servers {
		lastBackup := ""
		if len(srv.Backups) > 0 {
			b := srv.Backups[len(srv.Backups)-1]
			lastBackup = b.ModTime.Format(server.BackupTimeLayout)
		}
		rows = append(rows, []interface{}{
			srv.Name,
			srv.Version,
			srv.Port,
			srv.IsRunningStr(running),
			srv.IsEnabledStr(enabled),
			srv.Memory(),
			len(srv.Backups),
			lastBackup,
		})
	}
	return printab.Fprint(w, layout, rows)
}
