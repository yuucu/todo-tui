# tuido 開発タスクリスト

## 0. プロジェクト初期設定
- [x] Go modules 初期化 (`go mod init github.com/yuucu/tuido`)
- [x] 依存ライブラリ導入
    - [x] todotxt (`github.com/1set/todotxt`)
    - [x] Bubble Tea (`github.com/charmbracelet/bubbletea`, `bubbles`, `lipgloss`)
    - [x] fsnotify (`github.com/fsnotify/fsnotify`)
- [x] ディレクトリ構成を作成 (`cmd/`, `internal/ui`, `internal/todo`, `internal/infra`)
- [x] Git リポジトリ初期化 & `.gitignore` 追加

## 1. ストレージ層
- [x] `internal/todo` パッケージ作成
    - [x] `Load(path string) (*todo.TaskList, error)` 実装
    - [x] `Save(list *todo.TaskList, path string) error` 実装
- [x] ユニットテスト作成 (`_test.go`)

## 2. 最小 UI (1 ペイン)
- [x] `list.Model` で未完了タスク表示
- [x] `j` / `k` でカーソル上下
- [x] `q` または `ctrl+c` で終了
- [x] `model` 構造体と `main.go` の雛形作成

## 3. CRUD 操作
- [ ] `a` 新規タスク入力 (textarea)
- [ ] `e` 選択タスク編集
- [ ] `enter or x` 完了トグル
- [ ] `d` 削除
- [ ] `p` 優先度サイクル (A→Z)
- [ ] 各操作後に `Save()` 実行

## 4. フィルタ／検索
- [ ] `/` インクリメンタル検索
- [ ] `+project` フィルタ
- [ ] `@context` フィルタ
- [ ] 期日ベースのクイックフィルタ (今日／期限切れ)

## 5. マルチペイン
- [ ] `lipgloss.JoinHorizontal/Vertical` でレイアウト作成
    - [ ] 左: 未完了
    - [ ] 右: 完了
    - [ ] 下: ヘルプ
- [ ] `Tab` でペイン切替
- [ ] 各ペイン専用の `list.Model` 管理

## 6. ファイルウォッチ
- [ ] `fsnotify` でタスクファイル変更検知
- [ ] 更新時に `Load()` 再実行

## 7. 拡張機能
- [ ] 期日が近いタスクをピン留め
- [ ] 通知・アラート連携 (notify-send など)
- [ ] Git 操作ショートカット (commit 等)

## 8. ドキュメント & CI
- [ ] README 更新 (インストール方法・キーバインド一覧)
- [ ] GitHub Actions で `go test` / `golangci-lint` 実行
- [ ] リリースビルド配布 (Homebrew tap / goreleaser)
