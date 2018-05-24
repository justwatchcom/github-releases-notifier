EXECUTABLE := github-releases-notifier

VERSION = $(shell cat VERSION)
SHA = $(shell cat COMMIT 2>/dev/null || git rev-parse --short=8 HEAD)
DATE = $(shell date -u '+%FT%T%z')

GO ?= go
GOLDFLAGS += -X "main.version=$(VERSION)"
GOLDFLAGS += -X "main.date=$(DATE)"
GOLDFLAGS += -X "main.commit=$(SHA)"
GOLDFLAGS += -extldflags '-static'

GOOS ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f2)

PACKAGES ?= $(shell go list ./... | grep -v /vendor/ | grep -v /tests)

TAGS ?= netgo

.PHONY: all
all: clean test build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	find . -type f -name "coverage.out" -delete

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: tests
tests: test vet lint errcheck megacheck

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/golang/lint/golint; \
	fi
	golint -set_exit_status $(PACKAGES)

.PHONY: errcheck
errcheck:
	@which errcheck > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

.PHONY: megacheck
megacheck:
	@which megacheck > /dev/null; if [ $$? -ne 0  ]; then \
		$(GO) get -u honnef.co/go/tools/cmd/megacheck; \
	fi
	megacheck $(PACKAGES)

.PHONY: test
test:
	STATUS=0;
	for PKG in $(PACKAGES); do
		$(GO) test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || STATUS=1;
	done; 
	exit $$STATUS

.PHONY: build
build: $(EXECUTABLE)-$(GOOS)-$(GOARCH)

$(EXECUTABLE)-$(GOOS)-$(GOARCH): $(wildcard *.go)
	$(GO) build -tags '$(TAGS)' -ldflags '-s -w $(GOLDFLAGS)' -o $(EXECUTABLE)
