.PHONY: all build run gotool clean help
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o comic
detect:
	go run main.go detect
run:
	go run main.go server
gotool:
	go fmt ./...
	go vet ./
tidy:
	go mod tidy
