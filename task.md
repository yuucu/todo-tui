
# memo
思考時間: 1m 42s
以下では todo.txt フォーマットでタスクを保存し、lazygit風の操作感をもつ Go 製 TUI アプリを作るための設計指針とサンプルコードの最短ルートをまとめます。最終的なゴールは「1 ファイルですべてのタスクが完結し、キーボードだけで高速に一覧・フィルタ・編集ができる」環境です。

todo.txt フォーマットの要点
各タスクは 1 行のプレーンテキストで、先頭に (A) ～ (Z) の優先度、+project や @context、日付などを付与できる 。

フィールドは順不同なので、「タスクを文字列で持たせつつ正規表現で解析する」より「既存ライブラリを使う」ほうが安全です。Go なら github.com/1set/todotxt が多数のフィルタ／ソート API を備えています 。

ストレージ層：パーサ & ファイル I/O
go
コピーする
編集する
import todo "github.com/1set/todotxt"

func load(path string) (*todo.TaskList, error) {
    return todo.LoadFromPath(path) // パース＋バリデーション
}

func save(list *todo.TaskList, path string) error {
    return list.WriteToPath(path)
}
上記ライブラリは「未完了だけ抽出」「今日期限のものだけ抽出」など典型的なクエリをワンライナーで書けるので、ビジネスロジックがほぼ DSL になります 。

UI フレームワーク選定
候補	特徴	lazygit との親和性
Bubble Tea	Elm Architecture 準拠・モデル/更新/ビューを分離。Bubbles/Lipgloss がリッチな部品と装飾を提供	lazygit（旧 gocui 製）に近いキーバインドと画面分割を Lipgloss で柔軟に再現可能
tview	ウィジェットが豊富で学習コスト低	決め打ちレイアウトが多く、lazygit 的「パネルを矢印で行き来」は実装コスト高
gocui	最低限の View とキーバインドだけ提供。lazygit が採用	同じ API なので挙動は最も近いが、部品が少なく自前実装が増える

結論：Bubble Tea を第一候補

list・textarea などの既成コンポーネントで 8 割を書ける 。

Lipgloss でフレームや色を付けると “lazygit 風” の境界線やヘッダーが数行で引ける 。

Bubble Tea 製 ToDo サンプルも公開されており学習資料が揃う 。

アーキテクチャ概要
ディレクトリ構成例
csharp
コピーする
編集する
cmd/
  todos/
    main.go        // CLI 起動
internal/
  ui/              // Bubble Tea モデル群
  todo/            // ドメインロジック (wrap todotxt lib)
  infra/           // ファイルウォッチャなど
Bubble Tea の 3 分割
レイヤ	責務
Model	TaskList と UI 状態（選択インデックス、現在のフィルタ）
Update	キー入力を受け取り、リスト操作／ファイル書込を実行
View	Pane を Lipgloss で分割：
左=未完了、右=完了、下=ヘルプ

キーバインド案（lazygit に寄せる）
Key	動作
j/k	カーソル上下
Tab	パネル切替
c	新規タスク入力
e	選択行を編集
x	完了トグル
d	削除
p	優先度変更 (A→Z cycle)
/	インクリメンタルフィルタ

実装ステップ
todo.txt 読み込み／保存 – 上記ラッパを書いてテスト駆動で安定化。

最小 UI – list.New() で未完了タスクだけ描画し j/k で移動 。

CRUD 操作 – Update 内で task.Complete() などライブラリ関数を呼び、都度 save()。

フィルタ／検索 – tasks.Filter(todo.FilterByProject("+proj")) などをキーに割当 。

マルチペイン化 – lipgloss.JoinHorizontal/Vertical を使い、パネルごとに list.Model を持たせる 。

ファイルウォッチ – fsnotify で外部変更を検知→再読み込み。

拡張 – 期日が近いタスクをトップに pin する、通知用に go-todo CLI を呼ぶ等 。

参考スニペット（Bubble Tea 最小例）
go
コピーする
編集する
type model struct {
    list list.Model
}

func initialModel(tasks []list.Item) model {
    l := list.New(tasks, list.NewDefaultDelegate(), 0, 0)
    l.Title = "TODO"
    return model{list: l}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    }
    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m model) View() string { return lipgloss.NewStyle().Padding(1).Render(m.list.View()) }
わずか 80 行ほどで「読み込み → 画面表示 → キー操作 → 保存」まで形になります。

まとめ & 次のステップ
todo.txt はプレーンテキストゆえ Git で履歴管理が容易 。

Bubble Tea + Lipgloss なら、lazygit 級の操作感を Go らしい型安全さとテスト容易性で構築できる。

まずは 1 ペイン + CRUD で最小動作を確認し、段階的にフィルタ／マルチペイン／通知を足すと安全。

これで “テキスト 1 枚、操作はキーボード完結、Git で履歴も追える” 開発者向け TODO 管理ツールの骨格が整いました。Happy hacking! 🎉
