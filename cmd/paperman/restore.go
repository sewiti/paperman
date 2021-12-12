package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sewiti/paperman/internal/server"
)

func restore(ctx context.Context, srvDir, name, archive string) error {
	running, err := server.Server{Name: name}.IsRunningStandalone(ctx)
	if err != nil {
		return err
	}
	if running {
		return fmt.Errorf("%s is running", name)
	}

	err = verifyArchive(ctx, archive)
	if err != nil {
		return err
	}
	dir := filepath.Join(srvDir, name)

	fmt.Printf("Deleting %s\n", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() && f.Name() == server.BackupsDir {
			continue // skip backups dir
		}
		err = os.RemoveAll(filepath.Join(dir, f.Name()))
		if err != nil {
			return err
		}
	}

	fmt.Printf("Restoring %s from %s\n", name, archive)
	cmd := exec.CommandContext(ctx, "tar", "xvf", archive, "-C", dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Restored %s successfully\n", name)
	return nil
}

func verifyArchive(ctx context.Context, path string) error {
	return exec.CommandContext(ctx, "tar", "tf", path).Run()
}
