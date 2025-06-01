# 📋 Todo TUI

⌨️ A Vim-like TUI that honors the simplicity of todo.txt 📝


<table>
<tr>
<td width="50%">
<img width="100%" alt="todotui_top" src="https://github.com/user-attachments/assets/2dcd692f-bd63-442b-af37-91197f738feb" />
</td>
<td width="50%">
<img width="100%" alt="todotui_help" src="https://github.com/user-attachments/assets/4c1965ef-5dfa-4d11-8b34-4d714569d668" />
</td>
</tr>
</table>

---

## ✨ Features

- ⚡ **Vim-like TUI** — Navigate with intuitive, familiar keybindings (`j`, `k`, etc.)
- 📄 **todo.txt Compatible** — Fully supports the standard todo.txt format
- 🔍 **Powerful Filtering** — Filter by `+project`, `@context`, due dates, and more
- 📋 **Clipboard Support** — Easily copy task text with visual confirmation

---

## 🚀 Installation

```bash
go install github.com/yuucu/todotui/cmd/todotui@latest
```

## 🚀 Quick Start

```bash
# Basic usage
todotui ~/todo.txt

# If default_todo_file is set in config file, no argument needed
todotui

# With custom configuration
todotui --config config.yaml ~/todo.txt
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

Create a YAML configuration file at `~/.config/todotui/config.yaml`:

```yaml
# ~/.config/todotui/config.yaml
theme: catppuccin                # Available: catppuccin, nord, everforest-dark, everforest-light
priority_levels: ["", A, B, C, D]
default_todo_file: ~/todo.txt
ui:
  left_pane_ratio: 0.33
  min_left_pane_width: 18
  min_right_pane_width: 28
  vertical_padding: 2
```

**Usage:**
```bash
# Automatic config detection (loads ~/.config/todotui/config.yaml)
todotui ~/my-todo.txt

# If default_todo_file is set in config, no CLI argument needed
todotui

# Custom config file
todotui --config /path/to/config.yaml
```

Supported formats: YAML, JSON, TOML

## 🏗️ Development

### Quick Setup

```bash
git clone https://github.com/yuucu/todotui.git
cd todotui
make install  # Install development tools and set up pre-commit hooks
make build    # Build the application
make coverage # Run tests with coverage report
```

