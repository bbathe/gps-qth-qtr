package := $(shell basename `pwd`)

.PHONY: default get codetest build run test fmt lint vet

default: fmt codetest

get:
	GOOS=windows GOARCH=amd64 go get -v ./...
	go get github.com/akavel/rsrc
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.20.0

codetest: lint vet test

build: codetest
	mkdir -p target
	rm -f target/*
	rsrc -manifest gps-qth-qtr.manifest -ico gps-qth-qtr.ico -o gps-qth-qtr.syso
	GOOS=windows GOARCH=amd64 go build -v -ldflags -H=windowsgui -o target/$(package).exe
	cp $(package).yaml target/

test:
	rsrc -manifest gps-qth-qtr_test.manifest -ico gps-qth-qtr.ico -o gps-qth-qtr.syso
	GOOS=windows GOARCH=amd64 go test
	rm gps-qth-qtr.syso
	
fmt:
	GOOS=windows GOARCH=amd64 go fmt ./...

lint:
	GOOS=windows GOARCH=amd64 $(shell go env GOPATH)/bin/golangci-lint run --fix
	
vet:
	GOOS=windows GOARCH=amd64 go vet -all .