IGNORED_PATH := "pkg/dependencytrack/client"
GO_PACKAGES := $(shell go list ./... | grep -v $(IGNORED_PATH))

.PHONY: bootstrap
bootstrap:
	go build -o bin/bootstrap cmd/bootstrap/*.go

integration_test:
	go test -v -count=1 ./... -tags integration_test -run TestIntegration
test:
	go test ./... -coverprofile cover.out -short

local:
	go run cmd/bootstrap/main.go

compose:
	docker-compose build && docker-compose up

check: vet staticcheck vuln deadcode helm-lint goimport # gosec

vet:
	go vet ./...

goimport:
	@echo "Running goimport..."
	find . -type f -name '*.go' ! -path './pkg/dependencytrack/client/*' -exec go run golang.org/x/tools/cmd/goimports@latest -l -w  {} +

fmt:
	@echo "Running go fmt..."
	go fmt $(GO_PACKAGES)

staticcheck:
	@echo "Running staticcheck..."
	go run honnef.co/go/tools/cmd/staticcheck@latest -f=stylish  $(GO_PACKAGES)

vuln:
	@echo "Running vulncheck..."
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./... || true
	@echo "Govulncheck finished – vulnerabilities reported"

deadcode:
	@echo "Running deadcode..."
	go run golang.org/x/tools/cmd/deadcode@latest -filter "pkg/client/client" -test ./...

gosec:
	@echo "Running gosec..."
	go run github.com/securego/gosec/v2/cmd/gosec@latest --exclude-generated -terse ./...

helm-lint:
	@echo "Running helm lint..."
	helm lint --strict ./charts

# make connect-db I=dependencytrack-61ed1cda P=nais-management-4203 S=dependencytrack
# psql -U postgres -h localhost dependencytrack
connect-db:
	@CONNECTION_NAME=$$(gcloud sql instances describe $(I) \
	  --format="get(connectionName)" \
	  --project $(P)) && \
	cloud-sql-proxy $$CONNECTION_NAME \
	    --auto-iam-authn \
	    --impersonate-service-account="$(S)@$(P).iam.gserviceaccount.com"

generate: generate-client generate-mocks

generate-client:
	@echo "Generating Go code from the OpenAPI specification..."
	@openapi-generator generate \
        -i schema/dtrack.json \
        -g go \
        -o pkg/dependencytrack/client \
        --global-property apiTests=false,modelTests=false \
        --package-name client \
        --additional-properties=withGoMod=false \
        --additional-properties=generateInterfaces=true \
        --additional-properties=packageName=client || { \
			echo "Error: openapi-generator is not installed or failed to execute."; \
			echo "Please visit https://openapi-generator.tech/docs/installation/ for installation instructions."; \
			exit 1; \
		}

generate-mocks:
	#find pkg -type f -name "mock_*.go" -delete
	go run github.com/vektra/mockery/v3 --config ./mockery.yaml
	#find pkg -type f -name "mock_*.go" -exec go run mvdan.cc/gofumpt@latest -w {} \;
