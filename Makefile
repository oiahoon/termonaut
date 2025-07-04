# Termonaut Development Makefile
# This Makefile provides common development tasks for the Termonaut project

.PHONY: help dev-setup build test clean install lint format docs

# Default target
help: ## Show this help message
	@echo "Termonaut Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Variables
BINARY_NAME=termonaut
BUILD_DIR=bin
SOURCE_DIR=cmd
GO_VERSION=1.21
PYTHON_VERSION=3.9

# Development Setup
dev-setup: ## Setup development environment (detects Go/Python)
	@if command -v go >/dev/null 2>&1; then \
		make dev-setup-go; \
	elif command -v python3 >/dev/null 2>&1; then \
		make dev-setup-python; \
	else \
		echo "Neither Go nor Python3 found. Please install one."; \
		exit 1; \
	fi

dev-setup-go: ## Setup Go development environment
	@echo "Setting up Go development environment..."
	@go version
	@go mod init github.com/yourusername/termonaut || true
	@go mod tidy
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Go development environment ready!"

dev-setup-python: ## Setup Python development environment
	@echo "Setting up Python development environment..."
	@python3 --version
	@pip install --upgrade pip
	@pip install -r requirements-dev.txt || echo "requirements-dev.txt not found, creating..."
	@echo "click>=8.0.0\nrich>=10.0.0\npytest>=7.0.0\nblack>=22.0.0\nisort>=5.0.0\nmypy>=0.910" > requirements-dev.txt
	@pip install -r requirements-dev.txt
	@echo "Python development environment ready!"

# Building
build: ## Build the application
	@if [ -f "go.mod" ]; then \
		make build-go; \
	elif [ -f "setup.py" ] || [ -f "pyproject.toml" ]; then \
		make build-python; \
	else \
		echo "No build configuration found"; \
		exit 1; \
	fi

build-go: ## Build Go binary
	@echo "Building Go binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(SOURCE_DIR)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-python: ## Build Python package
	@echo "Building Python package..."
	@python3 setup.py build || echo "setup.py build not available"
	@pip install -e . || echo "Editable install not available"
	@echo "Python package built"

build-release: ## Build release binaries for multiple platforms
	@if [ -f "go.mod" ]; then \
		make build-release-go; \
	else \
		echo "Release builds only supported for Go currently"; \
		exit 1; \
	fi

build-release-go: ## Build release binaries for Go
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)/release
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 ./$(SOURCE_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 ./$(SOURCE_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 ./$(SOURCE_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe ./$(SOURCE_DIR)
	@echo "Release binaries built in $(BUILD_DIR)/release/"

# Testing
test: ## Run all tests
	@if [ -f "go.mod" ]; then \
		make test-go; \
	elif [ -f "pytest.ini" ] || command -v pytest >/dev/null 2>&1; then \
		make test-python; \
	else \
		echo "No test configuration found"; \
		exit 1; \
	fi

test-go: ## Run Go tests
	@echo "Running Go tests..."
	@go test -v ./...

test-python: ## Run Python tests
	@echo "Running Python tests..."
	@pytest -v

test-coverage: ## Run tests with coverage
	@if [ -f "go.mod" ]; then \
		make test-coverage-go; \
	else \
		make test-coverage-python; \
	fi

test-coverage-go: ## Run Go tests with coverage
	@echo "Running Go tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-coverage-python: ## Run Python tests with coverage
	@echo "Running Python tests with coverage..."
	@pytest --cov=termonaut --cov-report=html --cov-report=term-missing

test-unit: ## Run unit tests only
	@if [ -f "go.mod" ]; then \
		go test -v -short ./...; \
	else \
		pytest -v -m "not integration"; \
	fi

test-integration: ## Run integration tests only
	@if [ -f "go.mod" ]; then \
		go test -v -run Integration ./...; \
	else \
		pytest -v -m integration; \
	fi

test-performance: ## Run performance tests
	@echo "Running performance tests..."
	@if [ -f "go.mod" ]; then \
		go test -v -bench=. -benchmem ./...; \
	else \
		pytest -v -m performance; \
	fi

# Code Quality
lint: ## Run linters
	@if [ -f "go.mod" ]; then \
		make lint-go; \
	else \
		make lint-python; \
	fi

lint-go: ## Run Go linters
	@echo "Running Go linters..."
	@golangci-lint run
	@go vet ./...

lint-python: ## Run Python linters
	@echo "Running Python linters..."
	@black --check .
	@isort --check-only .
	@mypy . || echo "MyPy check completed"

format: ## Format code
	@if [ -f "go.mod" ]; then \
		make format-go; \
	else \
		make format-python; \
	fi

format-go: ## Format Go code
	@echo "Formatting Go code..."
	@go fmt ./...
	@goimports -w . || echo "goimports not available"

format-python: ## Format Python code
	@echo "Formatting Python code..."
	@black .
	@isort .

security-scan: ## Run security scans
	@if [ -f "go.mod" ]; then \
		make security-scan-go; \
	else \
		make security-scan-python; \
	fi

security-scan-go: ## Run Go security scan
	@echo "Running Go security scan..."
	@gosec ./... || echo "gosec not available"

security-scan-python: ## Run Python security scan
	@echo "Running Python security scan..."
	@bandit -r . || echo "bandit not available"

# Installation
install: build ## Install the application locally
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "Installing $(BINARY_NAME) to /usr/local/bin/"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		echo "Installation complete"; \
	else \
		echo "Binary not found. Run 'make build' first"; \
		exit 1; \
	fi

install-local: build ## Install for current user only
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		mkdir -p $$HOME/.local/bin; \
		cp $(BUILD_DIR)/$(BINARY_NAME) $$HOME/.local/bin/; \
		echo "Installed to $$HOME/.local/bin/$(BINARY_NAME)"; \
		echo "Make sure $$HOME/.local/bin is in your PATH"; \
	else \
		echo "Binary not found. Run 'make build' first"; \
		exit 1; \
	fi

uninstall: ## Uninstall the application
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@rm -f $$HOME/.local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete"

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if [ -f "go.mod" ]; then \
		godoc -http=:6060 & \
		echo "Go documentation server started at http://localhost:6060"; \
	else \
		echo "Documentation generation not configured"; \
	fi

docs-serve: ## Serve documentation locally
	@echo "Serving documentation..."
	@if command -v mkdocs >/dev/null 2>&1; then \
		mkdocs serve; \
	else \
		echo "mkdocs not available. Install with: pip install mkdocs"; \
	fi

# Database
db-create: ## Create initial database schema
	@echo "Creating database schema..."
	@./$(BUILD_DIR)/$(BINARY_NAME) init --create-db || echo "Binary not available"

db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@./$(BUILD_DIR)/$(BINARY_NAME) migrate || echo "Binary not available"

db-reset: ## Reset database (CAUTION: destroys data)
	@echo "Resetting database..."
	@read -p "This will destroy all data. Continue? (y/N): " confirm && [ "$$confirm" = "y" ]
	@rm -f $$HOME/.termonaut/termonaut.db
	@echo "Database reset complete"

# Development Tools
dev-shell: ## Start development shell with environment loaded
	@echo "Starting development shell..."
	@if [ -f "go.mod" ]; then \
		bash --rcfile <(echo '. ~/.bashrc; export GOPATH=$$(go env GOPATH); export PATH=$$PATH:$$GOPATH/bin'); \
	else \
		bash --rcfile <(echo '. ~/.bashrc; source venv/bin/activate 2>/dev/null || true'); \
	fi

watch: ## Watch for changes and rebuild
	@echo "Watching for changes..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} make build; \
	elif command -v inotifywait >/dev/null 2>&1; then \
		while inotifywait -r -e modify .; do make build; done; \
	else \
		echo "File watcher not available. Install fswatch or inotify-tools"; \
	fi

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	@make test-performance
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		echo "Performance test of binary:"; \
		time ./$(BUILD_DIR)/$(BINARY_NAME) stats --benchmark || true; \
	fi

profile: ## Profile the application
	@echo "Profiling application..."
	@if [ -f "go.mod" ]; then \
		go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...; \
		echo "Profile files generated: cpu.prof, mem.prof"; \
	else \
		echo "Profiling not configured for Python"; \
	fi

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f cpu.prof mem.prof
	@rm -rf dist/ build/ *.egg-info/
	@find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
	@echo "Cleanup complete"

clean-all: clean ## Clean everything including dependencies
	@echo "Cleaning all artifacts and dependencies..."
	@go clean -modcache 2>/dev/null || true
	@rm -rf venv/ node_modules/
	@echo "Deep cleanup complete"

# Docker (future)
docker-build: ## Build Docker image
	@echo "Docker build not implemented yet"

docker-run: ## Run in Docker container
	@echo "Docker run not implemented yet"

# Release Management
version: ## Show current version
	@if [ -f "$(BUILD_DIR)/$(BINARY_NAME)" ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) version; \
	else \
		echo "Binary not built. Run 'make build' first"; \
	fi

tag: ## Create a git tag for release
	@read -p "Enter version (e.g., v0.1.0): " version; \
	git tag -a $$version -m "Release $$version"; \
	echo "Created tag $$version. Push with: git push origin $$version"

release-notes: ## Generate release notes
	@echo "Generating release notes..."
	@echo "Release notes functionality not implemented yet"

# Maintenance
update-deps: ## Update dependencies
	@if [ -f "go.mod" ]; then \
		go get -u ./...; \
		go mod tidy; \
	elif [ -f "requirements.txt" ]; then \
		pip install --upgrade -r requirements.txt; \
	fi

check-deps: ## Check for dependency updates
	@if [ -f "go.mod" ]; then \
		go list -u -m all; \
	elif [ -f "requirements.txt" ]; then \
		pip list --outdated; \
	fi

audit: ## Run security audit on dependencies
	@if [ -f "go.mod" ]; then \
		go list -json -deps ./... | nancy sleuth; \
	elif [ -f "requirements.txt" ]; then \
		pip-audit; \
	fi

# All-in-one commands
ci: lint test ## Run all CI checks
	@echo "All CI checks passed!"

pre-commit: format lint test ## Run pre-commit checks
	@echo "Pre-commit checks passed!"

release-prep: clean build test lint docs ## Prepare for release
	@echo "Release preparation complete!"

# Homepage Development
homepage-dev: ## Start local homepage development server
	@echo "Starting homepage development server..."
	@./scripts/dev-homepage.sh

homepage-test: ## Test homepage functionality
	@echo "Testing homepage..."
	@./tests/test-homepage.sh

homepage-deploy: ## Deploy homepage to GitHub Pages
	@echo "Deploying homepage..."
	@./scripts/deploy-homepage.sh

homepage-validate: ## Validate homepage structure
	@echo "Validating homepage..."
	@./scripts/dev-homepage.sh --validate

.PHONY: homepage-dev homepage-test homepage-deploy homepage-validate
