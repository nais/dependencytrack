.PHONY: bootstrap
bootstrap:
	go build -o bin/bootstrap cmd/bootstrap/*.go

integration_test: fmt vet
	go test ./... -tags integration_test -run TestIntegration
test: fmt vet
	go test ./... -coverprofile cover.out -short
fmt:
	go run mvdan.cc/gofumpt -w ./
vet:
	go vet ./...

local:
	go run cmd/bootstrap/main.go

compose:
	docker-compose build && docker-compose up