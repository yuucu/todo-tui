package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestSimpleList_SetSelectedIndexPreserveScroll(t *testing.T) {
	tests := []struct {
		name            string
		items           []string
		height          int
		initialSelected int
		initialOffset   int
		newSelected     int
		expectedOffset  int
	}{
		{
			name:            "selection within visible range preserves offset",
			items:           []string{"item1", "item2", "item3", "item4", "item5"},
			height:          3,
			initialSelected: 1,
			initialOffset:   1,
			newSelected:     2,
			expectedOffset:  1, // Should preserve offset since new selection is visible
		},
		{
			name:            "selection outside visible range adjusts offset",
			items:           []string{"item1", "item2", "item3", "item4", "item5"},
			height:          3,
			initialSelected: 1,
			initialOffset:   1,
			newSelected:     4,
			expectedOffset:  2, // Should adjust to show new selection
		},
		{
			name:            "selection above visible range adjusts offset up",
			items:           []string{"item1", "item2", "item3", "item4", "item5"},
			height:          3,
			initialSelected: 2,
			initialOffset:   2,
			newSelected:     0,
			expectedOffset:  0, // Should adjust to show new selection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &SimpleList{
				items:    tt.items,
				height:   tt.height,
				selected: tt.initialSelected,
				offset:   tt.initialOffset,
			}

			list.SetSelectedIndexPreserveScroll(tt.newSelected)

			assert.Equal(t, tt.newSelected, list.selected, "Selected index should be updated")
			assert.Equal(t, tt.expectedOffset, list.offset, "Offset should be preserved or minimally adjusted")
		})
	}
}

func TestSimpleList_styleTaskContentInternal(t *testing.T) {
	// Create a mock theme for testing
	theme := &Theme{
		PriorityHigh:   lipgloss.Color("#FF0000"),
		PriorityMedium: lipgloss.Color("#FF8800"),
		PriorityLow:    lipgloss.Color("#FFFF00"),
		Secondary:      lipgloss.Color("#0088FF"),
		Primary:        lipgloss.Color("#8800FF"),
		Danger:         lipgloss.Color("#FF0000"),
		Warning:        lipgloss.Color("#FFAA00"),
		Success:        lipgloss.Color("#00FF00"),
		SelectionFg:    lipgloss.Color("#FFFFFF"),
	}

	list := &SimpleList{theme: theme}

	tests := []struct {
		name           string
		input          string
		withBackground bool
		expectContains []string
	}{
		{
			name:           "simple text without priority",
			input:          "Buy groceries",
			withBackground: false,
			expectContains: []string{"Buy", "groceries"},
		},
		{
			name:           "task with priority",
			input:          "(A) Important task",
			withBackground: false,
			expectContains: []string{"(A)", "Important", "task"},
		},
		{
			name:           "task with project and context",
			input:          "Review code +project @context",
			withBackground: false,
			expectContains: []string{"Review", "code", "+project", "@context"},
		},
		{
			name:           "task with due date",
			input:          "Finish report due:2024-12-31",
			withBackground: false,
			expectContains: []string{"Finish", "report", "due:2024-12-31"},
		},
		{
			name:           "complex task with all elements",
			input:          "(B) Complete project +work @office due:2024-01-15",
			withBackground: false,
			expectContains: []string{"(B)", "Complete", "project", "+work", "@office", "due:2024-01-15"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.withBackground {
				bgColor := lipgloss.Color("#444444")
				result = list.styleTaskContentInternal(tt.input, &bgColor)
			} else {
				result = list.styleTaskContentInternal(tt.input, nil)
			}

			// Verify that the result is not empty and contains expected content
			assert.NotEmpty(t, result, "Styled content should not be empty")

			// For simple validation, check that all expected strings are present in some form
			// (Note: actual content will have ANSI escape sequences for styling)
			for _, expected := range tt.expectContains {
				assert.Contains(t, result, expected, "Result should contain expected component")
			}
		})
	}
}

func TestSimpleList_BasicOperations(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	list := &SimpleList{}

	// Test SetItems
	list.SetItems(items)
	assert.Equal(t, items, list.items, "Items should be set correctly")
	assert.Equal(t, 0, list.selected, "Selected should default to 0")

	// Test GetSelectedItem
	assert.Equal(t, "item1", list.GetSelectedItem(), "Should return first item")

	// Test MoveDown
	list.MoveDown()
	assert.Equal(t, 1, list.selected, "Selected should move to 1")
	assert.Equal(t, "item2", list.GetSelectedItem(), "Should return second item")

	// Test MoveUp
	list.MoveUp()
	assert.Equal(t, 0, list.selected, "Selected should move back to 0")

	// Test boundary conditions
	list.MoveUp() // Should not go below 0
	assert.Equal(t, 0, list.selected, "Selected should stay at 0")

	list.SetSelectedIndex(2)
	list.MoveDown() // Should not go above last index
	assert.Equal(t, 2, list.selected, "Selected should stay at last index")
}

func TestSimpleList_ThemeHandling(t *testing.T) {
	list := &SimpleList{}
	theme := &Theme{
		Primary: lipgloss.Color("#0088FF"),
	}

	// Test theme setting
	list.SetTheme(theme)
	assert.Equal(t, theme, list.theme, "Theme should be set correctly")

	// Test task list setting
	list.SetTaskList(true)
	assert.True(t, list.isTaskList, "Should be marked as task list")

	// Test completion status setting
	completed := []bool{true, false, true}
	list.SetCompletedItems(completed)
	assert.Equal(t, completed, list.completedItems, "Completion status should be set")

	// Test checkbox colors setting
	colors := []lipgloss.Color{"#FF0000", "#00FF00", "#0000FF"}
	list.SetCheckboxColors(colors)
	assert.Equal(t, colors, list.checkboxColors, "Checkbox colors should be set")
}
