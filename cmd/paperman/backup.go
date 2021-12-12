package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sewiti/paperman/internal/server"
	"github.com/sewiti/paperman/pkg/screen"
)

func backup(ctx context.Context, srvDir, name string) error {
	running, err := server.Server{Name: name}.IsRunningStandalone(ctx)
	if err != nil {
		return err
	}
	if running {
		fmt.Printf("%s is running, turning off world saving for the backup\n", name)
		screen := screen.Screen("paperman-" + name)
		err := screen.SendStuffContext(ctx, "save-off")
		if err != nil {
			return err
		}
		defer screen.SendStuffContext(ctx, "save-on")
		err = screen.SendStuffContext(ctx, "save-all")
		if err != nil {
			return err
		}
		time.Sleep(3 * time.Second) // Some delay for save-all
	}

	instDir := filepath.Join(srvDir, name)
	instBackups := filepath.Join(instDir, server.BackupsDir)

	stat, err := os.Stat(instBackups)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.Mkdir(instBackups, 0750)
		if err != nil {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("%s is a file", instBackups)
	}

	fname := time.Now().Format(server.BackupTimeLayout) + ".tar.gz"
	path := filepath.Join(instBackups, fname)

	fmt.Printf("Creating backup %s\n", path)
	cmd := exec.CommandContext(ctx, "tar", "cvzf", path, "-C", instDir, "--exclude", server.BackupsDir, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Created %s\n", path)

	if running {
		fmt.Println("Re-enabling world saving") // it's deferred
	}
	return nil
}

func purgeBackups(srvDir, name string, count int) error {
	srv, err := server.Read(filepath.Join(srvDir, name))
	if err != nil {
		return err
	}
	for i, b := range srv.Backups {
		if i+count >= len(srv.Backups) {
			break // Leave `count` Backups
		}
		err = os.Remove(b.Path)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted %s\n", b.Path)
	}
	return nil
}
