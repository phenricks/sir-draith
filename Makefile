.PHONY: all build run test clean lint

# Variáveis
BINARY_NAME=sir-draith
MAIN_FILE=cmd/bot/main.go

all: clean build

build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

run:
	@go run $(MAIN_FILE)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

lint:
	@echo "Running linter..."
	@golangci-lint run

deps:
	@echo "Downloading dependencies..."
	@go mod download

tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME) .

docker-run:
	@echo "Running Docker container..."
	@docker run --rm $(BINARY_NAME)

# Development helpers
watch:
	@echo "Watching for changes..."
	@air -c .air.toml

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build files"
	@echo "  make lint       - Run linter"
	@echo "  make deps       - Download dependencies"
	@echo "  make tidy       - Tidy go.mod"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo "  make watch      - Watch for changes (requires air)"
	@echo "  make help       - Show this help" 