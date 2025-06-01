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

Todo TUI can be configured using a `config.yaml` file for detailed customization.

### 📍 Configuration File Location

- **macOS/Linux**: `~/.config/todotui/config.yaml`
- **Windows**: `%APPDATA%/todotui/config.yaml`

### 🎨 Basic Configuration Example

```yaml
# ~/.config/todotui/config.yaml
theme: catppuccin                # Available: catppuccin, nord, everforest-dark, everforest-light
priority_levels: ["", A, B, C, D]
default_todo_file: ~/todo.txt
```

**📚 Complete Configuration Reference**: 
- **[sample-config.yaml](sample-config.yaml)** - Comprehensive configuration file with all options and detailed explanations

### 💻 Usage

```bash
# Automatic config detection (loads ~/.config/todotui/config.yaml)
todotui ~/my-todo.txt

# If default_todo_file is set in config, no CLI argument needed
todotui

# Custom config file
todotui --config /path/to/config.yaml
```

**Supported formats**: YAML, JSON, TOML

## 🏗️ Development

### Quick Setup

```bash
git clone https://github.com/yuucu/todotui.git
cd todotui
make install  # Install development tools and set up pre-commit hooks
make build    # Build the application
```

