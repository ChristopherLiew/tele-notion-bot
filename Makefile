.PHONY: test build run all

fmt:
	gofmt -s -w .

test:
	go test -v ./test

test_coverage:
	go test ./test -coverprofile=./test/coverage.out
	go tool cover -html ./test/coverage.out -o ./test/coverage.html
	open ./test/coverage.html

build:
	cd ./cmd && go build -o ../bin/main && ../bin/main

run:
	go run ./cmd/main.go

dep:
	go mod download

vet:
	go vet ./...

clean:
	go mod tidy -v
	rm -rf bin
	rm -rf ./test/coverage.*

all:
	make clean
	make fmt
	make vet
	make test
	make build
