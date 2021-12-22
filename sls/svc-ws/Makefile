SHELL := /bin/bash
.PHONY: clean build test deploy destroy deploy-% log-% tail-%

clean:
	rm -rf ./bin

build:
	@echo "Building all handlers:"
	@for fn in $(shell ls handlers); do \
		echo -n "	- building $${fn}..."; \
		env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/$${fn} handlers/$${fn}/$${fn}.go; \
		echo "	done"; \
	done

test:
	go test ./handlers/...

deploy: clean build test
	sls deploy --verbose

destroy:
	sls remove

deploy-%: clean
	go test ./handlers/$*/...
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/$* handlers/$*/$*.go
	sls deploy -f $*

log-%:
	sls logs -f $*

tail-%:
	sls logs -f $* -t
