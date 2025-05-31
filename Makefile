# Go プロジェクト用 Makefile (シンプル版)
# プロジェクト情報
BINARY_NAME=todotui
MAIN_PATH=./cmd/todotui
BUILD_DIR=bin
REQUIRED_GO_VERSION=1.24

# Go コマンド
GO=go

# カラー出力用
GREEN=\033[0;32m
BLUE=\033[0;34m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build run test fmt lint clean deps tidy all lint-install check-go-version release-snapshot release-test

# デフォルトターゲット
all: check-go-version fmt test build

# Goバージョンチェック
check-go-version: ## 必要なGoバージョンをチェック
	@echo "$(BLUE)Checking Go version...$(NC)"
	@GO_VERSION=$$($(GO) version | sed 's/go version go\([0-9]*\.[0-9]*\).*/\1/'); \
	REQUIRED_VERSION=$(REQUIRED_GO_VERSION); \
	if [ "$$(echo "$$GO_VERSION $$REQUIRED_VERSION" | awk '{print ($$1 >= $$2)}')" = "1" ]; then \
		echo "$(GREEN)✓ Go version $$GO_VERSION is compatible (required: $$REQUIRED_VERSION+)$(NC)"; \
	else \
		echo "$(RED)✗ Go version $$GO_VERSION is not compatible. Required: $$REQUIRED_VERSION+$(NC)"; \
		echo "$(YELLOW)Please update Go to version $$REQUIRED_VERSION or higher$(NC)"; \
		exit 1; \
	fi

# ヘルプ表示
help: ## 使用可能なコマンドを表示
	@echo "$(BLUE)よく使うコマンド:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

# ビルド
build: check-go-version ## バイナリをビルド
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

lint: lint-fast ## 高速リンター（CI用）を実行

lint-fast: ## 高速リンター（CI用）を実行
	@echo "$(BLUE)Running fast linter (CI config)...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config=.golangci.yml --timeout=5m --fast; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		make lint-install; \
		golangci-lint run --config=.golangci.yml --timeout=5m --fast; \
	fi

lint-full: ## 詳細リンター（開発用）を実行
	@echo "$(BLUE)Running full linter (development config)...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config=.golangci-dev.yml --timeout=10m; \
	else \
		echo "$(YELLOW)golangci-lint not found. Installing...$(NC)"; \
		make lint-install; \
		golangci-lint run --config=.golangci-dev.yml --timeout=10m; \
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

# Goreleaser関連タスク
goreleaser-install: ## GoReleaserをインストール
	@echo "$(BLUE)Installing GoReleaser...$(NC)"
	@if command -v goreleaser >/dev/null 2>&1; then \
		echo "$(GREEN)GoReleaser is already installed$(NC)"; \
	else \
		if command -v brew >/dev/null 2>&1; then \
			brew install goreleaser; \
		else \
			go install github.com/goreleaser/goreleaser@latest; \
		fi; \
	fi

release-snapshot: goreleaser-install ## GoReleaserでスナップショットリリースをテスト
	@echo "$(BLUE)Building snapshot release with GoReleaser...$(NC)"
	goreleaser release --snapshot --clean
	@echo "$(GREEN)Snapshot release completed. Check dist/ directory$(NC)"

release-test: goreleaser-install ## GoReleaser設定をテスト（リリースなし）
	@echo "$(BLUE)Testing GoReleaser configuration...$(NC)"
	goreleaser check
	@echo "$(GREEN)GoReleaser configuration is valid$(NC)"

release-clean: ## GoReleaserの成果物をクリーンアップ
	@echo "$(BLUE)Cleaning up GoReleaser artifacts...$(NC)"
	rm -rf dist/
	@echo "$(GREEN)Cleanup completed$(NC)"

# クリーンアップ（GoReleaser対応）
clean: release-clean ## ビルド成果物をクリーンアップ
	@echo "$(BLUE)Cleaning up build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Cleanup completed$(NC)" 
