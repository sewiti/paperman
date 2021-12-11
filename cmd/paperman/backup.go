package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/sewiti/paperman/pkg/atomicfs"
)

func backup(srvDir, name string) error {
	backupsDir := filepath.Join(srvDir, name, "backups")
	stat, err := os.Stat(backupsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.Mkdir(backupsDir, 0750)
		if err != nil {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("%s is a file", backupsDir)
	}

	path := time.Now().Format(time.RFC3339) + ".tar.gz"
	f, err := atomicfs.NewWriter(path, 0640)
	if err != nil {
		return err
	}
	err = compressGzip(f, filepath.Join(srvDir, name), gzip.DefaultCompression)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

func compressGzip(w io.Writer, dir string, level int) error {
	g, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return err
	}
	err = archiveTar(g, dir)
	if err1 := g.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

func archiveTar(w io.Writer, dir string) error {
	t := tar.NewWriter(w)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(relPath)

		err = t.WriteHeader(hdr)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(t, f)
		return err
	})

	if err1 := t.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}
