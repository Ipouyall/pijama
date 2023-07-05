.PHONY: all deps install test

all: deps test install

deps:
	go get -u ./...

install:
	go install ./...

test:
	go test -v ./...
