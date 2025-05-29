package ui

import (
	"strings"
)

// SimpleList represents a simple list with selection
type SimpleList struct {
	items    []string
	selected int
	offset   int
	height   int
}

func (l *SimpleList) SetItems(items []string) {
	l.items = items
	if l.selected >= len(items) {
		l.selected = 0
	}
	l.adjustOffset()
}

func (l *SimpleList) SetHeight(height int) {
	l.height = height
	l.adjustOffset()
}

func (l *SimpleList) MoveUp() {
	if l.selected > 0 {
		l.selected--
		l.adjustOffset()
	}
}

func (l *SimpleList) MoveDown() {
	if l.selected < len(l.items)-1 {
		l.selected++
		l.adjustOffset()
	}
}

func (l *SimpleList) adjustOffset() {
	if l.height <= 0 {
		return
	}

	if l.selected < l.offset {
		l.offset = l.selected
	} else if l.selected >= l.offset+l.height {
		l.offset = l.selected - l.height + 1
	}
}

func (l *SimpleList) View() string {
	if len(l.items) == 0 {
		return ""
	}

	var lines []string
	start := l.offset
	end := l.offset + l.height

	if end > len(l.items) {
		end = len(l.items)
	}

	for i := start; i < end; i++ {
		line := l.items[i]
		if i == l.selected {
			// Note: Theme will be applied by the parent component
			line = "▶ " + line
		} else {
			line = "  " + line
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
