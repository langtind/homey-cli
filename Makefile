.PHONY: build test fmt lint tools

build:
	go build -o homeyctl .

test:
	go test ./...

tools:
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

fmt:
	gofumpt -w .

lint:
	golangci-lint run
