#!/bin/bash

# todotui プロジェクト用 pre-commit フック
# コミット前に make lint, make fmt, make test を実行する

set -e

# プロジェクトルートに移動して共通ライブラリを読み込み
cd "$(git rev-parse --show-toplevel)"
source "scripts/common.sh"

echo -e "${BLUE}=== Pre-commit hooks for todotui ===${NC}"

log_info "Running make fmt..."
if ! make fmt; then
    log_error "make fmt failed"
    exit 1
fi
log_success "make fmt completed"

log_info "Running make lint..."
if ! make lint; then
    log_error "make lint failed"
    log_warning "Please fix linting errors before committing"
    exit 1
fi
log_success "make lint completed"

log_info "Running make test..."
if ! make test; then
    log_error "make test failed"
    log_warning "Please fix failing tests before committing"
    exit 1
fi
log_success "make test completed"

# フォーマットの変更があった場合、ステージングエリアに追加
if ! git diff --quiet --exit-code; then
    log_warning "Code formatting changes detected. Adding to staging area..."
    git add .
fi

echo -e "${GREEN}=== All pre-commit checks passed! ===${NC}" 