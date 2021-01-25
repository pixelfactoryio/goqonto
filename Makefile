.PHONY: fmt test lint
SHELL := /bin/bash

fmt:
	@diff -u <(echo -n) <(gofmt -d -s .)

test:
	@go test -v -race -coverprofile coverage.txt -covermode atomic ./...

lint:
	@golangci-lint run ./...
