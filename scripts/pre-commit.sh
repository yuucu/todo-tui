#!/bin/bash

# todotui プロジェクト用 pre-commit フック
# コミット前に make lint, make fmt, make test を実行する

set -e

# カラー出力用
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Pre-commit hooks for todotui ===${NC}"

# プロジェクトルートに移動
cd "$(git rev-parse --show-toplevel)"

echo -e "${BLUE}Running make fmt...${NC}"
if ! make fmt; then
    echo -e "${RED}✗ make fmt failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ make fmt completed${NC}"

echo -e "${BLUE}Running make lint...${NC}"
if ! make lint; then
    echo -e "${RED}✗ make lint failed${NC}"
    echo -e "${YELLOW}Please fix linting errors before committing${NC}"
    exit 1
fi
echo -e "${GREEN}✓ make lint completed${NC}"

echo -e "${BLUE}Running make test...${NC}"
if ! make test; then
    echo -e "${RED}✗ make test failed${NC}"
    echo -e "${YELLOW}Please fix failing tests before committing${NC}"
    exit 1
fi
echo -e "${GREEN}✓ make test completed${NC}"

# フォーマットの変更があった場合、ステージングエリアに追加
if ! git diff --quiet --exit-code; then
    echo -e "${YELLOW}Code formatting changes detected. Adding to staging area...${NC}"
    git add .
fi

echo -e "${GREEN}=== All pre-commit checks passed! ===${NC}" 