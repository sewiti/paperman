package server

import (
	"strings"
)

type Server struct {
	Name    string
	Version string
	Port    int
	Backups []Backup

	Java     string
	JavaArgs []string
	Jar      string
	JarArgs  []string
}

func (s Server) Memory() string {
	const maxMemPrefix = "-Xmx"
	for _, arg := range s.JavaArgs {
		if !strings.HasPrefix(arg, maxMemPrefix) {
			continue
		}
		mem := strings.TrimPrefix(arg, maxMemPrefix)
		return strings.TrimSpace(mem)
	}
	return ""
}
