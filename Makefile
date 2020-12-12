package := $(shell basename `pwd`)

.PHONY: default get codetest build setup test fmt lint vet

default: fmt codetest

get:
	GOOS=windows GOARCH=amd64 go get -v ./...
	go get github.com/akavel/rsrc
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.31.0

codetest: lint vet test

build:
	mkdir -p target
	rm -f target/*
	$(shell go env GOPATH)/bin/rsrc -arch amd64 -manifest gps-qth-qtr.manifest -ico gps-qth-qtr.ico -o gps-qth-qtr.syso
	GOOS=windows GOARCH=amd64 go build -v -ldflags "-s -w -H=windowsgui" -o target/$(package).exe

setup: default build
	cp $(package).yaml target/

test:
	go test

fmt:
	GOOS=windows GOARCH=amd64 go fmt ./...

lint:
	GOOS=windows GOARCH=amd64 $(shell go env GOPATH)/bin/golangci-lint run --fix

vet:
	GOOS=windows GOARCH=amd64 go vet -all .