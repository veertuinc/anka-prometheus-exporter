PREFIX := github.com/veertuinc/anka-prometheus
VERSION := $(shell cat VERSION)
BIN := anka-prometheus

build:
	GOOS=darwin go build -o bin/$(BIN) $(PREFIX)

build-and-run:
	kill -15 $$(pgrep "[a]nka-prometheus") || true
	$(MAKE) build
	./bin/anka-prometheus --controller_address http://anka.controller:8090