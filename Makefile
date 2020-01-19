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
	GOOS=windows GOARCH=amd64 go build -v -o target/$(package).exe

test:
	GOOS=windows GOARCH=amd64 go test -v -cover
	
fmt:
	GOOS=windows GOARCH=amd64 go fmt ./...

lint:
	GOOS=windows GOARCH=amd64 golangci-lint run --fix
	
vet:
	GOOS=windows GOARCH=amd64 go vet -all .