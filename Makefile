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
	$(error VERSION is required. Usage: make release VERSION=0.1.0)
endif
	git tag v$(VERSION)
	git tag tui/v$(VERSION)
	git push origin v$(VERSION) tui/v$(VERSION)

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix
	cd tui && golangci-lint run --fix

.PHONY: format
format:
	golangci-lint fmt
	cd tui && golangci-lint fmt
