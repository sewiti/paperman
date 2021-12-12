package printab

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFprint(t *testing.T) {
	tests := []struct {
		layout    string
		rows      [][]interface{}
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			layout: "l\tl\n",
			rows: [][]interface{}{
				{"HEAD", "TAIL"},
				{"<", ">"},
				{"LONG_BOI", ">>"},
			},
			wantW: "" +
				"HEAD    \tTAIL\n" +
				"<       \t>   \n" +
				"LONG_BOI\t>>  \n",
			assertion: assert.NoError,
		},
		{
			layout: "|l|r|\n",
			rows: [][]interface{}{
				{"ABCDEF", "JIFF no JEFF"},
				{"abc", "gify"},
				{"def", "jef"},
			},
			wantW: "" +
				"|ABCDEF|JIFF no JEFF|\n" +
				"|abc   |        gify|\n" +
				"|def   |         jef|\n",
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			var w strings.Builder
			tt.assertion(t, Fprint(&w, tt.layout, tt.rows))
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}
