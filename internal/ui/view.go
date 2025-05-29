package ui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// updatePaneSizes calculates and applies proper sizes to both panes
func (m *Model) updatePaneSizes() {
	// Set minimum dimensions if not yet initialized
	if m.width <= 0 {
		m.width = 80 // Default terminal width
	}
	if m.height <= 0 {
		m.height = 24 // Default terminal height
	}

	// Calculate pane sizes (1/3 left, 2/3 right like lazygit)
	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth

	// Ensure minimum widths
	if leftWidth < 20 {
		leftWidth = 20
		rightWidth = m.width - leftWidth
	}
	if rightWidth < 30 {
		rightWidth = 30
		leftWidth = m.width - rightWidth
	}

	// Reserve space for:
	// - Custom title (1 line each pane)
	// - Main help bar (1 line)
	// - Status bar (1 line)
	// - Borders around panes (2 lines total for top/bottom)
	titleHeight := 1
	mainHelpHeight := 1
	statusBarHeight := 1
	verticalBorderHeight := 2

	availableHeight := m.height - titleHeight - mainHelpHeight - statusBarHeight - verticalBorderHeight // Remove the +2 to ensure proper display

	if availableHeight <= 2 { // Ensure at least 2 lines for list content
		availableHeight = 2
	}

	// Set the calculated height for both lists
	m.filterList.SetHeight(availableHeight)
	m.taskList.SetHeight(availableHeight)
}

// View renders the UI
func (m *Model) View() string {
	// Ensure pane sizes are set
	m.updatePaneSizes()

	// If in add/edit mode, show textarea (keeping existing behavior as full screen)
	if m.currentMode == modeAdd || m.currentMode == modeEdit {
		title := "Add New Task"
		if m.currentMode == modeEdit {
			title = "Edit Task"
		}

		titleStyle := lipgloss.NewStyle().
			Foreground(m.currentTheme.Primary).
			Bold(true).
			Padding(0, 1)

		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(m.currentTheme.Primary).
			Padding(0, 1)

		helpStyle := lipgloss.NewStyle().
			Foreground(m.currentTheme.TextSubtle).
			Padding(0, 1)

		return lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(title),
			inputStyle.Render(m.textarea.View()),
			helpStyle.Render("Enter/Ctrl+S: 保存 | Esc/Ctrl+C: キャンセル"),
		)
	}

	// Always render the main view (panes, help bar, etc.)
	mainView := m.renderMainView()

	// If in delete confirmation mode, overlay the dialog on top of the main view
	if m.currentMode == modeDeleteConfirm {
		if m.deleteIndex < len(m.filteredTasks) {
			task := &m.filteredTasks[m.deleteIndex]

			// Create dialog content
			dialogTitleStyle := lipgloss.NewStyle().
				Foreground(m.currentTheme.Danger).
				Bold(true)

			taskStyle := lipgloss.NewStyle().
				Foreground(m.currentTheme.TextMuted).
				Italic(true)

			dialogHelpStyle := lipgloss.NewStyle().
				Foreground(m.currentTheme.TextSubtle)

			dialogContent := lipgloss.JoinVertical(lipgloss.Left,
				dialogTitleStyle.Render("Delete Task?"),
				"",
				taskStyle.Render(task.Todo),
				"",
				dialogHelpStyle.Render("y: Yes, delete | n: No, cancel | Esc: Cancel"),
			)

			// Style the dialog with background to make it opaque
			dialogStyle := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(m.currentTheme.Danger).
				Padding(1, 2).
				Background(m.currentTheme.Surface).
				Foreground(m.currentTheme.Text)

			dialogBox := dialogStyle.Render(dialogContent)

			// Overlay the dialog on the main view
			return m.overlayDialog(mainView, dialogBox)
		}
	}

	return mainView
}

// renderMainView renders the main application view (panes, help bar, status bar)
func (m *Model) renderMainView() string {
	// Calculate dimensions for panels
	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth

	// Ensure minimum widths
	if leftWidth < 20 {
		leftWidth = 20
		rightWidth = m.width - leftWidth
	}
	if rightWidth < 30 {
		rightWidth = 30
		leftWidth = m.width - rightWidth
	}

	// Calculate content height (reserve space for help and status bar)
	contentHeight := m.height - 4 // Reserve for help bar, status bar, and ensure top elements are visible
	
	// Ensure minimum content height
	if contentHeight < 3 {
		contentHeight = 3
	}

	// Define styles for the panels (add height back to use full terminal height)
	activeBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.currentTheme.BorderActive).
		Width(leftWidth).
		Height(contentHeight)

	inactiveBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.currentTheme.BorderInactive).
		Width(leftWidth).
		Height(contentHeight)

	activeRightBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.currentTheme.BorderActive).
		Width(rightWidth).
		Height(contentHeight)

	inactiveRightBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.currentTheme.BorderInactive).
		Width(rightWidth).
		Height(contentHeight)

	var leftPane, rightPane string

	// Create custom titles
	titleStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Primary).
		Bold(true)

	filterTitle := titleStyle.Render("Workspaces")
	taskTitle := titleStyle.Render("Todos")

	if m.activePane == paneFilter {
		leftPaneContent := lipgloss.JoinVertical(lipgloss.Left, filterTitle, m.filterList.View())
		rightPaneContent := lipgloss.JoinVertical(lipgloss.Left, taskTitle, m.taskList.View())
		leftPane = activeBorderStyle.Render(leftPaneContent)
		rightPane = inactiveRightBorderStyle.Render(rightPaneContent)
	} else {
		leftPaneContent := lipgloss.JoinVertical(lipgloss.Left, filterTitle, m.filterList.View())
		rightPaneContent := lipgloss.JoinVertical(lipgloss.Left, taskTitle, m.taskList.View())
		leftPane = inactiveBorderStyle.Render(leftPaneContent)
		rightPane = activeRightBorderStyle.Render(rightPaneContent)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	// Help bar with theme colors - different help based on active pane and current filter
	var helpText string
	if m.activePane == paneFilter {
		helpText = "j/k: navigate | Enter: select filter & move to tasks | Tab/h/l: switch panes | a: add | q: quit"
	} else {
		// Check if we're viewing deleted tasks
		isViewingDeleted := m.filterList.selected < len(m.filters) && m.filters[m.filterList.selected].name == "Deleted Tasks"
		if isViewingDeleted {
			helpText = "j/k: navigate | r: restore task | Tab/h/l: switch panes | a: add | q: quit"
		} else {
			helpText = "j/k: navigate | Enter: complete task | e: edit | p: priority toggle | d: delete | Tab/h/l: switch panes | a: add | q: quit"
		}
	}

	help := lipgloss.NewStyle().
		Foreground(m.currentTheme.Text).
		Background(m.currentTheme.Background).
		Width(m.width).
		Render(helpText)

	// Status bar with actual information
	statusBar := lipgloss.NewStyle().
		Foreground(m.currentTheme.TextMuted).
		Background(m.currentTheme.Background).
		Width(m.width).
		Render(m.getStatusInfo())

	return lipgloss.JoinVertical(lipgloss.Left, content, help, statusBar)
}

// overlayDialog overlays a dialog box on top of the main view
func (m *Model) overlayDialog(mainView, dialog string) string {
	// First, create a full-screen layout with the dialog centered
	dialogOverlay := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialog)
	
	// Split both views into lines
	mainLines := strings.Split(mainView, "\n")
	dialogLines := strings.Split(dialogOverlay, "\n")
	
	// Ensure both have enough lines
	maxLines := m.height
	if len(mainLines) > maxLines {
		maxLines = len(mainLines)
	}
	if len(dialogLines) > maxLines {
		maxLines = len(dialogLines)
	}
	
	// Pad lines to match screen height
	for len(mainLines) < maxLines {
		mainLines = append(mainLines, "")
	}
	for len(dialogLines) < maxLines {
		dialogLines = append(dialogLines, "")
	}
	
	result := make([]string, maxLines)
	
	// Combine the views: use dialog overlay where it has non-space content
	for i := 0; i < maxLines; i++ {
		mainLine := mainLines[i]
		dialogLine := dialogLines[i]
		
		// If dialog line has any non-space content, use it; otherwise use main line
		hasDialogContent := false
		for _, r := range dialogLine {
			if r != ' ' && r != '\t' {
				hasDialogContent = true
				break
			}
		}
		
		if hasDialogContent {
			result[i] = dialogLine
		} else {
			result[i] = mainLine
		}
	}
	
	return strings.Join(result, "\n")
} 