#!/bin/bash

# todotui 共通スクリプトライブラリ
# 他のスクリプトから source されることを想定

# カラー出力用
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 共通関数
log_info() {
    echo -e "${BLUE}$1${NC}"
}

log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}$1${NC}"
}

log_error() {
    echo -e "${RED}✗ $1${NC}"
}

# プロジェクトルートディレクトリを取得
get_project_root() {
    cd "$(dirname "${BASH_SOURCE[1]}")/.." && pwd
} 