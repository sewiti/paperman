package server

import (
	"github.com/sewiti/paperman/pkg/tmux"
)

type Server struct {
	Name    string
	Version string
	Port    int

	Java     string
	JavaArgs []string
	Jar      string
	JarArgs  []string
}

func (s Server) IsRunning(sessions []tmux.Session) bool {
	for _, session := range sessions {
		if s.Name == string(session) {
			return true
		}
	}
	return false
}
