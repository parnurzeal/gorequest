GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /examples/)
GOFILES := $(shell find . -name "*.go")

.PHONY: init
init:
	@if [ $(GO_VERSION) -gt 15 ]; then \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0; \
	elif [ $(GO_VERSION) -lt 16 ]; then \
		$(GO) get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1; \
	fi

.PHONY: dep
dep:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: lint
lint:
	export GOFLAGS=-mod=vendor
	golangci-lint run

.PHONY: test
test:
	$(GO) test -mod=vendor . -covermode=count -coverprofile .coverage.cov
