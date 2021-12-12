package printab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseLayout(t *testing.T) {
	tests := []struct {
		layout string
		want   []layoutItem
		wantN  int
	}{
		{
			layout: "lr",
			want: []layoutItem{
				{t: alignLeft},
				{t: alignRight},
			},
			wantN: 2,
		},
		{
			layout: "l\tr\n",
			want: []layoutItem{
				{t: alignLeft},
				{t: text, val: "\t"},
				{t: alignRight},
				{t: text, val: "\n"},
			},
			wantN: 2,
		},
		{
			layout: "|l|r|\n",
			want: []layoutItem{
				{t: text, val: "|"},
				{t: alignLeft},
				{t: text, val: "|"},
				{t: alignRight},
				{t: text, val: "|\n"},
			},
			wantN: 2,
		},
		{
			layout: "===\n",
			want: []layoutItem{
				{t: text, val: "===\n"},
			},
			wantN: 0,
		},
		{
			layout: "|l|\n",
			want: []layoutItem{
				{t: text, val: "|"},
				{t: alignLeft},
				{t: text, val: "|\n"},
			},
			wantN: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			li, lrs := parseLayout(tt.layout)
			assert.Equal(t, tt.want, li)
			assert.Equal(t, tt.wantN, lrs)
		})
	}
}
