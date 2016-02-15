.PHONY: all fmt vet lint test test-full
.DEFAULT: default
all: AUTHORS fmt vet lint test test-full

AUTHORS: .mailmap .git/HEAD
	 git log --format='%aN <%aE>' | sort -fu > $@

MY_REPO ?= "dmp42/scape-go"

# Package list
MY_PKGS ?= $(shell go list -tags "$MY_BUILDTAGS" ./... | grep -v ^github.com/$MY_REPO/vendor)

# Resolving binary dependencies for specific targets
GOLINT_BIN := $(GOPATH)/bin/golint
GOLINT := $(shell [ -x $(GOLINT_BIN) ] && echo $(GOLINT_BIN) || echo '')

# Required for go 1.5 to build
GO15VENDOREXPERIMENT := 1

vet:
	@echo "+ $@"
	go vet -tags "${MY_BUILDTAGS}" ${MY_PKGS}

fmt:
	@echo "+ $@"
	test -z "$$(gofmt -s -l . 2>&1 | grep -v vendor/ | tee /dev/stderr)" || \
		(echo >&2 "+ please format Go code with 'gofmt -s'" && false)

lint:
	@echo "+ $@"
	$(if $(GOLINT), , \
		$(error Please install golint: `go get -u github.com/golang/lint/golint`))
	test -z "$$($(GOLINT) ./... 2>&1 | grep -v vendor/ | tee /dev/stderr)"

test:
	@echo "+ $@"
	go test -test.short -tags "${MY_BUILDTAGS}" ./...

test-full:
	@echo "+ $@"
	go test -tags "${MY_BUILDTAGS}" ./...
