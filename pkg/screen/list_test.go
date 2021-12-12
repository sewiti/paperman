package screen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseScreens(t *testing.T) {
	tests := []struct {
		name      string
		stdout    string
		want      []Screen
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			stdout: "There are screens on:\n" +
				"        370230.smh boi2 (12/12/21 02:18:17)     (Attached)\n" +
				"        370189.smh boi  (12/12/21 02:18:12)     (Detached)\n" +
				"2 Sockets in /run/screen/S-mindaugas.\n",
			want: []Screen{
				"smh boi2",
				"smh boi",
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseScreens(strings.NewReader(tt.stdout))
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
