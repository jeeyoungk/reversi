#!/usr/bin/env bash
GOPATH=`pwd` go fmt ./... && go run reversi.go
