# Go プロジェクト用 Makefile (シンプル版)
# プロジェクト情報
BINARY_NAME=todo-tui
MAIN_PATH=./cmd/todo-tui
BUILD_DIR=bin

# Go コマンド
GO=go

# カラー出力用
GREEN=\033[0;32m
BLUE=\033[0;34m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: help build run test fmt lint clean deps tidy all lint-install

# デフォルトターゲット
all: fmt test build

# ヘルプ表示
help: ## 使用可能なコマンドを表示
	@echo "$(BLUE)よく使うコマンド:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

# ビルド
build: ## バイナリをビルド
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# 実行
run: ## アプリケーションを実行
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	$(GO) run $(MAIN_PATH) sample.todo.txt -c sample-config.yaml

# テスト
test: ## テストを実行
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test -v ./...

# コード品質
fmt: ## コードをフォーマット
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	$(GO) mod tidy

# golangci-lint インストール
lint-install: ## golangci-lint をインストール
	@echo "$(BLUE)Installing golangci-lint...$(NC)"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: ## リンターを実行 (CIと同じ設定)
	@echo "$(BLUE)Running linter with .golangci.yml config...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config=.golangci.yml --timeout=10m; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		make lint-install; \
		golangci-lint run --config=.golangci.yml --timeout=10m; \
	fi

install: ## バイナリをシステムにインストール
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	$(GO) install $(MAIN_PATH)
	@echo "$(GREEN)Installation completed$(NC)"

# CI/CD用のターゲット
ci-lint: ## CI用リンター（GitHub Actionsと同じ）
	@echo "$(BLUE)Running CI linter...$(NC)"
	golangci-lint run --config=.golangci.yml --timeout=10m

ci-test: ## CI用テスト（カバレッジ付き）
	@echo "$(BLUE)Running CI tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out ./...

ci-build: ## CI用ビルド（複数プラットフォーム）
	@echo "$(BLUE)Building for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)Multi-platform build completed$(NC)" 
