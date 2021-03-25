lint:
	@go fmt ./...

test:
	@ginkgo ./...

.PHONY: test lint
