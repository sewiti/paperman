package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sewiti/paperman/internal/conf"
)

const (
	PapermanConf = "paperman.conf"
	Properties   = "server.properties"
)

func ReadAll(dirPath string) ([]Server, error) {
	dirs, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var servers []Server
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		s, err := Read(filepath.Join(dirPath, dir.Name()))
		if err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func Read(path string) (Server, error) {
	paper, err := conf.Read(filepath.Join(path, PapermanConf))
	if err != nil {
		return Server{}, err
	}
	props, err := conf.Read(filepath.Join(path, Properties))
	if err != nil {
		return Server{}, err
	}
	backups, err := readBackups(filepath.Join(path, BackupsDir))
	if err != nil {
		return Server{}, err
	}

	port, err := strconv.Atoi(props.Get("server-port"))
	if err != nil {
		return Server{}, fmt.Errorf("port: %w", err)
	}
	return Server{
		Name:    filepath.Base(path),
		Port:    port,
		Version: paper.Get("papermc-version"),
		Backups: backups,

		Java:     paper.Get("java"),
		JavaArgs: paper["java-args"],
		Jar:      paper.Get("jar"),
		JarArgs:  paper["jar-args"],
	}, nil
}
