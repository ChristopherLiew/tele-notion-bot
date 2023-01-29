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
	go build -o ./bin/main && ./bin/main

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

run:
	go run ./main.go

dep:
	go mod download

vet:
	go vet

clean:
	rm -rf bin
	rm -rf ./test/coverage.*

all:
	make clean
	make fmt
	make vet
	make test
	make build
