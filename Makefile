.PHONY: test
.PHONY: fmt

run:
	go run ./main.go

test:
	go test ./test

fmt:
	gofmt -s -w .
