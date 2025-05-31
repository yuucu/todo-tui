# Todo TUI - 設定仕様書

## 1. 設定システム概要

### 1.1 設定ファイルの場所

#### 優先順位（高い順）
1. **コマンドライン指定**: `--config /path/to/config.yaml`
2. **環境変数**: `TODO_TUI_CONFIG`で指定されたパス
3. **プロジェクト設定**: `./todotui.yaml` (カレントディレクトリ)
4. **ユーザー設定**: `~/.config/todotui/config.yaml`
5. **システム設定**: `/etc/todotui/config.yaml`
6. **デフォルト**: 組み込み設定値

#### 設定ファイル形式
サポートされる形式：
- **YAML**: `.yaml`, `.yml` (推奨)
- **JSON**: `.json`
- **TOML**: `.toml`
- **HCL**: `.hcl`

## 2. 設定項目一覧

### 2.1 完全な設定ファイル例

```yaml
# Todo TUI Configuration File
# ~/.config/todotui/config.yaml

# =====================================
# テーマ設定
# =====================================
theme: catppuccin  # catppuccin, nord, default

# =====================================
# 優先度レベル設定
# =====================================
priority_levels:
  - ""    # 優先度なし
  - A     # 最高優先度
  - B     # 高優先度  
  - C     # 中優先度
  - D     # 低優先度

# =====================================
# ファイル設定
# =====================================
default_todo_file: ~/todo.txt
backup:
  enabled: true
  directory: ~/.config/todotui/backups
  keep_days: 30

# =====================================
# UI設定
# =====================================
ui:
  # レイアウト設定
  left_pane_ratio: 0.33        # 左ペインの幅比率 (0.0-1.0)
  min_left_pane_width: 18      # 左ペイン最小幅
  min_right_pane_width: 28     # 右ペイン最小幅
  vertical_padding: 2          # 垂直パディング
  
  # 表示設定
  show_line_numbers: false     # 行番号表示
  show_task_count: true        # タスク数表示
  compact_mode: false          # コンパクト表示モード
  
  # 日付表示設定
  date_format: "2006-01-02"    # Go時間フォーマット
  relative_dates: true         # 相対日付表示（今日、明日等）
  
  # フィルター設定
  auto_filter_count: true      # フィルター数の自動更新
  show_empty_filters: false    # 空のフィルターを隠す

# =====================================
# キーバインディング設定
# =====================================
keybindings:
  # グローバル操作
  quit: "q"
  switch_pane: "tab"
  move_left: "h"
  move_right: "l"
  help: "?"
  
  # リスト操作
  up: ["j", "up"]
  down: ["k", "down"]
  select: "enter"
  home: "home"
  end: "end"
  
  # タスク操作
  add: "a"
  edit: "e"
  delete: "d"
  toggle_priority: "p"
  restore: "r"
  toggle_complete: "space"

# =====================================
# フィルター設定
# =====================================
filters:
  # デフォルトフィルター
  default: "all"
  
  # カスタムフィルター
  custom:
    urgent:
      name: "緊急"
      description: "優先度AまたはBで期限が今日または過去"
      expression: "(priority:A OR priority:B) AND (due:today OR overdue)"
    
    work:
      name: "仕事"
      description: "仕事関連のタスク"
      expression: "+work OR @office"

# =====================================
# 入力・言語設定
# =====================================
input:
  # IME設定
  ime_enabled: true
  fallback_locale: "ja_JP.UTF-8"
  
  # 入力補完
  auto_complete: true
  complete_projects: true      # プロジェクト名の補完
  complete_contexts: true      # コンテキスト名の補完

# =====================================
# 機能設定
# =====================================
features:
  # ファイル監視
  file_watcher: true
  auto_save: true
  save_delay: 500              # 自動保存遅延（ミリ秒）
  
  # 確認ダイアログ
  confirm_delete: true
  confirm_quit_unsaved: true
  
  # ソート設定
  default_sort: "priority"     # priority, date, alphabetical
  group_completed: true        # 完了タスクをグループ化

# =====================================
# 通知設定
# =====================================
notifications:
  enabled: false
  sound: false
  desktop: false
  
  # 期限通知
  due_reminders:
    - days: 1    # 1日前
    - hours: 2   # 2時間前

# =====================================
# デバッグ・ログ設定
# =====================================
debug:
  enabled: false
  log_file: ~/.config/todotui/debug.log
  log_level: info              # debug, info, warn, error
```

### 2.2 設定項目詳細

#### テーマ設定
```yaml
theme: catppuccin
# 利用可能な値: catppuccin, nord, default
# カスタムテーマファイルも指定可能: "/path/to/custom.yaml"
```

#### 優先度レベル
```yaml
priority_levels:
  - ""     # 必須: 優先度なしを表す空文字
  - A      # カスタマイズ可能
  - B
  - C
  - D
  - LOW    # 独自のレベル名も可能
  - HIGH
```

#### UI詳細設定
```yaml
ui:
  # レイアウト
  left_pane_ratio: 0.33       # 0.1-0.9の範囲
  min_left_pane_width: 18     # 最小15
  min_right_pane_width: 28    # 最小20
  vertical_padding: 2         # 1-5の範囲
  
  # 表示オプション
  show_line_numbers: false
  show_task_count: true
  compact_mode: false         # 行間を詰める
  word_wrap: true            # 長いタスクの折り返し
  
  # カラー設定
  use_colors: true
  force_colors: false        # カラー対応していない端末でも強制
  color_depth: auto          # auto, 256, truecolor
```

## 3. 環境変数による設定

### 3.1 サポートする環境変数

```bash
# 設定ファイル
export TODO_TUI_CONFIG="/path/to/config.yaml"

# テーマ設定
export TODO_TUI_THEME="nord"

# 優先度レベル
export TODO_TUI_PRIORITY_LEVELS="A,B,C,D"

# デフォルトファイル
export TODO_TUI_DEFAULT_FILE="~/Documents/todo.txt"

# デバッグモード
export TODO_TUI_DEBUG="true"
export TODO_TUI_LOG_LEVEL="debug"

# 言語・ロケール
export LANG="ja_JP.UTF-8"
export LC_CTYPE="ja_JP.UTF-8"
```

### 3.2 環境変数の優先順位

1. **コマンドラインオプション** (最高優先度)
2. **環境変数**
3. **設定ファイル**
4. **デフォルト値** (最低優先度)

## 4. カスタムテーマ作成

### 4.1 テーマファイル構造

```yaml
# custom-theme.yaml
name: "My Custom Theme"
description: "カスタムテーマの例"

colors:
  # 基本色
  background: "#1a1a1a"
  foreground: "#ffffff"
  primary: "#0078d4"
  secondary: "#6264a7"
  accent: "#00bcf2"
  
  # 状態色
  success: "#16a085"
  warning: "#f39c12"
  error: "#e74c3c"
  info: "#3498db"
  
  # 優先度色
  priority_a: "#e74c3c"    # 赤
  priority_b: "#f39c12"    # オレンジ
  priority_c: "#f1c40f"    # 黄
  priority_d: "#9b59b6"    # 紫
  
  # 特殊状態
  completed: "#7f8c8d"     # グレー
  overdue: "#c0392b"       # 暗赤
  
  # UI要素
  border: "#34495e"
  selected: "#2c3e50"
  focused: "#3498db"
  dimmed: "#95a5a6"

# ボーダースタイル
borders:
  style: rounded           # rounded, normal, thick, double
  filter_pane: true
  task_pane: true
  modal: true

# フォントスタイル
fonts:
  bold_priorities: true
  italic_projects: true
  underline_contexts: false
```

### 4.2 テーマの適用

```bash
# コマンドラインで指定
todotui --theme /path/to/custom-theme.yaml

# 設定ファイルで指定
echo 'theme: /path/to/custom-theme.yaml' >> ~/.config/todotui/config.yaml

# 環境変数で指定
export TODO_TUI_THEME="/path/to/custom-theme.yaml"
```

## 5. 設定ファイルの管理

### 5.1 設定ファイルの生成

```bash
# サンプル設定ファイルの生成
todotui --init-config

# 特定のパスに生成
todotui --init-config ~/.config/todotui/config.yaml

# 現在の設定をファイルに出力
todotui --export-config > current-config.yaml
```

### 5.2 設定の検証

```bash
# 設定ファイルの検証
todotui --validate-config ~/.config/todotui/config.yaml

# 現在の設定を表示
todotui --show-config

# 設定のテスト
todotui --test-config --config test-config.yaml
```

### 5.3 設定の移行

```bash
# 古い設定ファイルからの移行
todotui --migrate-config old-config.json new-config.yaml

# 設定のバックアップ
todotui --backup-config ~/.config/todotui/backup/
```

## 6. 高度な設定

### 6.1 プロファイル機能

```yaml
# ~/.config/todotui/profiles.yaml
profiles:
  work:
    theme: nord
    default_todo_file: ~/work/todo.txt
    filters:
      default: work
    ui:
      compact_mode: true
  
  personal:
    theme: catppuccin
    default_todo_file: ~/personal/todo.txt
    notifications:
      enabled: true
```

```bash
# プロファイルの使用
todotui --profile work
todotui --profile personal
```

### 6.2 条件付き設定

```yaml
# 環境別設定
conditions:
  terminal_width:
    - condition: "< 100"
      config:
        ui:
          compact_mode: true
          left_pane_ratio: 0.25
    - condition: "> 120"
      config:
        ui:
          show_line_numbers: true

  os:
    - condition: "darwin"
      config:
        keybindings:
          quit: "cmd+q"
    - condition: "linux"
      config:
        input:
          ime_enabled: true
```

### 6.3 プラグイン設定（将来拡張）

```yaml
plugins:
  enabled: true
  directory: ~/.config/todotui/plugins
  
  installed:
    - name: "sync-plugin"
      enabled: true
      config:
        provider: "dropbox"
        sync_interval: 300
    
    - name: "notification-plugin"
      enabled: false
      config:
        slack_webhook: "https://hooks.slack.com/..."
``` 