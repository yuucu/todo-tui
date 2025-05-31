# Todo TUI - アーキテクチャ概要

## プロジェクト構造

```
todotui/
├── cmd/todotui/          # エントリーポイント
│   └── main.go           # CLI実装
├── internal/
│   ├── ui/               # UIコンポーネント
│   │   ├── model.go      # アプリケーション状態
│   │   ├── view.go       # 画面描画
│   │   ├── filters.go    # フィルタリング機能
│   │   ├── config.go     # 設定管理
│   │   ├── colors.go     # テーマ管理
│   │   └── ime_helper.go # 日本語入力支援
│   └── todo/             # データ層
│       └── storage.go    # ファイル入出力
```

## アーキテクチャパターン

**Model-View-Update (MVU)**
- Model: アプリケーション状態 (`ui/model.go`)
- View: 画面表示 (`ui/view.go`)  
- Update: イベント処理とステート更新

## 主要ライブラリ

| ライブラリ | 用途 |
|------------|------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUIフレームワーク |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | スタイリング |
| [1set/todotxt](https://github.com/1set/todotxt) | todo.txt解析 |
| [Viper](https://github.com/spf13/viper) | 設定管理 |

## データフロー

```
todo.txt → storage.go → model.go → filters → view.go → ターミナル
                ↑                                        ↓
           キーボード入力 ←─────────────── Update ←─────┘
```

## 要件

- Go 1.24+
- カラー対応ターミナル
- Unix系OS (macOS/Linux) 