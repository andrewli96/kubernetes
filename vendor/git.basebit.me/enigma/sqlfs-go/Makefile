SHELL := /bin/bash

CURRENT_MAKEFILE_LIST := $(MAKEFILE_LIST)
MAKEFILE_DIR := $(shell cd $(dir $(firstword $(CURRENT_MAKEFILE_LIST))) && pwd)


ifneq "$(MAKEFILE_DIR)" "$(CURDIR)"
# $(error "Please change to $(MAKEFILE_DIR) then continue")
endif

SOURCES = $(shell find . -name '*.go')


ifeq ($(VERSION),)
VERSION := $(shell git describe --always --dirty)
endif


all: build

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: test
test:
	go test -v ./...


.PHONY: build
build:
	go build ./...

.PHONY: staticcheck
staticcheck:
	@staticcheck ./...

.PHONY: build-doc
build-doc:
	@cd docs && hugo

.PHONY: clean-doc
clean-doc:
	@cd docs && rm -rf public

.PHONY: run-doc-site
run-doc-site:
	@cd docs && hugo server

.PHONY: clean
clean: clean-doc
