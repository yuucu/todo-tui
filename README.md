# 📋 Todo TUI

A clean, efficient terminal-based todo manager that respects the [todo.txt](http://todotxt.org/) format.

![todo-tui](https://github.com/user-attachments/assets/8e7223f2-0429-4733-a128-53ef2935a6aa)

## ✨ Features

- **Clean TUI Interface** - Navigate with intuitive keyboard shortcuts
- **Todo.txt Compatible** - Works with your existing todo.txt files
- **Smart Filtering** - Filter by projects (`+project`), contexts (`@context`), due dates, and more
- **Japanese Input Support** - Full IME support for international users

## 🚀 Quick Start

```bash
# Install
go install github.com/yuucu/todo-tui/cmd/todo-tui@latest

# Run with your todo file
todo-tui ~/todo.txt

# Run with custom config
todo-tui --config config.yaml ~/todo.txt

# Or use the sample file
todo-tui sample.todo.txt
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
| `q` | Quit |

## 📝 Task Format

Works with standard todo.txt format:
```
(A) Call Mom @phone +family due:2025-01-15
Buy milk @store +groceries
x 2025-01-14 Clean garage @home +chores
```

## 🎨 Configuration

Supports multiple configuration formats. Choose your preferred format:

### YAML Configuration (`config.yaml`)
```yaml
theme: catppuccin
priority_levels:
  - ""
  - A
  - B
  - C
  - D
default_todo_file: ~/todo.txt
ui:
  left_pane_ratio: 0.33
  min_left_pane_width: 18
  min_right_pane_width: 28
  vertical_padding: 2
```

Place your config file in `~/.config/todo-tui/` or specify with `--config` flag.

## 🏗️ Build from Source

```bash
git clone https://github.com/yuucu/todo-tui.git
cd todo-tui
make build
```

## 📋 Requirements

- Go 1.24+
- Terminal with color support

---

*Simple. Fast. Distraction-free task management.* 
