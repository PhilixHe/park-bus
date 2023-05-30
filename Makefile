default: help

linux-build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/park-bus ./cmd/main.go
.PHONY: linux-build

mac-build:
	GOARCH=amd64 GOOS=darwin go build -o bin/park-bus ./cmd/main.go
.PHONY: mac-build

fmt:
	gofmt -s -l -w ./
.PHONY: fmt

run:
	./bin/park-bus -f ./config/config.yaml
.PHONY: run

help:
	@echo "make linux-build"
	@echo "make mac-build"
	@echo "make fmt"
	@echo "make run"