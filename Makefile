all: lint build test

build:
	go build -x ./

install:
	go install -x ${SRC}

lint:
	gometalinter --exclude=vendor --disable-all --enable=golint --enable=vet --enable=gofmt --enable=misspell ./...
	find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -w -s

test:
	go test ${SRC}

misspell:
	# misspell - requires that the following be run first:
	#    go get -u github.com/client9/misspell/cmd/misspell
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' | xargs misspell -error

clean:
	go clean -x -i ${SRC}

