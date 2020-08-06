.PHONY: test build/crud

ginkgo := go run github.com/onsi/ginkgo/ginkgo --race --randomizeAllSpecs

test: build/crud lint
	$(ginkgo) acceptance

build/crud:
	go build -o build/crud .

lint:
	golangci-lint run