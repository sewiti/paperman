package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sewiti/paperman/pkg/tmux"
	"github.com/stretchr/testify/assert"
)

func Test_listServers(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		sessions  []tmux.Session
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "servers list",
			path: filepath.Join("testdata", "list"),
			sessions: []tmux.Session{
				"gopher",
			},
			want: "Name  \tVersion\tPort\n" +
				"gopher\t 1.18.1\t25588\tRunning\n" +
				"lemur \t   1.18\t25565\n",
			assertion: assert.NoError,
		},
		{
			name: "no servers",
			path: filepath.Join("testdata", "emptylist"),
			sessions: []tmux.Session{
				"gopher",
			},
			want:      "No servers found\n",
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := strings.Builder{}
			tt.assertion(t, listServers(&sb, tt.sessions, tt.path))
			assert.Equal(t, tt.want, sb.String())
		})
	}
}
