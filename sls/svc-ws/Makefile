SHELL := /bin/bash

build:
	@echo "Building all handlers:"
	@for fn in $(shell ls handlers); do \
		echo -n "	- building $${fn}..."; \
		env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/$${fn} handlers/$${fn}/$${fn}.go; \
		echo "	done"; \
	done

test:
	go test ./handlers/...

clean:
	rm -rf ./bin

destroy-stack:
	sls remove

deploy: clean build test
	sls deploy --verbose

deploy-%: clean
	go test ./handlers/$*/...
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/$* handlers/$*/$*.go
	sls deploy -f $*

log-%:
	sls logs -f $*

tail-%:
	sls logs -f $* -t
