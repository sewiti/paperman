package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sewiti/paperman/pkg/atomicfs"
	"github.com/sewiti/paperman/pkg/papermc"
)

func (s Server) Launch(ctx context.Context, cwd string) error {
	_, err := os.Stat(cwd)
	if err != nil {
		return err
	}

	java := s.Java
	if java == "" {
		java = "java"
	}
	jar := s.Jar
	if jar == "" {
		// papermc
		jar, err = autoUpgrade(ctx, cwd, s.Version)
		if err != nil {
			return err
		}

		// delete old jars
		relJar, err := filepath.Rel(cwd, jar)
		if err != nil {
			return err
		}
		err = deleteFiles(cwd, "paper-", ".jar", relJar)
		if err != nil {
			return err
		}
	}
	args := append(s.JavaArgs, "-jar", jar)
	args = append(args, s.JarArgs...)

	cmd := exec.CommandContext(ctx, java, args...)
	cmd.Dir = cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func autoUpgrade(ctx context.Context, dir, version string) (jar string, err error) {
	ver, err := papermc.GetVersion(ctx, "paper", version)
	if err != nil {
		return "", err
	}
	if len(ver.Builds) == 0 {
		return "", errors.New("no builds found")
	}
	build := ver.Builds[len(ver.Builds)-1] // get latest build
	jar = filepath.Join(dir, fmt.Sprintf("paper-%s-%d.jar", version, build))

	_, err = os.Stat(jar)
	if err == nil {
		return jar, nil // already latest
	}
	if !os.IsNotExist(err) {
		return "", err
	}
	// new version/build

	r, err := papermc.Download(ctx, "paper", version, build)
	if err != nil {
		return "", err
	}
	defer r.Close()
	f, err := atomicfs.NewWriter(jar, 0640)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, r)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return jar, err
}

func deleteFiles(dir, prefix, suffix string, except ...string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
filesLoop:
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasPrefix(file.Name(), prefix) {
			continue
		}
		if !strings.HasSuffix(file.Name(), suffix) {
			continue
		}
		for _, exception := range except {
			if file.Name() == exception {
				continue filesLoop
			}
		}
		err = os.Remove(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
