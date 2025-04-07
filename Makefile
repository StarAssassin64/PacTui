# Makefile for PacTUI

DESTDIR = /usr/bin/
shell := /bin/bash

install:
	go install src/pactui.go GOBIN=$DESTDIR

build:
	go build src/pactui

clean:
	rm -rf pactui
	go clean
