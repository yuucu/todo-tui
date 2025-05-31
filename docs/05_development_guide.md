# Todo TUI - 開発ガイド

## 1. 開発環境のセットアップ

### 1.1 必要要件

#### システム要件
- **Go**: 1.24以降
- **Git**: バージョン管理
- **Terminal**: カラー対応推奨（256色以上）

#### 推奨開発ツール
- **エディタ**: VS Code, GoLand, Vim/Neovim
- **デバッガ**: Delve (dlv)
- **リンター**: golangci-lint
- **フォーマッター**: gofmt, goimports

### 1.2 プロジェクトのクローンとビルド

```bash
# プロジェクトのクローン
git clone https://github.com/yuucu/todotui.git
cd todotui

# 依存関係のダウンロード
go mod tidy

# ビルド
make build

# テスト実行
make test

# 開発モードでの実行
go run cmd/todotui/main.go sample.todo.txt
```

### 1.3 Makefile解説

```makefile
# 主要なMakeターゲット
.PHONY: build test lint clean install

# バイナリビルド
build:
	go build -o bin/todotui cmd/todotui/main.go

# テスト実行
test:
	go test -v ./...

# リンター実行
lint:
	golangci-lint run

# クリーンアップ
clean:
	rm -rf bin/

# システムインストール
install:
	go install cmd/todotui/main.go
```

## 2. プロジェクト構造と設計原則

### 2.1 ディレクトリ構造

```
todotui/
├── cmd/
│   └── todotui/          # メインエントリーポイント
│       └── main.go       # CLI解析とアプリケーション起動
├── internal/
│   ├── ui/               # UIレイヤー
│   │   ├── model.go      # メインアプリケーションモデル
│   │   ├── view.go       # ビューレンダリング
│   │   ├── types.go      # 型定義
│   │   ├── config.go     # 設定管理
│   │   ├── filters.go    # フィルタリング機能
│   │   ├── colors.go     # テーマ・カラー管理
│   │   └── ime_helper.go # 日本語入力支援
│   └── todo/             # データレイヤー
│       ├── storage.go    # ファイル入出力
│       └── storage_test.go
├── docs/                 # ドキュメント
├── sample.todo.txt       # サンプルファイル
├── sample-config.yaml    # サンプル設定
├── go.mod               # Go モジュール定義
├── go.sum               # 依存関係ハッシュ
├── Makefile            # ビルドスクリプト
└── README.md           # プロジェクト概要
```

### 2.2 設計原則

#### Clean Architecture
```
┌─────────────────┐
│   Presentation  │ ← UI層（internal/ui/）
│   (UI Layer)    │
├─────────────────┤
│   Business      │ ← ビジネスロジック
│   Logic         │   （フィルタリング、状態管理）
├─────────────────┤
│   Data Access   │ ← データアクセス層
│   (Storage)     │   （internal/todo/）
└─────────────────┘
```

#### 依存性の方向
- **上位層 → 下位層**: UIはビジネスロジックに依存
- **抽象化**: インターフェースによる疎結合
- **単一責任**: 各モジュールは明確な責任を持つ

## 3. コード規約とスタイル

### 3.1 Goコーディング規約

#### 命名規則
```go
// パッケージ名: 小文字、短縮形
package ui

// 型名: PascalCase
type TaskManager struct {}

// メソッド名: PascalCase（公開）、camelCase（非公開）
func (tm *TaskManager) AddTask() {}
func (tm *TaskManager) updateStatus() {}

// 定数: PascalCase または UPPER_CASE
const DefaultTimeout = 5
const MAX_RETRY_COUNT = 3

// インターフェース名: 通常は-er suffix
type TaskRenderer interface {}
```

#### コメント規則
```go
// Package ui provides terminal user interface components
// for the Todo TUI application.
package ui

// Model represents the main application state and handles
// all user interactions and business logic.
type Model struct {
    // tasks holds the current list of todo items
    tasks todotxt.TaskList
    // activePane indicates which pane currently has focus
    activePane pane
}

// NewModel creates a new Model instance with the specified
// todo file and configuration.
func NewModel(todoFile string, config AppConfig) (*Model, error) {
    // 実装...
}
```

### 3.2 エラーハンドリング

#### エラー定義
```go
// エラー型の定義
type TodoError struct {
    Type    ErrorType
    Message string
    Cause   error
}

func (e *TodoError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

// カスタムエラーの作成
func NewFileError(path string, err error) error {
    return &TodoError{
        Type:    ErrorFileAccess,
        Message: fmt.Sprintf("failed to access file: %s", path),
        Cause:   err,
    }
}
```

#### エラーラッピング
```go
// Go 1.13以降の標準的なエラーラッピング
func (m *Model) loadTasks() error {
    tasks, err := todo.Load(m.todoFile)
    if err != nil {
        return fmt.Errorf("loading tasks from %s: %w", m.todoFile, err)
    }
    m.tasks = tasks
    return nil
}
```

### 3.3 テストの書き方

#### テストファイル構造
```go
// storage_test.go
package todo

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
    tests := []struct {
        name     string
        setup    func(t *testing.T) string  // テストデータのセットアップ
        want     int                        // 期待する結果
        wantErr  bool                       // エラーが期待されるか
    }{
        {
            name: "valid todo file",
            setup: func(t *testing.T) string {
                tmpfile := createTempFile(t, validTodoContent)
                return tmpfile
            },
            want:    3,
            wantErr: false,
        },
        {
            name: "non-existent file",
            setup: func(t *testing.T) string {
                return "/non/existent/path"
            },
            want:    0,
            wantErr: false, // ファイルが存在しない場合は自動作成
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            path := tt.setup(t)
            defer os.Remove(path)

            tasks, err := Load(path)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Len(t, tasks, tt.want)
        })
    }
}
```

#### テストヘルパー関数
```go
// テスト用のヘルパー関数
func createTempFile(t *testing.T, content string) string {
    tmpfile, err := os.CreateTemp("", "todo-test-*.txt")
    require.NoError(t, err)
    
    _, err = tmpfile.WriteString(content)
    require.NoError(t, err)
    
    err = tmpfile.Close()
    require.NoError(t, err)
    
    return tmpfile.Name()
}

const validTodoContent = `(A) Call Mom @phone +family
Buy milk @store +groceries  
x 2025-01-14 Clean garage @home +chores`
```

## 4. Bubble Teaアーキテクチャ

### 4.1 MVUパターンの実装

#### Model
```go
// Model は アプリケーションの状態を保持
type Model struct {
    // 状態
    filterList    SimpleList
    taskList      SimpleList
    textarea      textarea.Model
    
    // データ
    tasks         todotxt.TaskList
    filteredTasks todotxt.TaskList
    
    // UI状態
    activePane    pane
    currentMode   mode
    width, height int
    
    // 設定
    currentTheme  Theme
    appConfig     AppConfig
}

// Init は初期化コマンドを返す
func (m Model) Init() tea.Cmd {
    return tea.Batch(
        m.loadTasksCmd(),
        m.startFileWatcherCmd(),
    )
}
```

#### Update
```go
// Update はメッセージを受け取って新しい状態を返す
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    
    case FileChangedMsg:
        return m, m.loadTasksCmd()
    
    default:
        // 子コンポーネントのUpdate
        m.taskList, cmd = m.taskList.Update(msg)
        cmds = append(cmds, cmd)
    }

    return m, tea.Batch(cmds...)
}
```

#### View
```go
// View は現在の状態をレンダリング
func (m Model) View() string {
    if m.width == 0 || m.height == 0 {
        return "Initializing..."
    }

    switch m.currentMode {
    case modeAdd, modeEdit:
        return m.renderEditMode()
    case modeDeleteConfirm:
        return m.renderDeleteConfirm()
    default:
        return m.renderMainView()
    }
}
```

### 4.2 カスタムメッセージ

```go
// カスタムメッセージの定義
type FileChangedMsg struct{}
type TaskUpdatedMsg struct {
    Task *todotxt.Task
    Index int
}
type FilterChangedMsg struct {
    FilterName string
}

// メッセージを生成するコマンド
func (m *Model) loadTasksCmd() tea.Cmd {
    return func() tea.Msg {
        tasks, err := todo.Load(m.todoFile)
        if err != nil {
            return ErrorMsg{err}
        }
        return TasksLoadedMsg{tasks}
    }
}
```

### 4.3 並行処理

```go
// ファイル監視の実装
func (m *Model) startFileWatcher() tea.Cmd {
    return func() tea.Msg {
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
            return ErrorMsg{err}
        }
        
        go func() {
            for {
                select {
                case event := <-watcher.Events:
                    if event.Op&fsnotify.Write == fsnotify.Write {
                        // ファイル変更を通知
                        tea.Send(FileChangedMsg{})
                    }
                case err := <-watcher.Errors:
                    tea.Send(ErrorMsg{err})
                }
            }
        }()
        
        return WatcherStartedMsg{watcher}
    }
}
```

## 5. デバッグとプロファイリング

### 5.1 デバッグ設定

#### ログ設定
```go
// デバッグログの設定
func setupLogging(config AppConfig) {
    if config.Debug.Enabled {
        logFile, err := os.OpenFile(config.Debug.LogFile, 
            os.O_CREATE|os.O_WRITEL|os.O_APPEND, 0666)
        if err == nil {
            log.SetOutput(logFile)
        }
        
        switch config.Debug.LogLevel {
        case "debug":
            log.SetLevel(log.DebugLevel)
        case "info":
            log.SetLevel(log.InfoLevel)
        }
    }
}
```

#### デバッグビルド
```bash
# デバッグ情報付きビルド
go build -gcflags="all=-N -l" -o bin/todotui-debug cmd/todotui/main.go

# Delveでデバッグ実行
dlv exec ./bin/todotui-debug -- sample.todo.txt
```

### 5.2 パフォーマンステスト

#### ベンチマークテスト
```go
// フィルタリング性能のベンチマーク
func BenchmarkFiltering(b *testing.B) {
    // 大量のタスクを生成
    tasks := generateLargeTasks(10000)
    filter := NewProjectFilter("work")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        filtered := filter.Apply(tasks)
        _ = filtered
    }
}

// メモリ使用量のテスト
func TestMemoryUsage(t *testing.T) {
    var m1, m2 runtime.MemStats
    
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    // 大量のタスクを処理
    tasks := generateLargeTasks(100000)
    model := NewModel("test.txt", DefaultConfig())
    model.SetTasks(tasks)
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    memUsed := m2.Alloc - m1.Alloc
    t.Logf("Memory used: %d bytes", memUsed)
    
    // メモリ使用量が閾値以下であることを確認
    assert.Less(t, memUsed, uint64(100*1024*1024)) // 100MB
}
```

#### プロファイリング
```go
// メインファイルにプロファイリング追加
func main() {
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
        
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    
    // 通常のアプリケーション実行
    // ...
}
```

```bash
# プロファイリング実行
go run cmd/todotui/main.go -cpuprofile=cpu.prof sample.todo.txt

# プロファイル解析
go tool pprof cpu.prof
```

## 6. リリースプロセス

### 6.1 バージョン管理

#### セマンティックバージョニング
- **Major (X.0.0)**: 破壊的変更
- **Minor (0.X.0)**: 新機能追加（後方互換）
- **Patch (0.0.X)**: バグフィックス

#### リリースタグ
```bash
# タグの作成
git tag -a v1.2.0 -m "Release version 1.2.0"
git push origin v1.2.0

# リリースノートの生成
git log --oneline v1.1.0..v1.2.0 > CHANGELOG.md
```

### 6.2 CI/CDパイプライン

#### GitHub Actions設定
```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, darwin, windows]
        goarch: [amd64, arm64]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.24
    
    - name: Build
      env:
        GOOS: ${{ matrix.goos }}
        GOARCH: ${{ matrix.goarch }}
      run: |
        go build -ldflags "-X main.version=${{ github.ref_name }}" \
          -o todotui-${{ matrix.goos }}-${{ matrix.goarch }} \
          cmd/todotui/main.go
    
    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: todotui-${{ matrix.goos }}-${{ matrix.goarch }}
        path: todotui-${{ matrix.goos }}-${{ matrix.goarch }}
```

### 6.3 コントリビューション

#### プルリクエストのガイドライン
1. **Issue作成**: バグ報告や機能要求
2. **フォーク**: リポジトリをフォーク
3. **ブランチ作成**: `feature/new-feature` or `fix/bug-name`
4. **コミット**: 説明的なコミットメッセージ
5. **テスト**: 既存テストがパスすることを確認
6. **プルリクエスト**: 詳細な説明とともに作成

#### コミットメッセージ規約
```
<type>(<scope>): <subject>

<body>

<footer>
```

例:
```
feat(ui): add new filter for overdue tasks

Implement filtering functionality to show tasks that are past their due date.
The filter appears in the filter pane and updates the task count automatically.

Closes #123
``` 