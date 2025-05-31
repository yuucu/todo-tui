# todotui Makefile

# プロジェクト情報
BINARY_NAME=todotui
MAIN_PATH=./cmd/todotui
BUILD_DIR=bin

# Go コマンド
GO=go

# カラー出力用
GREEN=\033[0;32m
BLUE=\033[0;34m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build run test fmt lint clean git-hooks-status install setup-hooks

# デフォルトターゲット
all: fmt test build

# ヘルプ表示
help: ## 使用可能なコマンドを表示
	@echo "$(BLUE)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

# 開発環境セットアップ
install: ## 開発に必要なツールとpre-commitフックをインストール
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@./scripts/install-dev-tools.sh
	@./scripts/setup-hooks.sh
	@echo "$(GREEN)✓ Development environment setup completed$(NC)"

# ビルド
build: ## バイナリをビルド
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# 実行
run: ## アプリケーションを実行
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	$(GO) run $(MAIN_PATH) -c sample-config.yaml

# テスト
test: ## テストを実行
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test -v ./...
	@echo "$(GREEN)✓ Tests completed$(NC)"

# コードフォーマット
fmt: ## コードをフォーマット
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	$(GO) mod tidy
	@echo "$(GREEN)✓ Formatting completed$(NC)"

# リンター
lint: ## リンターを実行
	@echo "$(BLUE)Running linter...$(NC)"
	golangci-lint run --config=.golangci.yml --timeout=5m
	@echo "$(GREEN)✓ Linting completed$(NC)"

# クリーンアップ
clean: ## ビルド成果物をクリーンアップ
	@echo "$(BLUE)Cleaning up...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -rf dist/
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

# Git フックの状態確認
git-hooks-status: ## Git フックの状態を確認
	@echo "$(BLUE)Git hooks status:$(NC)"
	@if [ -f .git/hooks/pre-commit ] && [ -x .git/hooks/pre-commit ]; then \
		echo "$(GREEN)✓ pre-commit hook is active$(NC)"; \
	elif [ -f .git/hooks/pre-commit.disabled ]; then \
		echo "$(YELLOW)○ pre-commit hook is disabled$(NC)"; \
	else \
		echo "$(RED)✗ pre-commit hook is not set up$(NC)"; \
	fi 
