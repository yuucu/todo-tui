package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// initializeHelpContent initializes the help content for the application
func (m *Model) initializeHelpContent() {
	m.helpContent = []HelpContent{
		{
			Category: "Global Commands",
			Items: []HelpItem{
				{"?", "Show/hide this help screen"},
				{"q / Ctrl+C", "Quit application"},
				{"a", "Add new task"},
				{"e", "Edit selected task"},
				{"d", "Delete selected task"},
				{"r", "Restore deleted/completed task"},
			},
		},
		{
			Category: "Navigation",
			Items: []HelpItem{
				{"Tab", "Switch between panes"},
				{"h", "Move to left pane (Workspaces)"},
				{"l", "Move to right pane (Todos)"},
				{"j / ↓", "Move down / Scroll down (in help)"},
				{"k / ↑", "Move up / Scroll up (in help)"},
				{"Enter", "Apply filter / Complete task"},
			},
		},
		{
			Category: "Help Navigation",
			Items: []HelpItem{
				{"j / ↓", "Scroll down"},
				{"k / ↑", "Scroll up"},
				{"g", "Go to top"},
				{"G", "Go to bottom"},
				{"Any other key", "Close help"},
			},
		},
		{
			Category: "Task Operations",
			Items: []HelpItem{
				{"y", "Copy task text to clipboard"},
				{"p", "Cycle task priority"},
				{"t", "Toggle due date to today"},
			},
		},
		{
			Category: "Edit Mode",
			Items: []HelpItem{
				{"Esc / Ctrl+C", "Cancel editing"},
				{"Enter / Ctrl+S", "Save task"},
			},
		},
	}
}

// renderHelpView renders the help screen with scrolling support
func (m *Model) renderHelpView() string {
	if len(m.helpContent) == 0 {
		m.initializeHelpContent()
	}

	// Calculate visible area for scrolling
	maxWidth := min(80, m.width-4) // Maximum 80 characters wide, leave 4 for margins
	visibleHeight := m.height - 6  // Leave space for borders and padding

	var content strings.Builder

	// Define styles (restored from original design)
	headerStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Primary).
		Bold(true).
		Padding(0, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(m.currentTheme.Primary).
		Align(lipgloss.Center)

	categoryStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Secondary).
		Bold(true).
		Underline(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Warning).
		Bold(true).
		Width(18).
		Align(lipgloss.Left).
		Background(m.currentTheme.Surface).
		Padding(0, 1).
		MarginRight(1)

	descStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Text)

	// Prepare content sections as simple strings for easier scrolling
	var allSections []string

	// Header and empty line
	allSections = append(allSections,
		headerStyle.Render("📚 TodoTUI - Keyboard Shortcuts Help"),
		"", // Empty line
	)

	// Content sections
	for i, category := range m.helpContent {
		if i > 0 {
			allSections = append(allSections, "") // Empty line between categories
		}

		// Category header
		allSections = append(allSections, categoryStyle.Render("▶ "+category.Category))

		// Category items in a box
		var itemsContent strings.Builder
		for _, item := range category.Items {
			keyPart := keyStyle.Render(item.Key)
			descPart := descStyle.Render(item.Description)
			itemsContent.WriteString(keyPart + " │ " + descPart + "\n")
		}

		// Box style for items
		itemBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.currentTheme.BorderInactive).
			Padding(0, 2).
			MarginLeft(2).
			Background(m.currentTheme.Surface)

		boxedItems := itemBoxStyle.Render(strings.TrimSuffix(itemsContent.String(), "\n"))
		boxLines := strings.Split(boxedItems, "\n")
		allSections = append(allSections, boxLines...)
	}

	// Footer
	allSections = append(allSections, "") // Empty line

	// Apply scrolling
	totalLines := len(allSections)
	maxScroll := max(0, totalLines-visibleHeight+2) // Leave space for scroll indicator
	if m.helpScroll > maxScroll {
		m.helpScroll = maxScroll
	}
	if m.helpScroll < 0 {
		m.helpScroll = 0
	}

	// Get visible lines
	startLine := m.helpScroll
	endLine := min(startLine+visibleHeight-2, totalLines) // Reserve 2 lines for scroll indicator
	if endLine > totalLines {
		endLine = totalLines
	}

	visibleSections := allSections[startLine:endLine]

	// Build final content
	for _, section := range visibleSections {
		content.WriteString(section + "\n")
	}

	// Add scroll indicator if needed
	if totalLines > visibleHeight-2 {
		content.WriteString("\n")
		scrollInfo := fmt.Sprintf("↑ %d-%d/%d ↓", startLine+1, endLine, totalLines)
		if m.helpScroll == 0 {
			scrollInfo = fmt.Sprintf("  %d-%d/%d ↓", startLine+1, endLine, totalLines)
		} else if m.helpScroll >= maxScroll {
			scrollInfo = fmt.Sprintf("↑ %d-%d/%d  ", startLine+1, endLine, totalLines)
		}

		scrollIndicator := lipgloss.NewStyle().
			Foreground(m.currentTheme.TextMuted).
			Align(lipgloss.Center).
			Width(maxWidth).
			Render(scrollInfo)
		content.WriteString(scrollIndicator)
	}

	// Style for the help content container
	helpContentStyle := lipgloss.NewStyle().
		Width(maxWidth).
		Align(lipgloss.Left).
		Padding(2).
		Background(m.currentTheme.Background)

	styledContent := helpContentStyle.Render(content.String())

	// Center the container on the screen
	centeredStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Background(m.currentTheme.Background)

	return centeredStyle.Render(styledContent)
}

// KeyBinding represents a key binding configuration
type KeyBinding struct {
	Key         string
	Description string
	Mode        ViewMode
	Pane        *Pane // nil means any pane
	Handler     func(m *Model) (tea.Model, tea.Cmd)
}
