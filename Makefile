PACKAGE = github.com/connesc/streamconv
BIN = streamconv

.PHONY: build

build:
	go build '${PACKAGE}/cmd/${BIN}'

clean:
	rm -f '${BIN}'
