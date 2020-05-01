#/usr/bin/env bash

NAME=lsstcp

fmt:
	goimports -l -w .

build:clean fmt
	go build -o bin/${NAME} .
clean:
	rm -rf bin/*
