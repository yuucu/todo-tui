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

.PHONY: help build run test fmt lint clean git-hooks-status

# デフォルトターゲット
all: fmt test build

# ヘルプ表示
help: ## 使用可能なコマンドを表示
	@echo "$(BLUE)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

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
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config=.golangci.yml --timeout=5m; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8; \
		golangci-lint run --config=.golangci.yml --timeout=5m; \
	fi
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
