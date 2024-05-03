.PHONY: lint test yaegi_test vendor clean

export GO111MODULE=on

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: fmt spell lint test yaegi_test

ci: inst tidy generate default vulncheck

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

generate:
	go generate ./...

tidy:
	go mod tidy
	cd tools && go mod tidy

spell:
	misspell -error -locale=US -w **.md

inst:
	cd tools && go install $(shell cd tools && go list -e -f '{{ join .Imports " " }}' -tags=tools)

vulncheck:
	govulncheck ./...
