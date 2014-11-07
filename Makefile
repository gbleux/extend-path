default: build

all: clean test build

clean:
	rm extend-path

build:
	go build .

test:
	go test .

.PHONY: default all build test clean
