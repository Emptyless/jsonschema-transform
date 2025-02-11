MAKEFLAGS := --no-print-directory --silent

default: help

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z\._-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.PHONY: help

fmt: ## Format go code & tidy the go.mod file
	@gofumpt -l -w .
	go mod tidy
.PHONY: fmt

t: test
test: ## Run unit tests, alias: t
	go test ./...
.PHONY: test

ci: fmt test ## simulate pipeline checks
.PHONY: ci

tools: ## Install extra tools for development
	go install mvdan.cc/gofumpt@latest
	brew install d2
.PHONY: tools

