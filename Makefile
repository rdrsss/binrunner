# binrunner makefile
#
# Author: Manuel A. Rodriguez (manuel.rdrs@gmail.com)
# 	
# Targets:
# 	all: builds installs tests
# 	build: builds the code
# 	install: installs the code to the GOPATH
#	test: runs tests
#	bench: runs bench tests
# 	clean: cleans the code
# 	fmt: formats source files
#

BIN_NAME := binrunner 
BIN_DIR := ./bin
SOURCE_DIR := ./src

CLI_NAME := cli
CLI_SOURCE_DIR := ./cli

GOFLAGS ?= $(GOFLAGS:)

.PHONY: all build install test bench clean fmt

all: install test

build: 
	@go build $(GOFLAGS) -o $(BIN_DIR)/$(BIN_NAME) $(SOURCE_DIR)/...
	@go build $(GOFLAGS) -o $(BIN_DIR)/$(CLI_NAME) $(CLI_SOURCE_DIR)/...

install:
	@go get $(GOFLAGS) $(SOURCE_DIR)/...

test: install
	@go test $(GOFLAGS) $(SOURCE_DIR)

bench: install
	@go test -run=NONE -bench=. $(GOFLAGS) $(SOURCE_DIR)

clean: 
	@go clean $(GOFLAGS) -i $(BIN_DIR)/...

fmt:
	@go fmt $(GOFLAGS) $(SOURCE_DIR)/...
	@go fmt $(GOFLAGS) $(CLI_SOURCE_DIR)/...

