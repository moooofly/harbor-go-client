GOPATH ?= $(shell go env GOPATH)

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

GO         := go
PKGS       := $(shell $(GO) list ./... | grep -v vendor)
# NOTE: '-race' requires cgo; enable cgo by setting CGO_ENABLED=1
#BUILD_FLAG := -race
GOBUILD    := CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)

LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.ClientVersion=$(shell cat VERSION)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GoVersion=$(shell go version)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.UTCBuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GitTag=$(shell git describe --tags)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GitHash=$(shell git rev-parse HEAD)"

.PHONY: all build install lint test pack docker misspell shellcheck clean

all: lint build test

build:
	@echo "==> Building ..."
	$(GOBUILD) -ldflags '$(LDFLAGS)' ./
	@echo ""

install:
	@echo "==> Installing ..."
	$(GO) install -x ${SRC}
	@echo ""

lint:
	@# gometalinter - Concurrently run Go lint tools and normalise their output
	@# - go get -u github.com/alecthomas/gometalinter  (Install from HEAD)
	@# - gometalinter --install  (Install all known linters)
	@echo "==> Running gometalinter ..."
	gometalinter --exclude=vendor --disable-all --enable=golint --enable=vet --enable=gofmt --enable=misspell ./...
	find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -w -s
	@echo ""

test:
	@echo "==> Testing ..."
	$(GO) test -short -race $(PKGS)
	@echo ""

pack: build
	@echo "==> Packing ..."
	@tar czvf $(shell cat VERSION)-bin.tar.gz harbor-go-client conf/*.yaml
	@echo ""
	@rm harbor-go-client
	@echo ""

docker:
	@echo "==> Creating docker image ..."
	docker build -t harbor-go-client:$(shell git rev-parse --short HEAD) .
	@echo ""

misspell:
	@# misspell - Correct commonly misspelled English words in source files
	@#    go get -u github.com/client9/misspell/cmd/misspell
	@echo "==> Runnig misspell ..."
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' | xargs misspell -error
	@echo ""

shellcheck:
	@# apt-get install -y shellcheck
	@echo "==> Runnig shellcheck ..."
	shellcheck ./scripts/*.sh
	@echo ""

clean:
	@echo "==> Cleaning ..."
	$(GO) clean -x -i ${SRC}
	rm -rf *.out
	rm -rf *.tar.gz
	@echo ""

