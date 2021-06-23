# Makefile based on https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705

GOCMD=go
GOTEST=$(GOCMD) test
PROJECT_NAME := $(shell basename "$(PWD)")
BINARY_NAME?=$(PROJECT_NAME)
BIN_DIR?=bin
TOOLS_DIR?=bin/tools

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

GOLANGCILINT_VERSION = 1.40.1

.PHONY: all test build vendor

all: help

## Build:
build: ## Build your project and put the output binary in bin/out
	mkdir -p $(BIN_DIR)
	$(GOCMD) build -o bin/out/$(BINARY_NAME) .

## Linting:
lint: $(TOOLS_DIR)/golangci-lint ## Run linters
	$(TOOLS_DIR)/golangci-lint run

## Test:
test: ## Run the tests of the project
	$(GOTEST) -v -race ./...

coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(TOOLS_DIR):
	mkdir -p $(TOOLS_DIR)

$(TOOLS_DIR)/golangci-lint: $(TOOLS_DIR)
	@wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINDIR=$(@D) sh -s v$(GOLANGCILINT_VERSION) > /dev/null 2>&1