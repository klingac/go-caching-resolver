GO=GO111MODULE=on go

.PHONY: test lint build format

all: format lint test

format:
	gofmt -s -l -w internal/ cmd/ *.go

build:
	$(GO) build 

test:
	$(GO) test ./...

lint:
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:v1.46.2 golangci-lint run --color always ${ARGS}

