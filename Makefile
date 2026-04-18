GOPATH := $(shell go env GOPATH)

.PHONY: setup
setup:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: lint
lint:
	golangci-lint run
	cd tui && golangci-lint run

.PHONY: test
test:
	go test ./...
	cd tui && go test ./...

.PHONY: release
release:
ifndef VERSION
	$(error VERSION is required. Usage: make release VERSION=v0.1.0)
endif
	git tag $(VERSION)
	git tag tui/$(VERSION)
	git push origin $(VERSION) tui/$(VERSION)

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix
	cd tui && golangci-lint run --fix

.PHONY: format
format:
	golangci-lint fmt
	cd tui && golangci-lint fmt
