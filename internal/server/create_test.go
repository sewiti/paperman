package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createServer(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		port      string
		server    Server
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:    "gopher",
			version: "1.18.1",
			port:    "25565",
			server: Server{
				Name:     "gopher",
				Version:  "1.18.1",
				Port:     25565,
				Java:     "",
				JavaArgs: []string{"-Xms1500M", "-Xmx1500M"},
				Jar:      "",
				JarArgs:  []string{"nogui"},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "Test_createServer-*")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			err = Create(dir, tt.name, tt.version, tt.port)
			tt.assertion(t, err)

			srv, err := Read(filepath.Join(dir, tt.name))
			tt.assertion(t, err)
			assert.Equal(t, tt.server, srv)
		})
	}
}
