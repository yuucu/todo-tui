#!/bin/bash

# todotui 開発ツールインストールスクリプト
set -e

# 共通ライブラリを読み込み
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

log_info "Installing development tools..."

# golangci-lint のインストール
log_info "Checking golangci-lint..."
if command -v golangci-lint >/dev/null 2>&1; then
    log_success "golangci-lint is already installed"
else
    log_warning "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
    log_success "golangci-lint installed"
fi

# その他の開発ツールがある場合はここに追加
# 例: gofumpt, govulncheck など

log_success "All development tools installed" 