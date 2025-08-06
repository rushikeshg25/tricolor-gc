# Makefile for tricolor-gc project

.PHONY: all test build clean help

all: build

help:
	@echo "Makefile targets:"
	@echo "  make build   - Build the project binary"
	@echo "  make test    - Run all Go tests"
	@echo "  make clean   - Remove built binaries"

build:
	go build -o tricolor-gc main.go gc.go

test:
	go test -v ./...

clean:
	rm -f tricolor-gc
