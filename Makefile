.PHONY: test build/crud

ginkgo := go run github.com/onsi/ginkgo/ginkgo --race --randomizeAllSpecs

test: build/crud
	$(ginkgo) acceptance

build/crud:
	go build -o build/crud .

