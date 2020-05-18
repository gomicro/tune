SHELL = bash

APP := $(shell basename $(PWD) | tr '[:upper:]' '[:lower:]')

GIT_COMMIT_HASH ?= $(shell git rev-parse HEAD)
GIT_SHORT_COMMIT_HASH := $(shell git rev-parse --short HEAD)

.PHONY: all
all: test

.PHONY: clean
clean: ## Cleans out all generated items
	-@rm -rf coverage
	-@rm -f coverage.txt

.PHONY: coverage
coverage: ## Generates the code coverage from all the tests
	docker run -v $$PWD:/go$${PWD/$$GOPATH} --workdir /go$${PWD/$$GOPATH} gomicro/gocover

.PHONY: help
help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

.PHONY: linters
linters: ## Run all the linters
	docker run -v $$PWD:/go$${PWD/$$GOPATH} --workdir /go$${PWD/$$GOPATH} gomicro/golinters

.PHONY: test
test: unit_test ## Run all available tests

.PHONY: unit_test
unit_test: ## Run unit tests
	go test ./...
