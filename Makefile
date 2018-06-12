all: lint build test

build:
	go build -x ./

install:
	go install -x ${SRC}

lint:
	gometalinter --exclude=vendor --exclude=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -w -s

test:
	go test ${SRC}

misspell:
	find . -name '*.go' -not -path './vendor/*' -not -path './_repos/*' | xargs misspell -error

clean:
	go clean -x -i ${SRC}

