VERSION := $(shell cat VERSION)
BIN := anka-prometheus-exporter
ARCH := amd64
OS_TYPE ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')

# CGO_ENABLED=0 needed to fix "sh: anka-prometheus-exporter: not found" in docker
build:
	GOARCH=$(ARCH) go build $(RACE) -ldflags "-X main.version=$(VERSION)" -o bin/$(BIN)_$(OS_TYPE)_$(ARCH)
	chmod +x bin/$(BIN)_$(OS_TYPE)_$(ARCH)

build-and-run:
	kill -15 $$(pgrep "[a]nka-prometheus") || true
	$(MAKE) build-mac
	./bin/$(BIN)_$(OS_TYPE)_$(ARCH) --controller-address http://anka.controller:8090

clean:
	rm -f $(BIN)*
	rm -f ./bin/$(BIN)*
	rm -f docker/scratch/$(BIN)_*
	
build-linux:
	GOOS=linux OS_TYPE=linux $(MAKE) build

build-mac:
	GOOS=darwin $(MAKE) build RACE="-race"

docker-build-and-push-scratch:
	cp -f ./bin/$(BIN)_linux_$(ARCH) docker/scratch/
	cd ./docker/scratch && \
	docker buildx build --platform linux/amd64 -t veertu/$(BIN):latest -t veertu/$(BIN):v$(VERSION) --push .
