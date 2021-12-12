package printab

import "strings"

type layoutItemType uint

const (
	alignLeft layoutItemType = iota
	alignRight
	text
)

type layoutItem struct {
	t   layoutItemType
	val string
}

func parseLayout(layout string) (li []layoutItem, lrs int) {
	i := strings.IndexAny(layout, "lr")
	if i < 0 {
		return []layoutItem{{t: text, val: layout}}, 0
	}
	var items []layoutItem
	for i >= 0 {
		lrs++
		if i > 0 {
			items = append(items, layoutItem{t: text, val: layout[:i]})
		}
		if layout[i] == 'l' {
			items = append(items, layoutItem{t: alignLeft})
		} else {
			items = append(items, layoutItem{t: alignRight})
		}
		layout = layout[i+1:]
		i = strings.IndexAny(layout, "lr")
	}
	if layout != "" {
		items = append(items, layoutItem{t: text, val: layout})
	}
	return items, lrs
}
