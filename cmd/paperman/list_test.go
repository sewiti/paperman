package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_listServers(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		running   []string
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "servers list",
			path: filepath.Join("testdata", "list"),
			running: []string{
				"gopher",
			},
			want: "" +
				"NAME  \t PORT\tVERSION\tRUNNING\tBACKUPS\n" +
				"gopher\t25588\t1.18.1 \t    yes\t      0\n" +
				"lemur \t25565\t1.18   \t     no\t      0\n",
			assertion: assert.NoError,
		},
		{
			name: "no servers",
			path: filepath.Join("testdata", "emptylist"),
			running: []string{
				"gopher",
			},
			want:      "No servers found\n",
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := strings.Builder{}
			tt.assertion(t, listInstances(&sb, tt.running, tt.path))
			assert.Equal(t, tt.want, sb.String())
		})
	}
}
