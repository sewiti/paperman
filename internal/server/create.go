package server

import (
	_ "embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sewiti/paperman/internal/atomicfs"
	"github.com/sewiti/paperman/internal/conf"
)

//go:embed template/server.properties
var srvPropsTpl string

func Create(srvDir, name, version, port string) error {
	err := os.Mkdir(filepath.Join(srvDir, name), 0750)
	if err != nil {
		return err
	}

	// server.properties
	tpl, err := template.New(Properties).Parse(srvPropsTpl)
	if err != nil {
		return err
	}
	f, err := atomicfs.NewWriter(filepath.Join(srvDir, name, Properties), 0640)
	if err != nil {
		return err
	}
	err = tpl.Execute(f, map[string]string{
		"name": name,
		"port": port,
	})
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err != nil {
		return err
	}

	// paperman.conf
	err = conf.Write(filepath.Join(srvDir, name, PapermanConf), conf.Values{
		"papermc-version": []string{version},
		"java-args":       []string{"-Xms1500M", "-Xmx1500M"},
		"jar-args":        []string{"nogui"},
	})
	if err != nil {
		return err
	}

	// eula.txt
	err = atomicfs.WriteFile(filepath.Join(srvDir, name, "eula.txt"), []byte("eula=true"), 0640)
	return err
}
