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

check: static vuln deadcode gosec helm-lint goimport

static:
	@echo "Running staticcheck..."
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

vuln:
	@echo "Running vulncheck..."
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

deadcode:
	@echo "Running deadcode..."
	go run golang.org/x/tools/cmd/deadcode@latest -filter "pkg/client/client" -test ./...

gosec:
	@echo "Running gosec..."
	go run github.com/securego/gosec/v2/cmd/gosec@latest --exclude G404,G101 --exclude-generated -terse ./...

helm-lint:
	@echo "Running helm lint..."
	helm lint --strict ./charts

goimport:
	@echo "Running goimport..."
	find . -name '*.go' -exec go run golang.org/x/tools/cmd/goimports@latest -l -w {} +