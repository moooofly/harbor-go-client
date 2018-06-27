GOPATH ?= $(shell go env GOPATH)

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif


GO        := go
GOBUILD   := CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)


LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.ClientVersion=$(shell cat VERSION)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GoVersion=$(shell go version)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.UTCBuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GitTag=$(shell git describe --tags)"
LDFLAGS += -X "github.com/moooofly/harbor-go-client/utils.GitHash=$(shell git rev-parse HEAD)"

all: lint build test

build:
	$(GOBUILD) -ldflags '$(LDFLAGS)' ./

install:
	$(GO) install -x ${SRC}

lint:
	gometalinter --exclude=vendor --disable-all --enable=golint --enable=vet --enable=gofmt --enable=misspell ./...
	find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -w -s

test:
	$(GO) test ${SRC}

misspell:
	# misspell - requires that the following be run first:
	#    go get -u github.com/client9/misspell/cmd/misspell
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' | xargs misspell -error

clean:
	$(GO) clean -x -i ${SRC}
	rm -rf *.out

