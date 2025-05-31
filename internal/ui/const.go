package ui

// ===============================
// ターミナル・UI関連定数
// ===============================

// デフォルトターミナルサイズ
const (
	DefaultTerminalWidth  = 80
	DefaultTerminalHeight = 24
)

// UI レイアウト定数
const (
	// パネル間のボーダー幅
	PaneBorderWidth = 4

	// ヘルプ・ステータスバーの高さ
	HelpStatusBarHeight = 1

	// 最小幅/高さ制約
	MinimumAvailableWidth = 20
	MinimumContentHeight  = 3
	MinimumListHeight     = 1
	MinimumHelpTextWidth  = 10
	MinimumTerminalWidth  = 20

	// テキストエリア関連
	TextAreaPadding   = 10
	TextAreaHeight    = 3
	TextAreaCharLimit = 0 // 無制限

	// パディング・マージン
	DefaultVerticalPadding = 2
	ListContentReserved    = 3 // ボーダー2行 + タイトル1行

	// 文字列切り詰め用
	EllipsisLength          = 3
	StatusTextSpacing       = 2
	WidthEstimateMultiplier = 2
)

// ===============================
// todo.txt フィールド定数
// ===============================

// タスクフィールド名
const (
	TaskFieldDeleted = "deleted_at"
	TaskFieldDue     = "due"
)

// タスクフィールドプレフィックス（コロン付き）
const (
	TaskFieldDeletedPrefix = TaskFieldDeleted + ":"
	TaskFieldDuePrefix     = TaskFieldDue + ":"
)

// ===============================
// 日付・時刻フォーマット定数
// ===============================

const (
	// Go標準の日付フォーマット（ISO 8601）
	DateFormat = "2006-01-02"

	// 時刻表示フォーマット
	TimeFormat = "15:04"
)

// ===============================
// フィルター名定数
// ===============================

const (
	FilterAllTasks       = "All Tasks"
	FilterCompletedTasks = "Completed Tasks"
	FilterDeletedTasks   = "Deleted Tasks"

	// セクションヘッダー
	FilterHeaderProjects = "── Projects ──"
	FilterHeaderContexts = "── Contexts ──"
)

// ===============================
// UI メッセージ・プレースホルダー定数
// ===============================

const (
	// テキストエリアプレースホルダー
	TaskInputPlaceholder = "タスクの説明を入力してください (例: '電話 @母 +home due:2025-01-15')"

	// 削除確認ダイアログ
	DeleteConfirmTitle = "Delete Task?"
	DeleteConfirmHelp  = "y: Yes, delete | n: No, cancel | Esc: Cancel"

	// 編集モードタイトル
	AddTaskTitle  = "Add New Task"
	EditTaskTitle = "Edit Task"

	// 編集モードヘルプ
	EditModeHelp = "Enter/Ctrl+S: 保存 | Esc/Ctrl+C: キャンセル"

	// ヘルプテキスト
	HelpFilterPane      = "?: help | j/k: navigate | Enter: select filter & move to tasks | Tab/h/l: switch panes | a: add | q: quit"
	HelpTaskPane        = "?: help | j/k: navigate | Enter: complete task | e: edit | p: priority toggle | t: toggle due today | d: delete | Tab/h/l: switch panes | a: add | q: quit"
	HelpDeletedTaskPane = "?: help | j/k: navigate | r: restore task | Tab/h/l: switch panes | a: add | q: quit"

	// パネルタイトル
	FilterPaneTitle = "Workspaces"
	TaskPaneTitle   = "Todos"

	// 省略記号
	Ellipsis = "..."
)

// ===============================
// 設定デフォルト値定数
// ===============================

const (
	// UI設定のデフォルト値
	DefaultLeftPaneRatio     = 0.33
	DefaultMinLeftPaneWidth  = 18
	DefaultMinRightPaneWidth = 28

	// ファイルパーミッション
	DefaultConfigDirMode = 0755
	DefaultFileDirMode   = 0755
)

// ===============================
// その他の定数
// ===============================

const (
	// 無効なインデックス
	InvalidIndex = -1

	// 開始インデックス
	StartIndex = 0

	// リストセレクター記号
	ListSelectorSymbol = "▶ "
	ListPaddingSymbol  = "  "
)
