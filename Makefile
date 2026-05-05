# Delegates to mise tasks. Use `mise run <task>` directly, or the aliases below.
# See mise/tasks/ for full task definitions.

.PHONY: test check build generate bootstrap tidy local compose connect-db integration_test

test:
	mise run test

check:
	mise run check

build:
	mise run build

generate:
	mise run generate

bootstrap:
	mise run build:bootstrap

tidy:
	@echo "Running go mod tidy for all modules..."
	find . -name go.mod -execdir go mod tidy \;

local:
	go run cmd/bootstrap/main.go

compose:
	docker-compose build && docker-compose up

integration_test:
	mise run test-integration

# make connect-db I=dependencytrack-61ed1cda P=nais-management-4203 S=dependencytrack
# psql -U postgres -h localhost dependencytrack
connect-db:
	@CONNECTION_NAME=$$(gcloud sql instances describe $(I) \
	  --format="get(connectionName)" \
	  --project $(P)) && \
	cloud-sql-proxy $$CONNECTION_NAME \
	    --auto-iam-authn \
	    --impersonate-service-account="$(S)@$(P).iam.gserviceaccount.com"
