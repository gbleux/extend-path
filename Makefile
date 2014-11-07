clean:
	rm extend-path

build:
	go build .

test:
	go test .

.PHONY: build test clean