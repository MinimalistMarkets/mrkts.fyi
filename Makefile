.PHONY: build run test

build:
	@go build -o bin/mrkts.fyi

run: build
	@./bin/mrkts.fyi

test:
	@go test -v ./...