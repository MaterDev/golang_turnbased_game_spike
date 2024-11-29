.PHONY: all build test clean dev install-tools frontend-install frontend-dev backend-dev start-dev stop-dev

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=gameserver
BINARY_UNIX=$(BINARY_NAME)_unix

# Main package path
MAIN_PATH=./cmd/gameserver

# Frontend parameters
FRONTEND_DIR=./frontend
NPM=npm

all: test build

build: frontend-build backend-build

frontend-install:
	cd $(FRONTEND_DIR) && $(NPM) install

frontend-build: frontend-install
	cd $(FRONTEND_DIR) && $(NPM) run build

frontend-dev:
	cd $(FRONTEND_DIR) && $(NPM) run dev

backend-build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

backend-dev:
	air

test:
	$(GOTEST) -v ./...

clean: frontend-clean backend-clean

frontend-clean:
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules

backend-clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -rf ./tmp

# Start both frontend and backend in development mode
start-dev:
	@echo "Starting development servers..."
	@make frontend-install
	@(make frontend-dev &)
	@(make backend-dev &)
	@echo "Both servers are starting..."
	@echo "Frontend will be available at http://localhost:5173"
	@echo "Backend will be available at http://localhost:3000"
	@echo "Frontend will be available at http://localhost:5173"
	@echo "Backend will be available at http://localhost:3000"

# Stop all development servers
stop-dev:
	@echo "Stopping development servers..."
	@pkill -f "npm run dev" || true
	@pkill -f "air" || true
	@echo "Development servers stopped"

# Install all development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	cd $(FRONTEND_DIR) && $(NPM) install

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) $(MAIN_PATH)

# Help command
help:
	@echo "Available commands:"
	@echo "  make build              - Build both frontend and backend"
	@echo "  make frontend-install   - Install frontend dependencies"
	@echo "  make frontend-dev       - Run frontend development server"
	@echo "  make frontend-build     - Build frontend for production"
	@echo "  make backend-dev        - Run backend with hot reload"
	@echo "  make backend-build      - Build backend"
	@echo "  make start-dev          - Start both frontend and backend in development mode"
	@echo "  make stop-dev           - Stop all development servers"
	@echo "  make test               - Run backend tests"
	@echo "  make clean              - Clean both frontend and backend"
	@echo "  make install-tools      - Install development tools"
	@echo "  make build-linux        - Cross compile backend for Linux"
	@echo "  make help               - Show this help message"
