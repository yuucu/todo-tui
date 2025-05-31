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

### 🔧 Go Install

```bash
go install github.com/yuucu/todotui/cmd/todotui@latest
```

## 🚀 Quick Start

```bash
# Check version
todotui --version

# Run with your todo file
todotui ~/todo.txt

# Run with custom config
todotui --config config.yaml ~/todo.txt

# Or use the sample file
todotui sample.todo.txt
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
| `q` | Quit |

## 📝 Task Format

Works with standard todo.txt format:
```
(A) Call Mom @phone +family due:2025-01-15
Buy milk @store +groceries
x 2025-01-14 Clean garage @home +chores
```

## 🎨 Configuration

Todo TUI uses YAML format for configuration. The recommended approach is to create a configuration file at `~/.config/todotui/config.yaml`.

### Quick Setup

```bash
# Create sample configuration file (recommended)
todotui --init-config

# This creates ~/.config/todotui/config.yaml with default settings
```

### Configuration File Location

You can specify a configuration file using the `--config` flag:

```bash
# Use specific config file
todotui --config /path/to/config.yaml
```

To create a configuration file, manually create a YAML file with your preferred settings based on the sample below.

### Sample Configuration

```yaml
# Todo TUI Configuration File
# ~/.config/todotui/config.yaml

# Theme settings (catppuccin, nord, default)
theme: catppuccin

# Priority levels (empty string for no priority)
priority_levels:
  - ""    # 優先度なし
  - A     # 最高優先度
  - B     # 高優先度
  - C     # 中優先度
  - D     # 低優先度

# Default todo.txt file path
default_todo_file: ~/todo.txt

# UI settings
ui:
  left_pane_ratio: 0.33        # Left pane width ratio (0.1-0.9)
  min_left_pane_width: 18      # Minimum left pane width
  min_right_pane_width: 28     # Minimum right pane width
  vertical_padding: 2          # Vertical padding (1-5)
```

### Supported Formats

While YAML is recommended, Todo TUI also supports:
- **YAML**: `.yaml`, `.yml` (recommended)
- **JSON**: `.json`
- **TOML**: `.toml`

## 🏗️ Build from Source

```bash
git clone https://github.com/yuucu/todotui.git
cd todotui
make build
```

## 📋 Requirements

- Go 1.24+
- Terminal with color support

---

*Simple. Fast. Distraction-free task management.* 

## Development

```bash
git clone https://github.com/yuucu/todotui.git
cd todotui
