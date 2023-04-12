.PHONY: bootstrap
bootstrap:
	go build -o bin/bootstrap cmd/bootstrap/*.go

test: fmt vet
	go test ./... -coverprofile cover.out -short
fmt:
	go fmt ./...
vet:
	go vet ./...

local:
	go run cmd/bootstrap/main.go