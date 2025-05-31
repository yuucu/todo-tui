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
			Category: "基本操作",
			Items: []HelpItem{
				{"?", "このヘルプを表示"},
				{"q, Ctrl+C", "アプリケーションを終了"},
				{"a", "新しいタスクを追加"},
				{"e", "選択中のタスクを編集"},
				{"d", "選択中のタスクを削除"},
				{"r", "削除済み/完了済みタスクを復元"},
			},
		},
		{
			Category: "ナビゲーション",
			Items: []HelpItem{
				{"Tab", "ペイン間を切り替え"},
				{"h", "左ペイン（フィルタ）に移動"},
				{"l", "右ペイン（タスク）に移動"},
				{"j, ↓", "下に移動"},
				{"k, ↑", "上に移動"},
				{"Enter", "フィルタを適用/タスクを完了"},
			},
		},
		{
			Category: "タスク操作",
			Items: []HelpItem{
				{"p", "優先度を切り替え"},
				{"t", "期限を今日に設定/解除"},
			},
		},
		{
			Category: "入力モード",
			Items: []HelpItem{
				{"Ctrl+C", "入力をキャンセル"},
				{"Enter, Ctrl+S", "タスクを保存"},
			},
		},
		{
			Category: "削除確認モード",
			Items: []HelpItem{
				{"y, Y", "削除を実行"},
				{"n, N, Esc, Ctrl+C", "削除をキャンセル"},
			},
		},
	}
}

// renderHelpView renders the help screen
func (m *Model) renderHelpView() string {
	if len(m.helpContent) == 0 {
		m.initializeHelpContent()
	}

	var b strings.Builder

	// Define styles
	headerStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Primary).
		Bold(true).
		Padding(0, 1)

	categoryStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Secondary).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Warning).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.Text)

	mutedStyle := lipgloss.NewStyle().
		Foreground(m.currentTheme.TextMuted).
		Italic(true)

	// Header
	b.WriteString(headerStyle.Render("📚 キーバインディング ヘルプ"))
	b.WriteString("\n\n")

	// Content
	for i, category := range m.helpContent {
		if i > 0 {
			b.WriteString("\n")
		}

		// Category header
		b.WriteString(categoryStyle.Render("▶ " + category.Category))
		b.WriteString("\n")

		// Category items
		for _, item := range category.Items {
			keyPart := keyStyle.Render(item.Key)
			descPart := descStyle.Render(item.Description)
			b.WriteString("  " + keyPart + " : " + descPart + "\n")
		}
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Press any key to close this help..."))

	return b.String()
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
			Description: "ヘルプを表示",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				m.currentMode = modeHelp
				return m, nil
			},
		},
		{
			Key:         "q",
			Description: "アプリケーションを終了",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				return m, tea.Quit
			},
		},
		{
			Key:         ctrlCKey,
			Description: "アプリケーションを終了",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				return m, tea.Quit
			},
		},

		// Task management keys
		{
			Key:         "a",
			Description: "新しいタスクを追加",
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
			Description: "選択中のタスクを編集",
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
			Description: "ペイン間を切り替え",
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
			Description: "左ペイン（フィルタ）に移動",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				m.activePane = paneFilter
				return m, nil
			},
		},
		{
			Key:         "l",
			Description: "右ペイン（タスク）に移動",
			Mode:        modeView,
			Handler: func(m *Model) (tea.Model, tea.Cmd) {
				m.activePane = paneTask
				return m, nil
			},
		},
	}
}
