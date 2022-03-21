ECHO_BEGIN	:= \033[93m---
ECHO_END	:= \033[0m

ORG			:= github.com/HyeonHo-Park
PROJECT		:= simple-mutating-webhook
REPO_PATH	:= $(ORG)/$(PROJECT)

GO_BASE		:= $(shell pwd)
GO_BIN		:= $(GO_BASE)/bin
GO_MAIN_PKG	:= $(GO_BASE)/cmd/$(PROJECT)
GO			:= GO111MODULE=on GOFLAGS=-mod=vendor go

SLACK_WEBHOOK=https://hooks.slack.com/services/OOOO/OOOO/OOOO

.DEFAULT_GOAL	:= help

.PHONY: update-dependency
update-dependency: ## updating go dependencies.
	@echo "$(ECHO_BEGIN) Updating dependencies$(ECHO_END)"
	@$(GO) mod tidy
	@$(GO) mod vendor

.PHONY: run
run: ## compile and start the application in development mode.
	@echo "$(ECHO_BEGIN) Starting the application$(ECHO_END)"
	@$(GO) run $(GO_MAIN_PKG) $(ARGS)

.PHONY: build
build: update-dependency ## goreleaser build artifacts and container image
	@echo "$(ECHO_BEGIN) Build $(ECHO_END)"
	#@export SLACK_WEBHOOK=$(SLACK_WEBHOOK)
	@goreleaser release -f .local-goreleaser.yaml --snapshot --rm-dist
	@echo "$(ECHO_BEGIN) Yay! Binary built successfully!: $(GO_BIN)/$(PROJECT)$(ECHO_END)"

.PHONY: test
test: update-dependency ## run tests for the application.
	@echo "$(ECHO_BEGIN) Running tests$(ECHO_END)"
	@$(GO) test -v $(GO_BASE)/...
	@echo "$(ECHO_BEGIN) Hooray! All tests passed!$(ECHO_END)"

.PHONY: lint
lint: ## run go lint.
	@echo "$(ECHO_BEGIN) Running go lint$(ECHO_END)"
	@docker run --rm \
		-v $(shell pwd):/app \
		-w /app golangci/golangci-lint:v1.32.2 \
		golangci-lint run --timeout 15m --issues-exit-code 0 -E goimports,golint,godot,testpackage,misspell
	@echo "$(ECHO_BEGIN) Go! Time to check lint$(ECHO_END)"

.PHONY: clean
clean: ## remove object files and chached files.
	@-rm $(GO_BIN)/$(PROJECT) 2> /dev/null
	@$(GO) clean

.PHONY: release
release: ## goreleaser release
	@echo "$(ECHO_BEGIN) Goreleaser release$(ECHO_END)"
	#@export SLACK_WEBHOOK=$(SLACK_WEBHOOK)
	@goreleaser release --rm-dist
	@echo "$(ECHO_BEGIN) Successfully Release!$(ECHO_END)"

.PHONY: install
install: ## delpoy webhook in kubernetes
	@echo "$(ECHO_BEGIN) install webhook $(ECHO_END)"
	@./hack/install.sh
	@echo "$(ECHO_BEGIN) Successfully Installed Webhoook!$(ECHO_END)"

.PHONY: remove
remove: ## remove webhook and examples deployments in kubernetes
	@echo "$(ECHO_BEGIN) remove webhook $(ECHO_END)"
	@./hack/remove.sh
	@echo "$(ECHO_BEGIN) Successfully Removed all resources!$(ECHO_END)"

.PHONY: deploy
deploy: ## deploy examples
	@echo "$(ECHO_BEGIN) deploy examples $(ECHO_END)"
	@./hack/deploy.sh
	@echo "$(ECHO_BEGIN) Successfully Deploy examples!$(ECHO_END)"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'
