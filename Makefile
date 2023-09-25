VERSION := $(shell cat VERSION)
BIN := anka-prometheus-exporter
ARCH ?= $(shell arch)
ifeq ($(ARCH), i386)
	ARCH = amd64
endif
ifeq ($(ARCH), x86_64)
	ARCH = amd64
endif
OS_TYPE ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

all: clean go.releaser

# CGO_ENABLED=0 needed to fix "sh: anka-prometheus-exporter: not found" in docker
go.build:
	GOARCH=$(ARCH) go build $(RACE) -ldflags "-X main.version=$(VERSION)" -o bin/$(BIN)_$(OS_TYPE)_$(ARCH)
	chmod +x bin/$(BIN)_$(OS_TYPE)_$(ARCH)

#go.lint:		@ Run `golangci-lint run` against the current code
go.lint:
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b /usr/local/bin v1.40.1
	golangci-lint run --fast

#go.releaser 	@ Run goreleaser release --clean for current version
go.releaser:
	git tag -d "$(VERSION)" 2>/dev/null || true
	git tag -a "$(VERSION)" -m "Version $(VERSION)"
	echo "LATEST TAG: $$(git describe --tags --abbrev=0)"
	goreleaser release --clean

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
	docker build ./docker/scratch -t $(BIN)-scratch
