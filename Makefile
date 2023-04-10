APP_NAME := banking
APP_CMD_PATH := cmd/$(APP_NAME)
BUILD_OUTPUT := build/$(APP_NAME)

.PHONY: build
build:
	@echo "Building the app..."
	@go build -o $(BUILD_OUTPUT) $(APP_CMD_PATH)/main.go

.PHONY: run
run:
	@echo "Running the app..."
	@./$(BUILD_OUTPUT)

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


