OUTPUT=go-sre-reports
GOPRIVATE=github.com/vend
GOFLAGS=-mod=vendor

COMMIT=$(shell git describe --always --abbrev=7 | cut -c 1-7)
.PHONY: all clean docker

#if make is typed with no further arguments, then show a list of available targets
default:
	@awk -F\: '/^[a-z_]+:/ && !/default/ {printf "- %-20s %s\n", $$1, $$2}' Makefile

all: $(OUTPUT)

vendor:
	go mod tidy
	go mod vendor

build:
	go env -w GOPRIVATE=$(GOPRIVATE)
	go env -w GOFLAGS=$(GOFLAGS)
	go build
	go env -u GOFLAGS

clean: #TODO