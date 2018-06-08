all: build

build:
	go install ${SRC}

clean:
	go clean -i ${SRC}

test:
	go test ${SRC}
