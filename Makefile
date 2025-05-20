GOPATH := $(shell go env GOPATH)

.PHONY: setup
setup:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: lint
lint:
	golangci-lint run

.PHONY: format
format:
	golangci-lint fmt
