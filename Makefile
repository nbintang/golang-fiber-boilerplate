.PHONY: dev migrate seed build clean

ENV_FILE := .env.local
APP_NAME := app
BUILD_DIR := bin

ifeq ($(OS),Windows_NT)
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -ExecutionPolicy Bypass -Command
AIR := air
LOAD_ENV := ./scripts/load-env.ps1 -EnvFile $(ENV_FILE);
else
SHELL := /bin/bash
.SHELLFLAGS := -lc
AIR := air
LOAD_ENV := set -a; [ -f $(ENV_FILE) ] && . $(ENV_FILE); set +a;
endif

dev:
	@echo "🚀 Running with Air..."
	@$(LOAD_ENV) $(AIR)

migrate:
	@$(LOAD_ENV) go run ./cmd/migrate

seed:
	@$(LOAD_ENV) go run ./cmd/seed

build:
	@echo "🔨 Building binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

clean:
	@rm -rf $(BUILD_DIR)
