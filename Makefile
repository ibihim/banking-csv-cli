APP_NAME := banking
APP_CMD_PATH := cmd/$(APP_NAME)
BUILD_OUTPUT := build/$(APP_NAME)
DB_FILE := transactions.db
CSV_FILE := data/transactions.csv

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_OUTPUT)
	@rm -rf $(DB_FILE)

.PHONY: build
build: clean
	@echo "Building the app..."
	@go build -o $(BUILD_OUTPUT) $(APP_CMD_PATH)/main.go

.PHONY: migrate
migrate:
	@echo "Setting up the app..."
	@./$(BUILD_OUTPUT) db migrate
	@./$(BUILD_OUTPUT) db load --filename $(CSV_FILE)

.PHONY: setup
setup: build migrate

.PHONY: run
run:
	@echo "Running the app..."
	@./$(BUILD_OUTPUT) app

.PHONY: lint
lint:
	@echo "Running linter..."
	@command -v golangci-lint > /dev/null || (echo "golangci-lint not found. Please install it: https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: watch
watch:
	@echo "Watching for file changes..."
	@command -v entr > /dev/null || (echo "entr not found. Please install it: http://entrproject.org" && exit 1)
	@find . -name "*.go" | entr -r make test run

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .


