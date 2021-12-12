package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

	fname := time.Now().Format("2006-01-02_15:04:05") + ".tar.gz"
	path := filepath.Join(instBackups, fname)

	cmd := exec.CommandContext(ctx, "tar", "czf", path, "-C", instDir, "--exclude", server.BackupsDir, ".")
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Created %s\n", path)
	return nil
}

func purgeBackups(srcDir, name string, count int) error {
	instBackups := filepath.Join(srvDir, name, server.BackupsDir)

	backups, err := os.ReadDir(instBackups)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	type pair struct {
		backup  fs.DirEntry
		modTime time.Time
	}
	pairs := make([]pair, 0, len(backups))
	for _, backup := range backups {
		stat, err := os.Stat(filepath.Join(instBackups, backup.Name()))
		if err != nil {
			return err
		}
		pairs = append(pairs, pair{
			backup:  backup,
			modTime: stat.ModTime(),
		})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].modTime.Before(pairs[j].modTime)
	})

	for i := 0; i < len(pairs)-count; i++ {
		file := filepath.Join(instBackups, pairs[i].backup.Name())
		err = os.Remove(file)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted %s\n", file)
	}
	return nil
}
