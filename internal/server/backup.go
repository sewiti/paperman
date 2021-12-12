package server

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	BackupsDir       = "backups"
	BackupTimeLayout = "2006-01-02_15-04-05"
)

type Backup struct {
	Path    string
	ModTime time.Time
}

func readBackups(path string) ([]Backup, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	// Collect tarballs & their mod times
	backups := make([]Backup, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !isTarball(f.Name()) {
			continue
		}

		stat, err := f.Info()
		if err != nil {
			return nil, err
		}
		backups = append(backups, Backup{
			Path:    filepath.Join(path, f.Name()),
			ModTime: stat.ModTime(),
		})
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].ModTime.Before(backups[j].ModTime)
	})
	return backups, nil
}

func isTarball(name string) bool {
	parts := strings.Split(name, ".")
	if len(parts) < 2 {
		return false
	}
	for i := len(parts) - 1; i > 0; i-- {
		if parts[i] == "tar" {
			return true
		}
	}
	return false
}
