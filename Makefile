.PHONY: lint test yaegi_test vendor clean

export GO111MODULE=on

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: fmt lint test yaegi_test

lint:
	golangci-lint run

fmt:
	gofmt -l -w $(SRC)

test:
	go test -race -cover ./...

yaegi_test:
	yaegi test .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor
