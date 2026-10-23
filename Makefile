# Set by Jenkins Version param / release tag. Local: make VERSION=4.3.6 ...
VERSION ?= dev
BIN := anka-prometheus-exporter
GOLANGCI_LINT_VERSION := v2.14.0
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint
ARCH ?= $(shell arch)
ifeq ($(ARCH), i386)
	ARCH = amd64
endif
ifeq ($(ARCH), x86_64)
    ARCH = amd64
endif
OS_TYPE ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

all: clean go.releaser

#setup:		@ Enable repo git hooks (required once per clone for pre-commit lint)
setup:
	git config --local core.hooksPath .githooks/
	@echo "Git hooks enabled (core.hooksPath=.githooks/). Pre-commit runs: make go.lint"

# CGO_ENABLED=0 needed to fix "sh: anka-prometheus-exporter: not found" in docker
go.build:
	GOARCH=$(ARCH) go build $(RACE) -ldflags "-X main.version=$(VERSION)" -o bin/$(BIN)_$(OS_TYPE)_$(ARCH)
	chmod +x bin/$(BIN)_$(OS_TYPE)_$(ARCH)

#go.test:		@ Run `go test` against the current code
go.test:
	@echo "=========================================="
	@echo "Running Go Tests"
	@echo "=========================================="
	go test -v -count=1 ./...
	@echo "=========================================="
	@echo "All tests passed!"
	@echo "=========================================="

#go.lint:		@ Run `golangci-lint run` against the current code
go.lint:
	go vet ./...
	@if ! "$(GOLANGCI_LINT)" version 2>/dev/null | grep -qF "$(GOLANGCI_LINT_VERSION:v%=%)"; then \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi
	"$(GOLANGCI_LINT)" run

#go.releaser 	@ Build + package only (no GitHub publish). Jenkins uses this; GHA publishes on tag push.
go.releaser:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "ERROR: VERSION is required for go.releaser (example: make VERSION=4.3.6 go.releaser)"; \
		exit 1; \
	fi
	git tag -d "v$(VERSION)" 2>/dev/null || true
	git tag -a "v$(VERSION)" -m "v$(VERSION)"
	echo "LATEST TAG: $$(git describe --tags --abbrev=0)"
	curl -SfL https://goreleaser.com/static/run | DISTRIBUTION=oss VERSION=v2.9.0 bash -s -- release --clean --skip=publish

build-and-run:
	kill -15 $$(pgrep "[a]nka-prometheus") || true
	$(MAKE) build-mac
	./bin/$(BIN)_$(OS_TYPE)_$(ARCH) --controller-address http://anka.controller:8090 $(ARGUMENTS)

clean:
	rm -f $(BIN)*
	rm -f ./bin/$(BIN)*
	rm -f docker/scratch/$(BIN)_*

build-linux:
	GOOS=linux OS_TYPE=linux $(MAKE) go.build

build-mac:
	GOOS=darwin $(MAKE) go.build RACE="-race"

docker-build-scratch:
	$(MAKE) build-linux
	mv ./bin/$(BIN)_linux_$(ARCH) docker/scratch/
	docker build --load ./docker/scratch -t $(BIN)-scratch
