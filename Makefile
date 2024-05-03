.PHONY: default ci lint test yaegi_test vendor clean generate tidy spell inst vulncheck

export GO111MODULE=on

default: spell lint test

ci: inst tidy generate default vulncheck

lint:
	golangci-lint run

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
