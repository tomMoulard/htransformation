.PHONY: default ci lint test yaegi_test vendor clean generate tidy spell vulncheck build

export GO111MODULE=on

default: spell lint build test

ci: tidy generate default vulncheck

lint:
	go tool goreleaser check
	go tool golangci-lint run

test:
	go test -race -cover ./...

yaegi_test:
	go tool yaegi test .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor

generate:
	go generate ./...

tidy:
	go mod tidy

spell:
	go tool misspell -error -locale=US -w **.md

vulncheck:
	go tool govulncheck ./...

build:
	go tool goreleaser build --clean --single-target --snapshot

