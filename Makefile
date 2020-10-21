VERSION := $(shell cat VERSION)
BIN := anka-prometheus-exporter
ARCH := amd64
OS_TYPE ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

# CGO_ENABLED=0 needed to fix "sh: anka-prometheus-exporter: not found" in docker
build:
	CGO_ENABLED=0 GOARCH=$(ARCH) go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BIN)_$(OS_TYPE)_$(ARCH)
	chmod +x bin/$(BIN)_$(OS_TYPE)_$(ARCH)

build-and-run:
	kill -15 $$(pgrep "[a]nka-prometheus") || true
	$(MAKE) build-mac
	./bin/$(BIN)_$(OS_TYPE)_$(ARCH) --controller_address http://anka.controller:8090
	
build-linux:
	GOOS=linux OS_TYPE=linux $(MAKE) build

build-mac:
	GOOS=darwin $(MAKE) build

gorelease:
	git tag -d v$(VERSION) || true
	git tag -a "v$(VERSION)" -m "Version $(VERSION)"
	goreleaser release --rm-dist --debug