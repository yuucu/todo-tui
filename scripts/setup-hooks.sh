#!/bin/bash

# Git pre-commit フックセットアップスクリプト
set -e

# 共通ライブラリを読み込み
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

log_info "Setting up pre-commit hooks..."

# プロジェクトルートディレクトリを取得
PROJECT_ROOT=$(get_project_root)

# pre-commitフックファイルが存在するかチェック
if [ -f "$PROJECT_ROOT/.git/hooks/pre-commit" ]; then
    log_warning "Pre-commit hook already exists. Updating..."
fi

# pre-commitフックファイルを作成/更新
cp "$PROJECT_ROOT/scripts/pre-commit.sh" "$PROJECT_ROOT/.git/hooks/pre-commit"
chmod +x "$PROJECT_ROOT/.git/hooks/pre-commit"

log_success "Pre-commit hook set up successfully"
log_warning "Pre-commit hook will run: make fmt, make lint, make test" 