package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Constants for list rendering
const (
	selectionIndicator = "▶ "
	spacing            = "  "
	checkboxCompleted  = "● "
	checkboxIncomplete = "○ "
)

// SimpleList represents a simple list with selection and enhanced styling
type SimpleList struct {
	items          []string
	selected       int
	offset         int
	height         int
	theme          *Theme           // Theme for styling
	isTaskList     bool             // Whether this is a task list (affects rendering)
	completedItems []bool           // Track which items are completed
	checkboxColors []lipgloss.Color // Track checkbox colors for incomplete tasks
}

// SetTheme sets the theme for styling
func (l *SimpleList) SetTheme(theme *Theme) {
	l.theme = theme
}

// SetTaskList sets whether this list displays tasks (affects rendering)
func (l *SimpleList) SetTaskList(isTaskList bool) {
	l.isTaskList = isTaskList
}

// SetCompletedItems sets which items are completed
func (l *SimpleList) SetCompletedItems(completed []bool) {
	l.completedItems = completed
}

// SetCheckboxColors sets colors for incomplete task checkboxes
func (l *SimpleList) SetCheckboxColors(colors []lipgloss.Color) {
	l.checkboxColors = colors
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

// GetSelectedIndex returns the currently selected index
func (l *SimpleList) GetSelectedIndex() int {
	return l.selected
}

// SetSelectedIndex sets the selected index and adjusts offset
func (l *SimpleList) SetSelectedIndex(index int) {
	if index >= 0 && index < len(l.items) {
		l.selected = index
		l.adjustOffset()
	}
}

// SetSelectedIndexPreserveScroll sets the selected index while trying to preserve scroll position
func (l *SimpleList) SetSelectedIndexPreserveScroll(index int) {
	if index >= 0 && index < len(l.items) {
		oldOffset := l.offset
		l.selected = index

		// Only adjust offset if the new selection is outside the current visible range
		if l.selected >= oldOffset && l.selected < oldOffset+l.height {
			// Selection is still visible, keep the old offset
			l.offset = oldOffset
		} else {
			// Selection is outside visible range, adjust minimally
			l.adjustOffset()
		}
	}
}

// GetSelectedItem returns the currently selected item
func (l *SimpleList) GetSelectedItem() string {
	if l.selected >= 0 && l.selected < len(l.items) {
		return l.items[l.selected]
	}
	return ""
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

	// Only adjust offset if the selected item is outside the visible range
	// This helps preserve the current scroll position
	if l.selected < l.offset {
		// Selected item is above the visible range - scroll up minimally
		l.offset = l.selected
	} else if l.selected >= l.offset+l.height {
		// Selected item is below the visible range - scroll down minimally
		l.offset = l.selected - l.height + 1
	}
	// If selected item is already visible, don't change offset
}

func (l *SimpleList) View() string {
	// Safety check for height
	if l.height <= 0 {
		return ""
	}

	var lines []string
	start := l.offset
	end := l.offset + l.height

	// Add items within the visible range
	if len(l.items) > 0 {
		if end > len(l.items) {
			end = len(l.items)
		}

		for i := start; i < end; i++ {
			line := l.items[i]

			// Apply different styling for task lists vs filter lists
			if l.isTaskList {
				line = l.renderTaskItem(line, i)
			} else {
				line = l.renderFilterItem(line, i)
			}

			lines = append(lines, line)
		}
	}

	// Fill remaining lines with empty content to match the set height
	for len(lines) < l.height {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// renderTaskItem renders a task item with modern styling
func (l *SimpleList) renderTaskItem(item string, index int) string {
	if l.theme == nil {
		// Fallback to simple rendering if no theme
		if index == l.selected {
			return selectionIndicator + item
		}
		return spacing + item
	}

	// Determine if this task is completed
	isCompleted := index < len(l.completedItems) && l.completedItems[index]

	// Create modern checkbox with circles (● / ○)
	var checkbox string
	if isCompleted {
		// Filled circle for completed tasks with green color
		checkboxStyle := lipgloss.NewStyle().
			Foreground(l.theme.Success).
			Bold(true)
		checkbox = checkboxStyle.Render(checkboxCompleted)
	} else {
		// Empty circle for incomplete tasks with dynamic color based on due date
		var checkboxColor lipgloss.Color
		if index < len(l.checkboxColors) && l.checkboxColors[index] != "" {
			// Use provided due date color
			checkboxColor = l.checkboxColors[index]
		} else {
			// Fallback to muted color
			checkboxColor = l.theme.TextMuted
		}

		checkboxStyle := lipgloss.NewStyle().
			Foreground(checkboxColor).
			Bold(false)
		checkbox = checkboxStyle.Render(checkboxIncomplete)
	}

	// Apply selection highlighting if this item is selected
	if index == l.selected {
		// Selection indicator with theme color
		indicator := lipgloss.NewStyle().
			Foreground(l.theme.Primary).
			Bold(true).
			Render(selectionIndicator)

		// For selected items, apply uniform background highlighting to entire content
		var content string
		if isCompleted {
			// For completed selected tasks: keep strikethrough but override colors
			contentStyle := lipgloss.NewStyle().
				Background(l.theme.SelectionBg).
				Foreground(l.theme.SelectionFg).
				Strikethrough(true).
				Bold(true)
			content = contentStyle.Render(item)
		} else {
			// For incomplete selected tasks: apply selection highlighting to entire content
			contentStyle := lipgloss.NewStyle().
				Background(l.theme.SelectionBg).
				Foreground(l.theme.SelectionFg).
				Bold(true)
			content = contentStyle.Render(item)
		}

		// Combine components: indicator + checkbox + highlighted content
		return indicator + checkbox + content
	}

	// Non-selected item - parse and style components individually
	var content string
	if isCompleted {
		// For completed tasks, apply both strikethrough and muted color to entire content
		contentStyle := lipgloss.NewStyle().
			Foreground(l.theme.TextMuted).
			Strikethrough(true)
		content = contentStyle.Render(item)
	} else {
		// For active tasks, parse and style individual components
		content = l.styleActiveTaskContent(item)
	}

	// Add spacing for non-selected items
	return spacing + checkbox + content
}

// styleActiveTaskContent parses and styles components of an active task
func (l *SimpleList) styleActiveTaskContent(item string) string {
	// Split the content to parse priority, todo text, and tags
	parts := strings.Fields(item)
	if len(parts) == 0 {
		return item
	}

	var styledParts []string

	for i, part := range parts {
		if i == 0 && strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") && len(part) == 3 {
			// This is a priority like "(A)"
			priority := strings.Trim(part, "()")
			priorityStyle := lipgloss.NewStyle().Bold(true)
			switch priority {
			case "A":
				priorityStyle = priorityStyle.Foreground(l.theme.PriorityHigh)
			case "B":
				priorityStyle = priorityStyle.Foreground(l.theme.PriorityMedium)
			case "C":
				priorityStyle = priorityStyle.Foreground(l.theme.PriorityLow)
			case "D":
				priorityStyle = priorityStyle.Foreground(l.theme.PriorityLowest)
			default:
				priorityStyle = priorityStyle.Foreground(l.theme.PriorityDefault)
			}
			styledParts = append(styledParts, priorityStyle.Render(part))
		} else if strings.HasPrefix(part, "+") {
			// Project tag
			projectStyle := lipgloss.NewStyle().Foreground(l.theme.Secondary)
			styledParts = append(styledParts, projectStyle.Render(part))
		} else if strings.HasPrefix(part, "@") {
			// Context tag
			contextStyle := lipgloss.NewStyle().Foreground(l.theme.Primary)
			styledParts = append(styledParts, contextStyle.Render(part))
		} else if strings.HasPrefix(part, "due:") {
			// Due date tag
			dueStyle := lipgloss.NewStyle()
			dueDate := strings.TrimPrefix(part, "due:")

			// Parse and color based on date
			now := time.Now()
			today := now.Format("2006-01-02")

			if dueDate < today {
				dueStyle = dueStyle.Foreground(l.theme.Danger) // Overdue
			} else if dueDate == today {
				dueStyle = dueStyle.Foreground(l.theme.Warning) // Due today
			} else {
				dueStyle = dueStyle.Foreground(l.theme.Success) // Future
			}

			styledParts = append(styledParts, dueStyle.Render(part))
		} else {
			// Regular text (todo content)
			styledParts = append(styledParts, part)
		}
	}

	return strings.Join(styledParts, " ")
}

// renderFilterItem renders a filter item with highlighting
func (l *SimpleList) renderFilterItem(item string, index int) string {
	if l.theme == nil {
		// Fallback to simple rendering if no theme
		if index == l.selected {
			return selectionIndicator + item
		}
		return spacing + item
	}

	// Apply selection highlighting for filter items
	if index == l.selected {
		// Selection indicator with theme color
		indicator := lipgloss.NewStyle().
			Foreground(l.theme.Primary).
			Bold(true).
			Render(selectionIndicator)

		// Apply selection highlighting to the content
		contentStyle := lipgloss.NewStyle().
			Background(l.theme.SelectionBg).
			Foreground(l.theme.SelectionFg).
			Bold(true)

		content := contentStyle.Render(item)
		return indicator + content
	}

	// Non-selected filter item
	return spacing + item
}
