#!/usr/bin/env bash
GOPATH=`pwd` go fmt ./... && go test -v ./...
