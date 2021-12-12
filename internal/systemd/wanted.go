package systemd

import (
	"os"
	"path/filepath"
)

type Service string
type Target string

func ListWanted(by Target) ([]Service, error) {
	wants := filepath.Join("/etc/systemd/system", string(by)+".wants")
	files, err := os.ReadDir(wants)
	if err != nil {
		return nil, err
	}
	services := make([]Service, 0, len(files))
	for _, wanted := range files {
		services = append(services, Service(wanted.Name()))
	}
	return services, nil
}
