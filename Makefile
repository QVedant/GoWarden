.PHONY: build run test

build:
	go build -o bin/gowarden ./cmd/gowarden

run: build
	./bin/gowarden

test:
	go test ./...