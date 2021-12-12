package printab

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

func Print(layout string, rows [][]interface{}) error {
	return Fprint(os.Stdout, layout, rows)
}

func Fprint(w io.Writer, layout string, rows [][]interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	li, lrs := parseLayout(layout)
	lens, err := getLens(lrs, rows)
	if err != nil {
		return err
	}
	format := buildFormat(lens, li)
	for _, row := range rows {
		_, err = fmt.Fprintf(w, format, row...)
		if err != nil {
			return err
		}
	}
	return nil
}

func getLens(n int, rows [][]interface{}) ([]int, error) {
	lens := make([]int, n)
	for _, row := range rows {
		if len(lens) != len(row) {
			return nil, errors.New("rows doesn't match layout")
		}
		for i, col := range row {
			l, err := getLen(col)
			if err != nil {
				return nil, err
			}
			if l > lens[i] {
				lens[i] = l
			}
		}
	}
	return lens, nil
}

func getLen(v interface{}) (int, error) {
	switch v := v.(type) {
	case string:
		return len(v), nil
	case int:
		return intLen(float64(v)), nil
	default:
		return 0, fmt.Errorf("unable to get length of %v", v)
	}
}

func intLen(v float64) int {
	if v < 0 {
		return int(math.Ceil(math.Log10(-v))) + 1
	}
	return int(math.Ceil(math.Log10(v)))
}

func buildFormat(colLens []int, layout []layoutItem) string {
	var format strings.Builder
	i := 0
	for _, li := range layout {
		switch li.t {
		case alignLeft, alignRight:
			l := colLens[i]
			i++
			if li.t == alignLeft {
				l = -l
			}
			fmt.Fprintf(&format, "%%%dv", l)
		case text:
			format.WriteString(li.val)
		}
	}
	return format.String()
}
