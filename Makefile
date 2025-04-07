# Makefile for PacTUI

shell := /bin/bash

install:
	go install src/pactui.go GOBIN=/usr/bin/
