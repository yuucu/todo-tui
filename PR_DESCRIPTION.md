# feat: implement slog-based logging system

## 概要

slogベースの構造化ログシステムを実装しました。CLIアプリケーションに適した人に見やすいログ形式で、デバッグログの切り替えとファイル出力機能を提供します。**ログファイル出力は常に有効**で、問題の調査やトレーサビリティを確保します。

## 主な機能

### 🆕 新機能
- **独立したloggerパッケージ** (`pkg/logger`) - 再利用可能な設計
- **CLIオプション**: `--debug/-d` (デバッグログ), `--log-file/-l` (カスタムログファイルパス)
- **クロスプラットフォーム対応**: macOS/Linux/Windowsの適切なログディレクトリ
- **設定ファイル対応**: ログレベル、カスタムパス、保持期間を設定可能
- **自動ローテーション**: 日付別ファイル + 古いログの自動削除
- **必須ファイル出力**: 常にログファイルに記録（無効化不可）

### 🔧 技術的改善
- **構造化ログ**: key-value形式の詳細な情報記録
- **人に見やすい形式**: 時刻は `HH:MM:SS` 形式
- **デバッグ時ソース情報**: ファイル名:行番号の表示
- **二重出力戦略**: ログ記録 + ユーザー向け直接メッセージ
- **運用重視**: 問題調査のため常にファイル出力

## ログファイル保存場所

- **macOS**: `~/Library/Logs/todotui/`
- **Linux**: `~/.local/share/todotui/logs/`  
- **Windows**: `%APPDATA%/todotui/logs/`

## 使用例

```bash
# デバッグログを有効にする（ファイル出力は自動で有効）
./todotui --debug ~/todo.txt

# カスタムログファイルパスを指定
./todotui --log-file /path/to/custom.log ~/todo.txt

# 通常使用（デフォルトのログディレクトリに自動出力）
./todotui ~/todo.txt
```

## 設定ファイル例

`sample-config.yaml` でログ設定の例を確認できます：

```yaml
logging:
  enable_debug: false       # デバッグログの有効/無効
  log_file_path: ""         # カスタムパス（空=デフォルト）
  max_log_days: 30         # ログ保持日数
```

## 出力例

### 通常ログ (INFO レベル)
```
01:18:50 level=INFO msg="todotui started" version=dev commit=none
01:18:50 level=INFO msg="Loaded tasks from file" file=sample.todo.txt task_count=1023
01:18:50 level=INFO msg="Starting Bubble Tea program"
01:18:53 level=INFO msg="todotui exited"
```

### デバッグログ (DEBUG レベル)
```
01:18:50 level=DEBUG source=model.go:45 msg="Creating new model" todo_file=sample.todo.txt theme=catppuccin
01:18:50 level=DEBUG source=model.go:84 msg="Initializing file watcher"
01:18:50 level=DEBUG source=model.go:95 msg="File watcher initialized successfully" file=sample.todo.txt
```

## 変更されたファイル

- `pkg/logger/logger.go` - 新規: 独立したログパッケージ
- `cmd/todotui/main.go` - CLIフラグとログ初期化を追加
- `internal/ui/config.go` - ログ設定の追加（ファイル出力必須化）
- `internal/ui/model.go` - ログ出力の統合
- `sample-config.yaml` - ログ設定セクションを追加

## 動作確認済み

✅ デバッグログの有効/無効切り替え  
✅ 常時ファイル出力機能  
✅ OSごとの適切なログディレクトリ作成  
✅ 設定ファイルからの設定読み込み  
✅ ログファイルのローテーション  
✅ 既存機能の後方互換性

## テスト

```bash
# ビルド確認
go build -o todotui cmd/todotui/main.go

# 通常使用（自動でログファイル出力）
timeout 2s ./todotui sample.todo.txt

# デバッグログテスト
timeout 2s ./todotui --debug sample.todo.txt

# ログファイル確認
ls -la ~/Library/Logs/todotui/
cat ~/Library/Logs/todotui/todotui-$(date +%Y-%m-%d).log
```

## 運用上の利点

- **必須ログ出力**: 問題発生時の調査に必要な情報が確実に記録
- **自動クリーンアップ**: ディスク容量を考慮した自動ローテーション
- **標準準拠**: OSごとの適切なログディレクトリ使用 