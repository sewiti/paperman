package server

import (
	"context"

	"github.com/sewiti/paperman/pkg/screen"
)

func Screen(name string) screen.Screen {
	return screen.Screen("paperman-" + name)
}

func (s Server) Screen() screen.Screen {
	return Screen(s.Name)
}

func (s Server) IsRunning(running []screen.Screen) bool {
	for _, running := range running {
		if s.Screen() == running {
			return true
		}
	}
	return false
}

func (s Server) IsRunningStr(running []screen.Screen) string {
	if s.IsRunning(running) {
		return "yes"
	}
	return "no"
}

func (s Server) IsRunningStandalone(ctx context.Context) (bool, error) {
	running, err := screen.ListContext(ctx, "paperman-")
	if err != nil {
		return false, err
	}
	return s.IsRunning(running), nil
}
