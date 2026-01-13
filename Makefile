.PHONY: run migrate seed build clean

ENV_FILE=.env.local
APP_NAME=app
BUILD_DIR=bin

run:
	@export $$(cat $(ENV_FILE) | xargs) && go run ./cmd/api/main.go

migrate:
	@export $$(cat $(ENV_FILE) | xargs) && go run ./cmd/migrate

seed:
	@export $$(cat $(ENV_FILE) | xargs) && go run ./cmd/seed

build:
	@echo "🔨 Building binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

clean:
	@rm -rf $(BUILD_DIR)
