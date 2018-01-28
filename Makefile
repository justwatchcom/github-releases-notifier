DIST := dist
BIN := bin

EXECUTABLE := github-releases-notifier

PWD := $(shell pwd)
VERSION := $(shell cat VERSION)
SHA := $(shell cat COMMIT 2>/dev/null || git rev-parse --short=8 HEAD)
DATE := $(shell date -u '+%FT%T%z')

GOLDFLAGS += -X "main.version=$(VERSION)"
GOLDFLAGS += -X "main.date=$(DATE)"
GOLDFLAGS += -X "main.commit=$(SHA)"
GOLDFLAGS += -extldflags '-static'

GO := CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
GOARM := CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go

GOOS ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell go version | cut -d' ' -f4 | cut -d'/' -f2)

PACKAGES ?= $(shell go list ./... | grep -v /vendor/ | grep -v /tests)

TAGS ?= netgo

.PHONY: all
all: clean test build buildarm

.PHONY: clean
clean:
	$(GO) clean -i ./...
	find . -type f -name "coverage.out" -delete
	if [ -f Dockerfile ]; then rm Dockerfile ; fi

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
	STATUS=0; for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || STATUS=1; done; exit $$STATUS

.PHONY: errcheck
errcheck:
	@which errcheck > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/kisielk/errcheck; \
	fi
	STATUS=0; for PKG in $(PACKAGES); do errcheck $$PKG || STATUS=1; done; exit $$STATUS

.PHONY: megacheck
megacheck:
	@which megacheck > /dev/null; if [ $$? -ne 0  ]; then \
		$(GO) get -u honnef.co/go/tools/cmd/megacheck; \
	fi
	STATUS=0; for PKG in $(PACKAGES); do megacheck $$PKG || STATUS=1; done; exit $$STATUS

.PHONY: test
test:
	STATUS=0; for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || STATUS=1; done; exit $$STATUS

.PHONE: buildlinux
buildlinux:
	$(GO) build -tags '$(TAGS)' -ldflags '-s -w $(GOLDFLAGS)' -o $(EXECUTABLE)

.PHONY: buildarm
buildarm:
	$(GOARM) build -tags '$(TAGS)' -ldflags '-s -w $(GOLDFLAGS)' -o $(EXECUTABLE)

.PHONY: build
build: buildlinux buildarm

.PHONY: releaseamd64
releaseamd64:
	if [ -f Dockerfile ]; then rm Dockerfile ; fi
	ln -s DockerfileAMD Dockerfile
	docker build . -t pcarranza/github-releases-notifier:latest
	docker push pcarranza/github-releases-notifier:latest

.PHONY: releasearm
releasearm:
	if [ -f Dockerfile ]; then rm Dockerfile ; fi
	ln -s DockerfileARM Dockerfile
	docker build . -t pcarranza/github-releases-notifier-armv6:latest
	docker tag pcarranza/github-releases-notifier-armv6:latest \
		pcarranza/github-releases-notifier-armv6:$(VERSION)
	docker push pcarranza/github-releases-notifier-armv6:latest
	docker push pcarranza/github-releases-notifier-armv6:$(VERSION)
