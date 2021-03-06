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
	fmt.Printf("Launching %s...\n", s.Name)
	_, err := os.Stat(cwd)
	if err != nil {
		return err
	}

	java := s.Java
	if java == "" {
		java, err = exec.LookPath("java")
		if err != nil {
			return err
		}
	}
	jar := s.Jar
	if jar == "" {
		// papermc
		jar, err = autoUpgrade(ctx, cwd, s.Version)
		if err != nil {
			return err
		}

		// delete old jars
		err = deleteFiles(cwd, "paper-", ".jar", jar)
		if err != nil {
			return err
		}
	}
	args := append(s.JavaArgs, "-jar", jar)
	args = append(args, s.JarArgs...)

	fmt.Printf("%s %s\n", java, strings.Join(args, " "))
	cmd := exec.Command(java, args...)
	cmd.Dir = cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}

func autoUpgrade(ctx context.Context, dir, version string) (jar string, err error) {
	fmt.Println("Checking for newer version")
	ver, err := papermc.GetVersion(ctx, "paper", version)
	if err != nil {
		return "", err
	}
	if len(ver.Builds) == 0 {
		return "", errors.New("no builds found")
	}
	build := ver.Builds[len(ver.Builds)-1] // get latest build
	jar = fmt.Sprintf("paper-%s-%d.jar", version, build)

	_, err = os.Stat(filepath.Join(dir, jar))
	if err == nil {
		fmt.Printf("Already latest %s\n", jar)
		return jar, nil // already latest
	}
	if !os.IsNotExist(err) {
		return "", err
	}
	// new version/build

	fmt.Printf("Downloading %s...\n", jar)
	r, err := papermc.Download(ctx, "paper", version, build)
	if err != nil {
		return "", err
	}
	defer r.Close()
	f, err := atomicfs.NewWriter(filepath.Join(dir, jar), 0640)
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
		jar := filepath.Join(dir, file.Name())
		fmt.Printf("Deleting %s\n", jar)
		err = os.Remove(jar)
		if err != nil {
			return err
		}
	}
	return nil
}
