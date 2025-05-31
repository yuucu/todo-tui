package ui

import (
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
				{"j / ↓", "Move down"},
				{"k / ↑", "Move up"},
				{"Enter", "Apply filter / Complete task"},
			},
		},
		{
			Category: "Task Operations",
			Items: []HelpItem{
				{"p", "Cycle task priority"},
				{"t", "Toggle due date to today"},
			},
		},
		{
			Category: "Edit Mode",
			Items: []HelpItem{
				{"Ctrl+C", "Cancel editing"},
				{"Enter / Ctrl+S", "Save task"},
			},
		},
		{
			Category: "Delete Confirmation",
			Items: []HelpItem{
				{"y / Y", "Confirm deletion"},
				{"n / N / Esc", "Cancel deletion"},
			},
		},
	}
}

// renderHelpView renders the help screen
func (m *Model) renderHelpView() string {
	if len(m.helpContent) == 0 {
		m.initializeHelpContent()
	}

	var content strings.Builder

	// Define styles
	headerStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Primary).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(m.currentTheme.Primary)

	categoryStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Secondary).
		Bold(true).
		Underline(true).
		MarginTop(1).
		MarginBottom(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Warning).
		Bold(true).
		Width(16).
		Align(lipgloss.Right).
		Background(m.currentTheme.Surface).
		Padding(0, 1).
		MarginRight(1)

	descStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Text)

	// Header
	header := headerStyle.Render("📚 TodoTUI - Keyboard Shortcuts Help")
	content.WriteString(header)
	content.WriteString("\n\n")

	// Content sections
	for i, category := range m.helpContent {
		if i > 0 {
			content.WriteString("\n")
		}

		// Category header
		categoryHeader := categoryStyle.Render("▶ " + category.Category)
		content.WriteString(categoryHeader)
		content.WriteString("\n")

		// Category items in a box
		var itemsContent strings.Builder
		for _, item := range category.Items {
			keyPart := keyStyle.Render(item.Key)
			descPart := descStyle.Render(item.Description)
			itemsContent.WriteString(keyPart + " : " + descPart + "\n")
		}

		// Box style for items
		itemBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(m.currentTheme.BorderInactive).
			PaddingLeft(2).
			MarginLeft(1).
			MarginBottom(1)

		boxedItems := itemBoxStyle.Render(strings.TrimSuffix(itemsContent.String(), "\n"))
		content.WriteString(boxedItems)
		content.WriteString("\n")
	}

	// Footer with additional info
	footerContent := strings.Builder{}
	footerContent.WriteString("Todo.txt Format Examples:\n")
	footerContent.WriteString("  (A) Call Mom                   - High priority task\n")
	footerContent.WriteString("  Buy milk @store +groceries     - Task with context and project\n")
	footerContent.WriteString("  Meeting prep due:2025-05-31    - Task with due date\n")
	footerContent.WriteString("  x 2025-05-30 Completed task    - Completed task\n\n")
	footerContent.WriteString("Press any key to close this help...")

	footerStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.TextMuted).
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.currentTheme.BorderInactive).
		Padding(1).
		MarginTop(1)

	footer := footerStyle.Render(footerContent.String())
	content.WriteString(footer)

	// Center the entire help content
	helpContent := content.String()
	centeredStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Background(m.currentTheme.Background)

	return centeredStyle.Render(helpContent)
}

// KeyBinding represents a key binding configuration
type KeyBinding struct {
	Key         string
	Description string
	Mode        mode
	Pane        *pane // nil means any pane
	Handler     func(m *Model) (tea.Model, tea.Cmd)
}

// getKeyBindings returns all key bindings for the application
func (m *Model) getKeyBindings() []KeyBinding {
	return []KeyBinding{
		// Global keys
		{
			Key:         "?",
			Description: "Show help",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				m.currentMode = modeHelp
				return m, nil
			},
		},
		{
			Key:         "q",
			Description: "Quit application",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				return m, tea.Quit
			},
		},
		{
			Key:         ctrlCKey,
			Description: "Quit application",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				return m, tea.Quit
			},
		},

		// Task management keys
		{
			Key:         "a",
			Description: "Add new task",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				m.currentMode = modeAdd
				m.textarea.SetValue("")
				m.textarea.Focus()
				return m, nil
			},
		},
		{
			Key:         "e",
			Description: "Edit selected task",
			Mode:        modeView,
			Pane:        &[]pane{paneTask}[0],
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				if m.taskList.selected < len(m.filteredTasks) {
					m.currentMode = modeEdit
					selectedTask := m.filteredTasks[m.taskList.selected]
					for i := 0; i < len(m.tasks); i++ {
						if m.tasks[i].String() == selectedTask.String() {
							m.editingTask = &m.tasks[i]
							break
						}
					}
					if m.editingTask != nil {
						m.textarea.SetValue(m.editingTask.String())
						m.textarea.Focus()
					}
					return m, nil
				}
				return m, nil
			},
		},

		// Navigation keys
		{
			Key:         "tab",
			Description: "Switch between panes",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				if m.activePane == paneFilter {
					m.activePane = paneTask
				} else {
					m.activePane = paneFilter
				}
				return m, nil
			},
		},
		{
			Key:         "h",
			Description: "Move to left pane (Workspaces)",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				m.activePane = paneFilter
				return m, nil
			},
		},
		{
			Key:         "l",
			Description: "Move to right pane (Todos)",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				m.activePane = paneTask
				return m, nil
			},
		},
	}
}
