package := $(shell basename `pwd`)

.PHONY: default get codetest build test fmt lint vet

default: fmt codetest

get:
	go get -v ./...
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.20.0

codetest: lint vet test

build: codetest
	mkdir -p target
	rm -f target/*
	GOOS=windows GOARCH=amd64 go build -v -o target/$(package)

test:
	go test -v -cover

fmt:
	go fmt ./...

lint:
	golangci-lint run --fix
	
vet:
	go vet -all .