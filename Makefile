SHELL := /bin/bash
.DEFAULT_GOAL := help

.PHONY: goimports
goimports: ## exec goimports
	goimports -w .

.PHONY: gazelle
gazelle: ## exec gazelle
	bazel run //:gazelle

.PHONY: update-repo
update-repo: ## update repository
	bazel run //:gazelle -- update-repos -from_file go.mod

.PHONY: update-deps
update-deps: ## Update deps and update repository
	go get -u -v
	go mod tidy
	make update-repo

.PHONY: test
test: ## test
	bazel test //...

.PHONY: build
build: ## build commands
	bazel build //cmd/...

.PHONY: serve
serve:
	bazel run //cmd/server

.PHONY: shutdown
shutdown: ## shutdown bazel
	bazel shutdown

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -h -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

# https://postd.cc/auto-documented-makefile/
help: ## Show help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

