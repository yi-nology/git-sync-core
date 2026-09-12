.PHONY: test tidy

tidy:
	@go mod tidy

test:
	@go test ./... -race -count=1
