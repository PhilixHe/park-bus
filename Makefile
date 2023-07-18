default: help

linux-build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/park-bus ./cmd/main.go
.PHONY: linux-build

mac-build:
	GOARCH=amd64 GOOS=darwin go build -o bin/park-bus ./cmd/main.go
.PHONY: mac-build

win-build:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/park-bus.exe ./cmd/main.go
.PHONY: win-build

fmt:
	gofmt -s -l -w ./
.PHONY: fmt

run:
	./bin/park-bus -f ./config/config.yaml
.PHONY: run

help:
	@echo "make linux-build"
	@echo "make mac-build"
	@echo "make win-build"
	@echo "make fmt"
	@echo "make run"