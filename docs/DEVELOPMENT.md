# Todo TUI - 開発ガイド

## 開発環境セットアップ

```bash
# クローン
git clone https://github.com/yuucu/todotui.git
cd todotui

# 依存関係取得
go mod tidy

# ビルド
make build

# 実行
./bin/todotui sample.todo.txt
```

## 開発コマンド

```bash
# ビルド
make build

# テスト実行
make test

# リント
make lint

# クリーン
make clean

# リリースビルド
make release
```

## ファイル構成

- `cmd/todotui/main.go` - エントリーポイント、CLI引数処理
- `internal/ui/model.go` - メインアプリケーションロジック
- `internal/ui/view.go` - 画面描画とレイアウト
- `internal/ui/filters.go` - タスクフィルタリング機能
- `internal/todo/storage.go` - ファイル入出力

## デバッグ

```bash
# デバッグモードでビルド
go build -tags debug ./cmd/todotui

# ログ出力付きで実行
DEBUG=1 ./todotui sample.todo.txt
```

## テスト

```bash
# 全テスト実行
go test ./...

# カバレッジ付き
go test -cover ./...

# 特定パッケージのテスト
go test ./internal/ui
```

## コントリビューション

1. フォーク
2. フィーチャーブランチ作成 (`git checkout -b feature/awesome-feature`)
3. コミット (`git commit -m 'Add awesome feature'`)
4. プッシュ (`git push origin feature/awesome-feature`)
5. Pull Request作成

## コーディング規約

- Go標準のコーディングスタイルに従う
- `gofmt`と`go vet`を実行
- コメントは日本語でも英語でもOK
- テストカバレッジ80%以上を目指す 