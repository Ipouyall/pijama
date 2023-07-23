.PHONY: all deps install test

all: deps test

deps:
	@go get -u ./...

install: deps
	go install ./...

test:
	@go test -v ./...
