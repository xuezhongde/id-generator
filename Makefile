GOCMD=go
GOBUILD=go build

default: build

build: build-linux

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/id-gen cmd/main.go