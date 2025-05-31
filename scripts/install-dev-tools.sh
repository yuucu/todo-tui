#!/bin/bash

# todotui 開発ツールインストールスクリプト
set -e

# カラー出力用
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}Installing development tools...${NC}"

# golangci-lint のインストール
echo -e "${BLUE}Checking golangci-lint...${NC}"
if command -v golangci-lint >/dev/null 2>&1; then
    echo -e "${GREEN}✓ golangci-lint is already installed${NC}"
else
    echo -e "${YELLOW}Installing golangci-lint...${NC}"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
    echo -e "${GREEN}✓ golangci-lint installed${NC}"
fi

# その他の開発ツールがある場合はここに追加
# 例: gofumpt, govulncheck など

echo -e "${GREEN}✓ All development tools installed${NC}" 