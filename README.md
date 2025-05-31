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
```

