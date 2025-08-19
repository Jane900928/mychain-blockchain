#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

# Build flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=MyChain \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=mychaind \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

###############################################################################
###                                  Build                                  ###
###############################################################################

all: install

build:
	@echo "Building mychaind binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/mychaind ./cmd/mychaind

install: go.sum
	@echo "Installing mychaind binary..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/mychaind

go.sum: go.mod
	@echo "Ensuring dependencies have not been modified..."
	@go mod verify

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/

###############################################################################
###                                 Tests                                   ###
###############################################################################

test:
	@echo "Running tests..."
	@go test -mod=readonly ./...

test-race:
	@echo "Running tests with race detection..."
	@go test -mod=readonly -race ./...

test-cover:
	@echo "Running tests with coverage..."
	@go test -mod=readonly -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

###############################################################################
###                              Development                                ###
###############################################################################

init: build
	@echo "Initializing blockchain..."
	@./scripts/deploy.sh init

start: 
	@echo "Starting blockchain..."
	@./scripts/deploy.sh start

stop:
	@echo "Stopping blockchain..."
	@./scripts/deploy.sh stop

restart: stop start

clean-data:
	@echo "Cleaning blockchain data..."
	@./scripts/deploy.sh clean

status:
	@echo "Checking blockchain status..."
	@./scripts/deploy.sh status

###############################################################################
###                               Frontend                                  ###
###############################################################################

install-js:
	@echo "Installing JavaScript dependencies..."
	@npm install

build-js:
	@echo "Building TypeScript client..."
	@npm run build

serve-explorer:
	@echo "Starting blockchain explorer..."
	@npm run explorer

###############################################################################
###                                Docker                                   ###
###############################################################################

docker-build:
	@echo "Building Docker image..."
	@docker build -t mychain:latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -d --name mychain-node \
		-p 26656:26656 \
		-p 26657:26657 \
		-p 1317:1317 \
		-p 9090:9090 \
		-p 9091:9091 \
		mychain:latest

docker-stop:
	@echo "Stopping Docker container..."
	@docker stop mychain-node || true
	@docker rm mychain-node || true

###############################################################################
###                              Utilities                                 ###
###############################################################################

format:
	@echo "Formatting Go code..."
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s

lint:
	@echo "Running linter..."
	@golangci-lint run

proto-gen:
	@echo "Generating protobuf files..."
	@buf generate

###############################################################################
###                               Help                                     ###
###############################################################################

help:
	@echo "Available commands:"
	@echo ""
	@echo "Build commands:"
	@echo "  build          - Build the mychaind binary"
	@echo "  install        - Install the mychaind binary"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "Test commands:"
	@echo "  test           - Run all tests"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-cover     - Run tests with coverage report"
	@echo ""
	@echo "Development commands:"
	@echo "  init           - Initialize a new blockchain"
	@echo "  start          - Start the blockchain"
	@echo "  stop           - Stop the blockchain"
	@echo "  restart        - Restart the blockchain"
	@echo "  clean-data     - Clean all blockchain data"
	@echo "  status         - Check blockchain status"
	@echo ""
	@echo "Frontend commands:"
	@echo "  install-js     - Install JavaScript dependencies"
	@echo "  build-js       - Build TypeScript client"
	@echo "  serve-explorer - Start blockchain explorer"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-stop    - Stop Docker container"
	@echo ""
	@echo "Utility commands:"
	@echo "  format         - Format Go code"
	@echo "  lint           - Run linter"
	@echo "  proto-gen      - Generate protobuf files"

.PHONY: all build install clean test test-race test-cover init start stop restart clean-data status install-js build-js serve-explorer docker-build docker-run docker-stop format lint proto-gen help
