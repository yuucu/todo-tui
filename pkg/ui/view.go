package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// calculatePaneWidths calculates left and right pane widths based on configuration
func (m *Model) calculatePaneWidths(availableWidth int) (int, int) {
	leftWidth := int(float64(availableWidth) * m.appConfig.UI.LeftPaneRatio)
	rightWidth := availableWidth - leftWidth

	// Ensure minimum widths from configuration
	minLeftWidth := m.appConfig.UI.MinLeftPaneWidth
	minRightWidth := m.appConfig.UI.MinRightPaneWidth

	// If minimum widths exceed available space, scale them down proportionally
	totalMinWidth := minLeftWidth + minRightWidth
	if totalMinWidth > availableWidth {
		leftWidth = int(float64(availableWidth) * float64(minLeftWidth) / float64(totalMinWidth))
		rightWidth = availableWidth - leftWidth
	} else {
		// Adjust widths to meet minimum requirements
		if leftWidth < minLeftWidth {
			leftWidth = minLeftWidth
		}
		if rightWidth < minRightWidth {
			rightWidth = minRightWidth
		}
		// Recalculate to ensure total doesn't exceed available width
		if leftWidth+rightWidth > availableWidth {
			// If both minimums can't fit, scale proportionally
			leftWidth = int(float64(availableWidth) * float64(minLeftWidth) / float64(totalMinWidth))
			rightWidth = availableWidth - leftWidth
		}
	}

	return leftWidth, rightWidth
}

// updatePaneSizes calculates and sets the sizes for UI panes
func (m *Model) updatePaneSizes() {
	// Set minimum dimensions if not yet initialized
	if m.width <= 0 {
		m.width = DefaultTerminalWidth // Default terminal width
	}
	if m.height <= 0 {
		m.height = DefaultTerminalHeight // Default terminal height
	}

	// Use actual terminal size for calculations - NEVER exceed actual size
	actualWidth := m.width
	actualHeight := m.height

	// Calculate pane sizes using configuration
	borderWidth := PaneBorderWidth
	availableWidth := actualWidth - borderWidth

	// Ensure minimum viable width but never exceed actual width
	if availableWidth < MinimumAvailableWidth {
		availableWidth = MinimumAvailableWidth
		if availableWidth > actualWidth-PaneBorderWidth {
			availableWidth = actualWidth - PaneBorderWidth
		}
	}

	// Calculate pane widths using the common function
	leftWidth, rightWidth := m.calculatePaneWidths(availableWidth)

	// Store calculated widths for potential future use
	_ = leftWidth
	_ = rightWidth

	// Calculate available height for list content using actual terminal size
	// Reserve space for:
	// - Combined help/status bar (1 line)
	// - Configurable vertical padding
	helpBarHeight := HelpStatusBarHeight
	verticalPadding := m.appConfig.UI.VerticalPadding

	// Ensure padding doesn't exceed reasonable limits for small terminals
	if verticalPadding > actualHeight/3 {
		verticalPadding = actualHeight / 3
	}

	// Available height for the entire content area using actual terminal size
	contentHeight := actualHeight - helpBarHeight - verticalPadding

	// Ensure we have at least minimal content height
	if contentHeight < MinimumContentHeight {
		contentHeight = MinimumContentHeight
	}

	// Reserve space for border (2 lines) and title (1 line) within content area
	listHeight := contentHeight - ListContentReserved // 2 for borders + 1 for title

	if listHeight < MinimumListHeight { // Ensure at least 1 line for list content
		listHeight = MinimumListHeight
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

	// If in help mode, show help screen
	if m.currentMode == modeHelp {
		return m.renderHelpView()
	}

	// If in add/edit mode, show textarea (keeping existing behavior as full screen)
	if m.currentMode == modeAdd || m.currentMode == modeEdit {
		title := AddTaskTitle
		if m.currentMode == modeEdit {
			title = EditTaskTitle
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
			helpStyle.Render(EditModeHelp),
		)
	}

	// Always render the main view (panes, help bar, etc.)
	mainView := m.renderMainView()

	return mainView
}

// renderMainView renders the main application view (panes, help bar, status bar)
func (m *Model) renderMainView() string {
	// Use actual terminal size - NEVER exceed actual dimensions
	actualWidth := m.width
	actualHeight := m.height

	// Calculate dimensions for panels using configuration
	borderWidth := PaneBorderWidth
	availableWidth := actualWidth - borderWidth

	// Ensure minimum viable width but never exceed actual width
	if availableWidth < MinimumAvailableWidth {
		availableWidth = MinimumAvailableWidth
		if availableWidth > actualWidth-PaneBorderWidth {
			availableWidth = actualWidth - PaneBorderWidth
		}
	}

	// Calculate pane widths using the common function
	leftWidth, rightWidth := m.calculatePaneWidths(availableWidth)

	// Calculate content height using actual terminal size
	helpBarHeight := HelpStatusBarHeight
	verticalPadding := m.appConfig.UI.VerticalPadding

	// Ensure padding doesn't exceed reasonable limits for small terminals
	if verticalPadding > actualHeight/3 {
		verticalPadding = actualHeight / 3
	}

	contentHeight := actualHeight - helpBarHeight - verticalPadding

	// Ensure we have at least minimal content height
	if contentHeight < MinimumContentHeight {
		contentHeight = MinimumContentHeight
	}

	// Define styles for the panels (strictly use calculated content height)
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

	filterTitle := titleStyle.Render(FilterPaneTitle)
	taskTitle := titleStyle.Render(TaskPaneTitle)

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
	// Ensure we have a valid width
	if m.width <= 0 {
		return ""
	}

	// Get help text based on active pane and current filter
	var helpText string
	if m.activePane == paneFilter {
		helpText = HelpFilterPane
	} else {
		// Check if we're viewing deleted tasks
		isViewingDeleted := m.filterList.selected < len(m.filters) && m.filters[m.filterList.selected].name == FilterDeletedTasks
		isViewingCompleted := m.filterList.selected < len(m.filters) && m.filters[m.filterList.selected].name == FilterCompletedTasks

		if isViewingDeleted {
			helpText = "j/k: navigate | r: restore task | y: copy task | Tab/h/l: switch panes | a: add | q: quit"
		} else if isViewingCompleted {
			helpText = "j/k: navigate | r: restore task | y: copy task | Tab/h/l: switch panes | a: add | q: quit"
		} else {
			helpText = "j/k: navigate | Enter: toggle completion | e: edit | p: priority toggle | d: delete | y: copy task | Tab/h/l: switch panes | a: add | q: quit"
		}
	}

	// Get status information
	statusText := m.getStatusInfo()

	// Calculate available width for help text, ensuring we don't exceed terminal width
	statusWidth := lipgloss.Width(statusText)

	// Reserve space for status text and some padding
	reservedWidth := statusWidth + StatusTextSpacing // 2 spaces for padding
	helpAvailableWidth := m.width - reservedWidth

	// Ensure we have at least some space for help text
	if helpAvailableWidth < MinimumHelpTextWidth {
		// If terminal is too narrow, prioritize status and truncate help heavily
		helpAvailableWidth = MinimumHelpTextWidth
		if m.width < MinimumTerminalWidth {
			helpAvailableWidth = m.width / 2
		}
	}

	// Truncate help text if necessary
	if lipgloss.Width(helpText) > helpAvailableWidth {
		truncateLen := helpAvailableWidth - EllipsisLength // Account for "..."
		if truncateLen > 0 {
			runes := []rune(helpText)
			if len(runes) > truncateLen {
				helpText = string(runes[:truncateLen]) + Ellipsis
			}
		} else {
			helpText = Ellipsis // Minimal text if extremely narrow
		}
	}

	// Recalculate status width in case it was too long
	if statusWidth > m.width/2 {
		// Truncate status text if it's taking up too much space
		statusRunes := []rune(statusText)
		maxStatusWidth := m.width / 2
		if len(statusRunes)*WidthEstimateMultiplier > maxStatusWidth { // Rough estimate for character width
			truncateLen := maxStatusWidth/WidthEstimateMultiplier - EllipsisLength
			if truncateLen > 0 {
				statusText = string(statusRunes[:truncateLen]) + Ellipsis
			} else {
				statusText = Ellipsis
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

	// Calculate spacing needed, ensuring total width doesn't exceed terminal width
	usedWidth := lipgloss.Width(leftPart) + lipgloss.Width(rightPart)
	spacingNeeded := m.width - usedWidth
	if spacingNeeded < 0 {
		spacingNeeded = 0
	}

	spacing := strings.Repeat(" ", spacingNeeded)

	// Combine with spacing
	combinedContent := leftPart + spacing + rightPart

	// Ensure the final content doesn't exceed terminal width
	if lipgloss.Width(combinedContent) > m.width {
		// Final fallback: truncate the entire content
		runes := []rune(combinedContent)
		if len(runes) > m.width {
			combinedContent = string(runes[:m.width-3]) + "..."
		}
	}

	// Apply background style to the entire bar with strict width control
	barStyle := lipgloss.NewStyle().
		Background(m.currentTheme.Background).
		Width(m.width).
		MaxWidth(m.width) // Ensure we never exceed the width

	return barStyle.Render(combinedContent)
}
