GO ?= go
OS = $(shell uname -s | tr A-Z a-z)
SERVICE_NAME = $(shell basename "$(PWD)")
ROOT = $(shell pwd)
export GOBIN = ${ROOT}/bin

PATH := $(PATH):$(GOBIN)

ENV_FILE = .env
EXPORT_ENV = export $(shell test -e ./$(ENV_FILE) || cp ./.env.example ./$(ENV_FILE) && grep -v '^#' ./$(ENV_FILE) | xargs -d '\n')

TPARSE = $(GOBIN)/tparse
TPARSE_DOWNLOAD = $(GO) install github.com/mfridman/tparse@latest

BUILD_OUTPUT = ./bin/${SERVICE_NAME}

.PHONY: help
help: ## Display help message
	@ cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: mod
mod: ## Get dependency packages
	@ $(GO) mod tidy

.PHONY: test-base
test-base: ## Run base of test
	@ test -e $(TPARSE) || $(TPARSE_DOWNLOAD)
	@ $(TPARSE) --version

.PHONY: test
test:test-base ## Run test
	@ $(EXPORT_ENV) && $(GO) test -timeout 1000s -short ./internal/... -json -cover | $(TPARSE) -all -smallscreen

.PHONY: build
build: ## Build app binary file
	@ $(GO) env -w GO111MODULE=on
	@ $(GO) env -w CGO_ENABLED=0
	@ $(GO) build -o ${BUILD_OUTPUT} ./cmd/search/*.go

.PHONY: run
run: ## Run app
	@ $(EXPORT_ENV) && $(GO) run ./cmd/search/*.go
