package ui

// ===============================
// ターミナル・UI関連定数
// ===============================

// Default terminal size
const (
	DefaultTerminalWidth  = 80
	DefaultTerminalHeight = 24
)

// UI layout constants
const (
	// Border width between panes
	PaneBorderWidth = 4

	// Height of help/status bar
	HelpStatusBarHeight = 1

	// Minimum width/height constraints
	MinimumAvailableWidth = 20
	MinimumContentHeight  = 3
	MinimumListHeight     = 1
	MinimumHelpTextWidth  = 10
	MinimumTerminalWidth  = 20

	// Text area related
	TextAreaPadding   = 10
	TextAreaHeight    = 3
	TextAreaCharLimit = 0 // Unlimited

	// Padding & margin
	DefaultVerticalPadding = 2
	ListContentReserved    = 3 // 2 border lines + 1 title line

	// Text truncation
	EllipsisLength          = 3
	StatusTextSpacing       = 2
	WidthEstimateMultiplier = 2
)

// ===============================
// todo.txt Field Constants
// ===============================

// Task field names
const (
	TaskFieldDeleted = "deleted_at"
	TaskFieldDue     = "due"
)

// Task field prefixes (with colon)
const (
	TaskFieldDeletedPrefix = TaskFieldDeleted + ":"
	TaskFieldDuePrefix     = TaskFieldDue + ":"
)

// ===============================
// Date & Time Format Constants
// ===============================

const (
	// Go standard date format (ISO 8601)
	DateFormat = "2006-01-02"

	// Time display format
	TimeFormat = "15:04"
)

// ===============================
// Filter Name Constants
// ===============================

const (
	FilterAllTasks       = "All Tasks"
	FilterCompletedTasks = "Completed Tasks"
	FilterDeletedTasks   = "Deleted Tasks"
	FilterNoProject      = "No Project"

	// Section headers
	FilterHeaderProjects = "── Projects ──"
	FilterHeaderContexts = "── Contexts ──"
)

// ===============================
// UI Message & Placeholder Constants
// ===============================

const (
	// Text area placeholder
	TaskInputPlaceholder = "Enter task description (e.g., 'call @mom +home due:2025-01-15')"

	// Edit mode titles
	AddTaskTitle  = "Add New Task"
	EditTaskTitle = "Edit Task"

	// Edit mode help
	EditModeHelp = "Enter/Ctrl+S: save | Esc/Ctrl+C: cancel"

	// Help text
	HelpFilterPane      = "?: help | j/k: navigate | Enter: select filter & move to tasks | Tab/h/l: switch panes | a: add | q: quit"
	HelpTaskPane        = "?: help | j/k: navigate | Enter: toggle completion | e: edit | p: priority toggle | t: toggle due today | d: delete | y: copy task | Tab/h/l: switch panes | a: add | q: quit"
	HelpDeletedTaskPane = "?: help | j/k: navigate | r: restore task | y: copy task | Tab/h/l: switch panes | a: add | q: quit"

	// Panel titles
	FilterPaneTitle = "Workspaces"
	TaskPaneTitle   = "Todos"

	// Ellipsis symbol
	Ellipsis = "..."
)

// ===============================
// Configuration Default Value Constants
// ===============================

const (
	// UI configuration default values
	DefaultLeftPaneRatio     = 0.33
	DefaultMinLeftPaneWidth  = 18
	DefaultMinRightPaneWidth = 28

	// File permissions
	DefaultConfigDirMode = 0755
	DefaultFileDirMode   = 0755
)

// ===============================
// Other Constants
// ===============================

const (
	// Invalid index
	InvalidIndex = -1

	// Start index
	StartIndex = 0
)

// ===============================
// Key String Constants
// ===============================

// よく使用されるキー文字列定数
const (
	ctrlCKey = "ctrl+c"
	ctrlSKey = "ctrl+s"
	escKey   = "esc"
	enterKey = "enter"
	tabKey   = "tab"
	spaceKey = " "

	// Direction keys
	upKey    = "up"
	downKey  = "down"
	leftKey  = "left"
	rightKey = "right"

	// Letter keys for navigation
	jKey = "j"
	kKey = "k"
	hKey = "h"
	lKey = "l"
	gKey = "g"
	GKey = "G"

	// Action keys
	qKey = "q"
	aKey = "a"
	eKey = "e"
	dKey = "d"
	pKey = "p"
	tKey = "t"
	rKey = "r"
	yKey = "y"

	// Help key
	helpKey = "?"
)

// ===============================
// Checkbox Style Constants
// ===============================

// Modern checkbox styles
const (
	CheckboxStyleCircle  = "circle"  // ● / ○
	CheckboxStyleSquare  = "square"  // ■ / □
	CheckboxStyleCheck   = "check"   // ✓ / ○
	CheckboxStyleDiamond = "diamond" // ◆ / ◇
	CheckboxStyleStar    = "star"    // ★ / ☆
)

// Default checkbox style
const DefaultCheckboxStyle = CheckboxStyleCircle
