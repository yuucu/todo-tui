package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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

	// Ensure absolute minimum terminal size
	minTerminalWidth := 40
	minTerminalHeight := 8

	if m.width < minTerminalWidth {
		m.width = minTerminalWidth
	}
	if m.height < minTerminalHeight {
		m.height = minTerminalHeight
	}

	// Calculate pane sizes using configuration
	borderWidth := 4
	availableWidth := m.width - borderWidth
	leftWidth := int(float64(availableWidth) * m.appConfig.UI.LeftPaneRatio)
	rightWidth := availableWidth - leftWidth

	// Ensure minimum widths from configuration
	if leftWidth < m.appConfig.UI.MinLeftPaneWidth {
		leftWidth = m.appConfig.UI.MinLeftPaneWidth
		rightWidth = availableWidth - leftWidth
	}
	if rightWidth < m.appConfig.UI.MinRightPaneWidth {
		rightWidth = m.appConfig.UI.MinRightPaneWidth
	}

	// Calculate available height for list content (consistent with renderMainView)
	// Reserve space for:
	// - Combined help/status bar (1 line)
	// - Configurable vertical padding
	helpBarHeight := 1
	verticalPadding := m.appConfig.UI.VerticalPadding

	// Available height for the entire content area (including borders and titles)
	contentHeight := m.height - helpBarHeight - verticalPadding

	// Ensure minimum content height
	if contentHeight < 5 { // Minimum 5 lines to show border + some content
		contentHeight = 5
	}

	// Reserve space for border (2 lines) and title (1 line) within content area
	listHeight := contentHeight - 3 // 2 for borders + 1 for title

	if listHeight < 1 { // Ensure at least 1 line for list content
		listHeight = 1
	}

	// Set the calculated height for both lists
	m.filterList.SetHeight(listHeight)
	m.taskList.SetHeight(listHeight)
}

// View renders the UI
func (m *Model) View() string {
	// Ensure pane sizes are set only if not initialized
	if m.width <= 0 || m.height <= 0 {
		m.updatePaneSizes()
	}

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
	// Calculate dimensions for panels using configuration
	borderWidth := 4
	availableWidth := m.width - borderWidth
	leftWidth := int(float64(availableWidth) * m.appConfig.UI.LeftPaneRatio)
	rightWidth := availableWidth - leftWidth

	// Ensure minimum widths from configuration
	if leftWidth < m.appConfig.UI.MinLeftPaneWidth {
		leftWidth = m.appConfig.UI.MinLeftPaneWidth
		rightWidth = availableWidth - leftWidth
	}
	if rightWidth < m.appConfig.UI.MinRightPaneWidth {
		rightWidth = m.appConfig.UI.MinRightPaneWidth
	}

	// Calculate content height (consistent with updatePaneSizes)
	helpBarHeight := 1
	verticalPadding := m.appConfig.UI.VerticalPadding
	contentHeight := m.height - helpBarHeight - verticalPadding

	// Ensure minimum content height
	if contentHeight < 5 { // Minimum 5 lines to show border + some content
		contentHeight = 5
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

	// Create combined help/status bar
	combinedBar := m.renderCombinedHelpStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left, content, combinedBar)
}

// renderCombinedHelpStatusBar creates a single bar with help text on the left and status on the right
func (m *Model) renderCombinedHelpStatusBar() string {
	// Get help text based on active pane and current filter
	var helpText string
	if m.activePane == paneFilter {
		helpText = "j/k: navigate | Enter: select filter & move to tasks | Tab/h/l: switch panes | a: add | q: quit"
	} else {
		// Check if we're viewing deleted tasks
		isViewingDeleted := m.filterList.selected < len(m.filters) && m.filters[m.filterList.selected].name == deletedTasksFilter
		if isViewingDeleted {
			helpText = "j/k: navigate | r: restore task | Tab/h/l: switch panes | a: add | q: quit"
		} else {
			helpText = "j/k: navigate | Enter: complete task | e: edit | p: priority toggle | d: delete | Tab/h/l: switch panes | a: add | q: quit"
		}
	}

	// Get status information
	statusText := m.getStatusInfo()

	// Calculate available width for help text
	statusWidth := lipgloss.Width(statusText)
	helpAvailableWidth := m.width - statusWidth - 2 // Leave some space between them

	// Truncate help text if necessary
	if lipgloss.Width(helpText) > helpAvailableWidth {
		truncateLen := helpAvailableWidth - 3 // Account for "..."
		if truncateLen > 0 {
			runes := []rune(helpText)
			if len(runes) > truncateLen {
				helpText = string(runes[:truncateLen]) + "..."
			}
		}
	}

	// Create styles for left and right parts
	leftStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Text).
		Background(m.currentTheme.Background)

	rightStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.TextMuted).
		Background(m.currentTheme.Background)

	// Create the combined bar with proper spacing
	leftPart := leftStyle.Render(helpText)
	rightPart := rightStyle.Render(statusText)

	// Calculate spacing needed
	usedWidth := lipgloss.Width(leftPart) + lipgloss.Width(rightPart)
	spacingNeeded := m.width - usedWidth
	if spacingNeeded < 0 {
		spacingNeeded = 0
	}

	spacing := strings.Repeat(" ", spacingNeeded)

	// Combine with spacing
	combinedContent := leftPart + spacing + rightPart

	// Apply background style to the entire bar
	barStyle := lipgloss.NewStyle().
		Background(m.currentTheme.Background).
		Width(m.width)

	return barStyle.Render(combinedContent)
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
