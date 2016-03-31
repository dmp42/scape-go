.PHONY: all fmt vet lint test test-full
.DEFAULT: default
all: fmt vet lint test test-full

AUTHORS: .mailmap .git/HEAD
	 git log --format='%aN <%aE>' | sort -fu > $@

PWD := $(shell pwd)
MY_REPO ?= $(notdir $(abspath $(dir $(abspath $(dir $(PWD))))))/$(notdir $(abspath $(dir $(PWD))))/$(notdir $(PWD))

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

info:
	echo $(mkfile_path)


# Package list
MY_PKGS ?= $(shell go list -tags "${MY_BUILDTAGS}" ./... | grep -v ^${MY_REPO}/vendor/)

# Resolving binary dependencies for specific targets
GOLINT_BIN := $(GOPATH)/bin/golint
GOLINT := $(shell [ -x $(GOLINT_BIN) ] && echo $(GOLINT_BIN) || echo '')

GODEP_BIN := $(GOPATH)/bin/godep
GODEP := $(shell [ -x $(GODEP_BIN) ] && echo $(GODEP_BIN) || echo '')

# Required for go 1.5 to build
GO15VENDOREXPERIMENT := 1

vet:
	@echo "+ $@"
	go vet -tags "${MY_BUILDTAGS}" ${MY_PKGS}

fmt:
	@echo "+ $@"
	test -z "$$(gofmt -s -l . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || \
		(echo >&2 "+ please format Go code with 'gofmt -s'" && false)

lint:
	@echo "+ $@"
	$(if $(GOLINT), , \
		$(error Please install golint: `go get -u github.com/golang/lint/golint`))
	test -z "$$($(GOLINT) ./... 2>&1 | grep -v ^vendor/ | tee /dev/stderr)"

test:
	@echo "+ $@"
	go test -test.short -tags "${MY_BUILDTAGS}" ${MY_PKGS}

test-full:
	@echo "+ $@"
	go test -tags "${MY_BUILDTAGS}" ${MY_PKGS}

validate-godep:
	$(if $(GODEP), , \
		$(error Please install godep: go get github.com/tools/godep))
	@$(GODEP) restore -v
	@echo "+ $@"
	@rm -Rf .vendor.bak
	@mv vendor .vendor.bak
	@rm -Rf Godeps
	@$(GODEP) save ./...
	@test -z "$$(diff -r vendor .vendor.bak 2>&1 | tee /dev/stderr)" || \
		(echo >&2 "+ borked dependencies! what you have in Godeps/Godeps.json does not match with what you have in vendor" && false)
	@rm -Rf .vendor.bak
