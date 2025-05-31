#!/bin/bash

# Git pre-commit フックセットアップスクリプト
set -e

# カラー出力用
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}Setting up pre-commit hooks...${NC}"

# プロジェクトルートディレクトリを取得
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# pre-commitフックファイルが存在するかチェック
if [ -f "$PROJECT_ROOT/.git/hooks/pre-commit" ]; then
    echo -e "${YELLOW}Pre-commit hook already exists. Updating...${NC}"
fi

# pre-commitフックファイルを作成/更新
cp "$PROJECT_ROOT/scripts/pre-commit.sh" "$PROJECT_ROOT/.git/hooks/pre-commit"
chmod +x "$PROJECT_ROOT/.git/hooks/pre-commit"

echo -e "${GREEN}✓ Pre-commit hook set up successfully${NC}"
echo -e "${YELLOW}Pre-commit hook will run: make fmt, make lint, make test${NC}" 