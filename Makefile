# circulator Makefile

.PHONY: build build-all run test clean docker proto help fmt fmt-strict fmt-imports fmt-all fmt-check fmt-fix fmt-keep-align check check-quick format-check-ci tools-install tools-check tools-clean workspace-init project-check project-list project-install lint-fix lint-strict
.DEFAULT_GOAL := help

# Variables
BINARY_DIR := bin
GO_VERSION := $(shell go version | cut -d' ' -f3)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Environment management (simplified)
env-up: ## Start all services (MySQL, Redis, Pulsar)
	@echo "Starting all services..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d; \
	else \
		docker compose up -d; \
	fi
	@echo "Waiting for services to be ready..."
	sleep 15
	@echo "Services started. Use 'make status' to check health."

env-down: ## Stop all services
	@echo "Stopping all services..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose stop; \
	else \
		docker compose stop; \
	fi
	@echo "All services stopped."

env-clean: ## Stop and remove all containers and volumes
	@echo "Cleaning up all containers and volumes..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down -v; \
	else \
		docker compose down -v; \
	fi
	@echo "Environment cleaned up."

# Legacy aliases for backward compatibility
dev-env: env-up ## Alias for env-up (deprecated, use env-up)

dev-full: env-up dev-server ## Start full development environment with serverhell go version | cut -d ' ' -f 3)
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

s:
	git add .
	commit-emoji
	git push origin main

# Help target
help: ## Show this help message
	@echo "Circulator - Available make targets:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Protocol Buffers
proto: ## Generate protobuf files
	@echo "Generating protobuf files..."
	@mkdir -p pkg/agent/gengrpc
	@if [ -d "pkg/entity/proto" ] && [ -n "$$(find pkg/entity/proto -name '*.proto' -type f 2>/dev/null)" ]; then \
		echo "Found .proto files, generating gRPC code..."; \
		protoc --go_out=pkg/agent/gengrpc --go_opt=paths=source_relative \
			--go-grpc_out=pkg/agent/gengrpc --go-grpc_opt=paths=source_relative \
			--proto_path=pkg/entity/proto pkg/entity/proto/*.proto; \
		echo "Proto files generated successfully."; \
	else \
		echo "No .proto files found in pkg/entity/proto, skipping generation."; \
		echo "Create .proto files in pkg/entity/proto/ to enable gRPC code generation."; \
	fi

proto-force: ## Force generate protobuf files (will fail if no .proto files exist)
	@echo "Force generating protobuf files..."
	@mkdir -p pkg/agent/gengrpc
	protoc --go_out=pkg/agent/gengrpc --go_opt=paths=source_relative \
		--go-grpc_out=pkg/agent/gengrpc --go-grpc_opt=paths=source_relative \
		--proto_path=pkg/entity/proto pkg/entity/proto/*.proto

proto-clean: ## Clean generated protobuf files
	@echo "Cleaning generated protobuf files..."
	@rm -rf pkg/agent/gengrpc/*.pb.go

proto-install: ## Install protobuf compiler and Go plugins
	@echo "Installing protobuf compiler..."
	@if ! command -v protoc >/dev/null 2>&1; then \
		echo "protoc not found. Installing via Homebrew..."; \
		brew install protobuf; \
	fi
	@echo "Installing Go protobuf plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Build targets
build-server: proto ## Build server binary (with proto generation)
	@echo "Building server..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BINARY_DIR)/server cmd/server/main.go

build-client: proto ## Build client binary (with proto generation)
	@echo "Building client..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BINARY_DIR)/client cmd/client/main.go

build-agent: proto ## Build agent binary (with proto generation)
	@echo "Building agent..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BINARY_DIR)/agent cmd/agent/main.go

build-simulator: proto ## Build simulator binary (with proto generation)
	@echo "Building simulator..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BINARY_DIR)/simulator cmd/simulator/main.go

build: build-server build-client build-agent build-simulator ## Build all binaries (with proto generation)

build-all: clean build ## Clean and build all binaries (with proto generation)

# Build without proto generation (for development speed)
build-fast: ## Build all binaries without proto generation (faster for development)
	@echo "Building all binaries (skipping proto generation)..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BINARY_DIR)/server cmd/server/main.go
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BINARY_DIR)/client cmd/client/main.go
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BINARY_DIR)/agent cmd/agent/main.go
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" -o $(BINARY_DIR)/simulator cmd/simulator/main.go

build-server-fast: ## Build server binary without proto generation
	@echo "Building server (fast mode)..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BINARY_DIR)/server cmd/server/main.go

build-agent-fast: ## Build agent binary without proto generation
	@echo "Building agent (fast mode)..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-X main.version=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BINARY_DIR)/agent cmd/agent/main.go

# Run targets
run-server: build-server ## Build and run server
	@echo "Starting server..."
	./$(BINARY_DIR)/server

run-client: build-client ## Build and run client
	@echo "Starting client..."
	./$(BINARY_DIR)/client

run-agent: build-agent ## Build and run agent
	@echo "Starting agent..."
	./$(BINARY_DIR)/agent

run-simulator: build-simulator ## Build and run simulator
	@echo "Starting simulator..."
	./$(BINARY_DIR)/simulator

# Development targets
dev-server: ## Run server in development mode (without build)
	@echo "Running server in development mode..."
	go run cmd/server/main.go

dev-client: ## Run client in development mode (without build)
	@echo "Running client in development mode..."
	go run cmd/client/main.go

dev-agent: ## Run agent in development mode (without build)
	@echo "Running agent in development mode..."
	go run cmd/agent/main.go

dev-simulator: ## Run simulator in development mode (without build)
	@echo "Running simulator in development mode..."
	go run cmd/simulator/main.go

# Development environment setup
dev-env: ## Start development environment (databases and Pulsar)
	@echo "Starting development environment..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d mysql redis pulsar; \
	else \
		docker compose up -d mysql redis pulsar; \
	fi
	@echo "Waiting for services to be ready..."
	sleep 15

dev-env-clean: ## Clean and restart development environment
	@echo "Cleaning and restarting development environment..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down; \
		docker-compose up -d mysql redis pulsar; \
	else \
		docker compose down; \
		docker compose up -d mysql redis pulsar; \
	fi
	@echo "Waiting for services to be ready..."
	sleep 15

dev-env-reset: ## Reset development environment (WARNING: This will delete all data)
	@echo "Resetting development environment (deleting all data)..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down -v --remove-orphans; \
		docker-compose up -d mysql redis pulsar; \
	else \
		docker compose down -v --remove-orphans; \
		docker compose up -d mysql redis pulsar; \
	fi
	@echo "Waiting for services to be ready..."
	sleep 15

dev-env-stop: ## Stop development environment
	@echo "Stopping development environment..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose stop; \
	else \
		docker compose stop; \
	fi

dev-env-logs: ## Show development environment logs
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose logs -f mysql redis pulsar; \
	else \
		docker compose logs -f mysql redis pulsar; \
	fi

dev-env-status: ## Show development environment status
	@echo "=== Development Environment Status ==="
	@echo ""
	@echo "=== Docker Containers ==="
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose ps; \
	else \
		docker compose ps; \
	fi
	@echo ""
	@echo "=== Service Health Checks ==="
	@echo -n "MySQL: "
	@docker exec circulator-mysql mysqladmin ping -h localhost -u root -ppassword 2>/dev/null && echo "Healthy" || echo "Not available"
	@echo -n "Redis: "
	@docker exec circulator-redis redis-cli ping 2>/dev/null || echo "Not available"
	@echo -n "Pulsar: "
	@curl -f http://localhost:8080/admin/v2/clusters/standalone 2>/dev/null && echo "Healthy" || echo "Not available"

dev-env-ports: ## Show development environment port usage
	@echo "=== Development Environment Ports ==="
	@echo "MySQL:  3306"
	@echo "Redis:  6379" 
	@echo "Pulsar: 6650 (binary), 8080 (admin)"
	@echo ""
	@echo "=== Port Status ==="
	@echo -n "3306 (MySQL): "
	@lsof -i :3306 >/dev/null 2>&1 && echo "In use" || echo "Available"
	@echo -n "6379 (Redis): "
	@lsof -i :6379 >/dev/null 2>&1 && echo "In use" || echo "Available"
	@echo -n "6650 (Pulsar): "
	@lsof -i :6650 >/dev/null 2>&1 && echo "In use" || echo "Available"
	@echo -n "8080 (Pulsar Admin): "
	@lsof -i :8080 >/dev/null 2>&1 && echo "In use" || echo "Available"

dev-full: dev-env dev-server ## Start full development environment with server

# Watch targets (requires entr or similar)
watch-server: ## Watch and restart server on file changes
	@echo "Watching server files for changes..."
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | grep -E "(cmd/server|pkg)" | entr -r make dev-server; \
	else \
		echo "entr not found. Install with: brew install entr"; \
		echo "Falling back to regular dev mode..."; \
		make dev-server; \
	fi

watch-agent: ## Watch and restart agent on file changes
	@echo "Watching agent files for changes..."
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | grep -E "(cmd/agent|pkg)" | entr -r make dev-agent; \
	else \
		echo "entr not found. Install with: brew install entr"; \
		echo "Falling back to regular dev mode..."; \
		make dev-agent; \
	fi

# Test targets
test: ## Run all tests
	@echo "Running tests..."
	go test ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	go test -race ./...

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Clean targets
clean: ## Clean build artifacts and temporary files
	@echo "Cleaning build artifacts..."
	rm -rf $(BINARY_DIR)/
	rm -f coverage.out coverage.html

clean-proto: proto-clean ## Clean generated protobuf files

clean-all: clean clean-proto env-clean ## Clean everything (builds, proto, environment)
	@echo "Cleaning docker images..."
	docker system prune -f

# Docker operations
docker-build: ## Build docker image
	@echo "Building docker image..."
	docker build -t circulator .

docker-run: ## Start docker compose services
	@echo "Starting docker services..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d; \
	else \
		docker compose up -d; \
	fi

docker-stop: ## Stop docker compose services
	@echo "Stopping docker services..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down; \
	else \
		docker compose down; \
	fi

docker-logs: ## Show docker logs
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose logs -f; \
	else \
		docker compose logs -f; \
	fi

# Database operations (simplified)
db-only: ## Start only MySQL and Redis (no Pulsar)
	@echo "Starting database services only..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d mysql redis; \
	else \
		docker compose up -d mysql redis; \
	fi

# Legacy aliases for backward compatibility  
db-up: env-up ## Alias for env-up (deprecated, use env-up)
db-up-minimal: db-only ## Alias for db-only (deprecated, use db-only)
db-down: env-down ## Alias for env-down (deprecated, use env-down)
db-reset: env-clean env-up ## Reset environment (deprecated, use env-clean then env-up)

# Development
dev:
	docker-compose up -d mysql
	sleep 5
	go run cmd/server/main.go

# Database migration (when implemented)
migrate: ## Run database migrations
	@echo "Database migration not implemented yet"
	@echo "TODO: Implement database migration system"

# Dependencies and maintenance
deps: ## Download and tidy dependencies
	@echo "Managing dependencies..."
	go mod tidy
	go mod download

deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# プロジェクト専用ツール管理
tools-install: ## Install project-specific development tools
	@echo "Installing project development tools..."
	@echo "Installing goimports..."
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing gofumpt..."
	go install mvdan.cc/gofumpt@latest
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing protobuf tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "All development tools installed!"

tools-check: ## Check if required tools are installed
	@echo "Checking development tools..."
	@command -v goimports >/dev/null 2>&1 || (echo "❌ goimports not found" && exit 1)
	@command -v gofumpt >/dev/null 2>&1 || (echo "❌ gofumpt not found" && exit 1)
	@command -v golangci-lint >/dev/null 2>&1 || (echo "❌ golangci-lint not found" && exit 1)
	@command -v protoc-gen-go >/dev/null 2>&1 || (echo "❌ protoc-gen-go not found" && exit 1)
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || (echo "❌ protoc-gen-go-grpc not found" && exit 1)
	@echo "✅ All development tools are installed!"

tools-clean: ## Clean tool caches
	@echo "Cleaning tool caches..."
	go clean -cache
	go clean -modcache
	@echo "Tool caches cleaned!"

workspace-init: tools-install deps ## Initialize complete development workspace
	@echo "🚀 Development workspace initialized!"
	@echo "Available commands:"
	@echo "  make fmt          - Format code"
	@echo "  make fmt-strict   - Strict formatting"
	@echo "  make lint         - Run linter"
	@echo "  make test         - Run tests"
	@echo "  make build        - Build binaries"

project-check: ## Check project structure and setup
	@echo "🔍 Checking project setup..."
	@go run scripts/project.go check

project-list: ## List all Go files in project
	@echo "📁 Listing project files..."
	@go run scripts/project.go list

project-install: ## Install project tools using Go
	@echo "🔧 Installing tools via Go..."
	@go run scripts/project.go install

# Code quality and formatting
fmt: ## Format code with go fmt
	@echo "Formatting code with go fmt..."
	go fmt ./...

fmt-strict: ## Format code with gofumpt (stricter formatting)
	@echo "Formatting code with gofumpt..."
	@if command -v gofumpt >/dev/null 2>&1; then \
		gofumpt -w .; \
	else \
		echo "gofumpt not found. Installing..."; \
		go install mvdan.cc/gofumpt@latest; \
		gofumpt -w .; \
	fi

fmt-imports: ## Format and organize imports
	@echo "Formatting imports..."
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' | xargs goimports -w; \
	else \
		echo "goimports not found. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' | xargs goimports -w; \
	fi

fmt-all: fmt fmt-imports ## Run all formatters (go fmt + goimports)
	@echo "All formatting completed!"

fmt-check: ## Check if code is formatted properly
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files need formatting:"; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "All Go files are properly formatted!"; \
	fi

fmt-fix: ## Fix common formatting issues automatically
	@echo "Fixing formatting issues..."
	@find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' -exec gofmt -w {} \;
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' | xargs goimports -w; \
	fi
	@echo "Formatting fixes applied!"

fmt-keep-align: ## Format code but preserve manual struct alignment (warning: manual maintenance required)
	@echo "Note: This target skips go fmt on files with manually aligned structs."
	@echo "Run 'make fmt' for standard Go formatting, or manually maintain struct alignment."
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' -not -name 'logger.go' | xargs goimports -w; \
	else \
		echo "goimports not found. Installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
		find . -name '*.go' -not -path './vendor/*' -not -path './.git/*' -not -name 'logger.go' | xargs goimports -w; \
	fi
	@echo "Formatting applied (excluding manually aligned files)!"

lint: ## Lint code
	@echo "Linting code with custom configuration..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config .env/lint/.golangci.yml --verbose; \
	else \
		echo "golangci-lint not found. Install with:"; \
		echo "  brew install golangci-lint"; \
		echo "  # or"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

lint-fix: ## Lint and fix issues automatically
	@echo "Linting and fixing code issues..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config .env/lint/.golangci.yml --fix --verbose; \
	else \
		echo "golangci-lint not found. Please install it first with 'make tools-install'"; \
	fi

lint-strict: ## Strict linting with maximum indent checks
	@echo "Running strict lint checks including indent depth..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --config .env/lint/.golangci.yml --verbose --enable-all --disable=gochecknoglobals,gochecknoinits; \
	else \
		echo "golangci-lint not found. Please install it first with 'make tools-install'"; \
	fi

vet: ## Vet code
	@echo "Vetting code..."
	go vet ./...

check: fmt-all vet lint test ## Run all code quality checks (with full formatting)

check-quick: fmt vet lint ## Run quick code quality checks (basic formatting only)

format-check-ci: fmt-check vet lint ## CI-friendly format and quality check

# Installation
install: build ## Install binaries to GOPATH/bin
	@echo "Installing binaries..."
	cp $(BINARY_DIR)/* $(GOPATH)/bin/

uninstall: ## Remove installed binaries
	@echo "Removing installed binaries..."
	rm -f $(GOPATH)/bin/server $(GOPATH)/bin/client $(GOPATH)/bin/agent $(GOPATH)/bin/simulator

# Development tools
tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	@if command -v brew >/dev/null 2>&1; then \
		brew install entr; \
	else \
		echo "Homebrew not found. Please install entr manually for watch functionality."; \
	fi

# Air (hot reload) targets
air-server: ## Run server with hot reload using Air
	@if command -v air >/dev/null 2>&1; then \
		air -c .air.server.toml; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air -c .air.server.toml; \
	fi

air-agent: ## Run agent with hot reload using Air
	@if command -v air >/dev/null 2>&1; then \
		air -c .air.agent.toml; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air -c .air.agent.toml; \
	fi

# Quick development shortcuts
quick-start: build-all dev-env run-server ## Quick start: build all + start env + run server

quick-reset: dev-env-reset build-all ## Quick reset: reset environment + rebuild all

quick-clean: dev-env-clean build-fast ## Quick clean: clean environment + fast build
	
full-reset: clean-all build-all ## Full reset: clean everything and rebuild

# Development environment aliases (shorter commands)
up: dev-env ## Alias for dev-env
down: dev-env-stop ## Alias for dev-env-stop  
restart: dev-env-clean ## Alias for dev-env-clean
reset: dev-env-reset ## Alias for dev-env-reset
logs: dev-env-logs ## Alias for dev-env-logs
ps: dev-env-status ## Alias for dev-env-status
ports: dev-env-ports ## Alias for dev-env-ports

# Status and info
status: ## Show development environment status
	@echo "=== Circulator Status ==="
	@echo "Go version: $(GO_VERSION)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build time: $(BUILD_TIME)"
	@echo ""
	@echo "=== Binary Status ==="
	@ls -la $(BINARY_DIR)/ 2>/dev/null || echo "No binaries built yet"
	@echo ""
	@echo "=== Docker Status ==="
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose ps 2>/dev/null || echo "Docker compose not running"; \
	else \
		docker compose ps 2>/dev/null || echo "Docker compose not running"; \
	fi
	@echo ""
	@echo "=== Pulsar Health ==="
	@curl -f http://localhost:8080/admin/v2/clusters/standalone 2>/dev/null && echo "Pulsar: Healthy" || echo "Pulsar: Not available"

info: status ## Alias for status

# Help
.PHONY: all build test coverage clean clean-proto clean-all
.PHONY: build-client build-server build-agent build-simulator
.PHONY: dev-client dev-server dev-agent dev-simulator dev-all
.PHONY: dev-env env-up env-down env-clean
.PHONY: proto proto-gen proto-clean
.PHONY: docker docker-build docker-up docker-down docker-clean
.PHONY: help status info
