# 📋 Todo TUI

A clean, efficient terminal-based todo manager that respects the [todo.txt](http://todotxt.org/) format.

![todotui](https://github.com/user-attachments/assets/8e7223f2-0429-4733-a128-53ef2935a6aa)

## ✨ Features

- **Clean TUI Interface** - Navigate with intuitive keyboard shortcuts
- **Todo.txt Compatible** - Works with your existing todo.txt files
- **Smart Filtering** - Filter by projects (`+project`), contexts (`@context`), due dates, and more
- **Task Copy** - Copy task text to clipboard with visual feedback
- **Japanese Input Support** - Full IME support for international users

## 🚀 Installation

```bash
go install github.com/yuucu/todotui/cmd/todotui@latest
```

## 🚀 Quick Start

```bash
# Basic usage
todotui ~/todo.txt

# With custom configuration
todotui --config config.yaml ~/todo.txt

# Check available themes
todotui --list-themes
```

## ⌨️ Key Bindings

| Key | Action |
|-----|--------|
| `j/k` | Navigate lists |
| `Tab` | Switch between filter and task panes |
| `Enter` | Apply filter / Complete task |
| `a` | Add new task |
| `e` | Edit task |
| `d` | Delete task |
| `p` | Cycle priority (A→B→C→D→none) |
| `r` | Restore deleted/completed task |
| `y` | Copy task text to clipboard |
| `?` | Show help |
| `q` | Quit |

## 📝 Task Format

Supports standard todo.txt format:
```
(A) Call Mom @phone +family due:2025-01-15
Buy milk @store +groceries
x 2025-01-14 Clean garage @home +chores
```

## ⚙️ Configuration

Create a YAML configuration file and specify it with `--config`:

```yaml
# config.yaml
theme: catppuccin                # Available: catppuccin, nord, default
priority_levels: ["", A, B, C, D]
default_todo_file: ~/todo.txt
ui:
  left_pane_ratio: 0.33
  min_left_pane_width: 18
  min_right_pane_width: 28
  vertical_padding: 2
```

Supported formats: YAML (recommended), JSON, TOML

## 🏗️ Development

```bash
git clone https://github.com/yuucu/todotui.git
cd todotui
make build
```

**Requirements:** Go 1.24+, color-capable terminal

### 🔄 Release Management

このプロジェクトはSemantic Versioningを採用し、Release Drafterで自動的にリリースノートを生成します。

#### コミットメッセージの規約

適切なリリースノート生成とバージョン決定のため、以下の形式でコミットしてください：

- `feat:` または `feature:` - 新機能 (minor version up)
- `fix:` または `bug:` - バグ修正 (patch version up)
- `BREAKING CHANGE:` または `!:` - 破壊的変更 (major version up)
- `docs:` - ドキュメント更新
- `chore:` - メンテナンス作業
- `deps:` - 依存関係の更新

#### プルリクエストラベル

PRには以下のラベルを適切に付けてください：

- **Semantic Versioning**: `breaking`, `feature`, `bug`, `patch`
- **Categories**: `docs`, `chore`, `dependencies`

#### リリースプロセス

1. **開発**: featureブランチで開発
2. **PR作成**: mainブランチへのPRを作成
3. **自動ラベリング**: Release Drafterが自動的にラベルを付与
4. **ドラフト更新**: PRマージ時にリリースドラフトが自動更新
5. **リリース**: GitHubでリリースドラフトを公開
6. **自動配布**: GoReleaserが自動的にバイナリをビルド・配布

---

*Simple. Fast. Distraction-free task management.*
