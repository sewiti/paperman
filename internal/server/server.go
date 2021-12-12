package server

import (
	"context"
	"strings"

	"github.com/sewiti/paperman/pkg/screen"
)

type Server struct {
	Name    string
	Version string
	Port    int
	Backups int

	Java     string
	JavaArgs []string
	Jar      string
	JarArgs  []string
}

func (s Server) IsRunning(running []string) bool {
	for _, running := range running {
		if s.Name == string(running) {
			return true
		}
	}
	return false
}

func (s Server) IsRunningStandalone(ctx context.Context) (bool, error) {
	running, err := screen.ListContext(ctx, "paperman-")
	if err != nil {
		return false, err
	}
	for _, running := range FilterScreens(running) {
		if s.Name == running {
			return true, nil
		}
	}
	return false, nil
}

func FilterScreens(screens []screen.Screen) []string {
	var filtered []string
	for _, screen := range screens {
		if !strings.HasPrefix(string(screen), "paperman-") {
			continue
		}
		filtered = append(filtered, strings.TrimPrefix(string(screen), "paperman-"))
	}
	return filtered
}
