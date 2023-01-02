.PHONY: test
.PHONY: fmt

run:
	go run cmd/travel-buddy-bot/main.go

test:
	go test ./test

fmt:
	gofmt -s -w .
