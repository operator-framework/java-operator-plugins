lint:
	@./hack/check-license.sh
	@go fmt ./...

test:
	@go test -coverprofile=coverage.out -covermode=count -short ./...

.PHONY: test lint
