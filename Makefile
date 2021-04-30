lint:
	@./hack/check-license.sh
	@go fmt ./...

test:
	@ginkgo ./...

.PHONY: test lint
