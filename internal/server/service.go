package server

import (
	"context"
	"fmt"

	"github.com/sewiti/paperman/internal/systemd"
)

const (
	PapermanService    = "paperman@.service"
	fmtPapermanService = "paperman@%s.service"
)

func Service(name string) systemd.Service {
	return systemd.Service(fmt.Sprintf(fmtPapermanService, name))
}

func (s Server) Service() systemd.Service {
	return Service(s.Name)
}

func (s Server) IsEnabled(enabled []systemd.Service) bool {
	for _, enabled := range enabled {
		if s.Service() == enabled {
			return true
		}
	}
	return false
}

func (s Server) IsEnabledStr(enabled []systemd.Service) string {
	if s.IsEnabled(enabled) {
		return "yes"
	}
	return "no"
}

func (s Server) IsEnabledStandalone(ctx context.Context) (bool, error) {
	enabled, err := systemd.ListWanted("multi-user.target")
	if err != nil {
		return false, err
	}
	return s.IsEnabled(enabled), nil
}
