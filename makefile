# Directory where migration files are located
MIGRATION_DIR := database/mysql/migration
IS_IN_PROGRESS = "is in progress ..."

# Get list of domains, compatible with Unix-like systems
OS := $(shell uname 2>/dev/null || echo Unknown)
ifeq ($(OS),Unknown)
    DOMAINS := $(shell for /d %D in (internal/app/*) do @echo %~nxD)  # Windows cmd (not perfect)
else
    DOMAINS := $(shell find internal/app -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)
endif

.PHONY: all
all: env install mod

## help: prints this help message
.PHONY: help
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## env: will setup env
.PHONY: env
env:
	@echo "make env ${IS_IN_PROGRESS}"
	@go env -w GO111MODULE=on
	@go env -w GOBIN=`go env GOPATH`/bin
	@go env -w GOPROXY=https://proxy.golang.org,direct

## mod: will pull all dependency
.PHONY: mod
mod:
	@echo "make mod ${IS_IN_PROGRESS}"
	@rm -rf ./vendor ./go.sum
	@go mod tidy
	@go mod vendor

## setup: Set up database temporary for local environtment
.PHONY: setup
setup:
	@echo "make setup ${IS_IN_PROGRESS}"
	@docker-compose -f ./infrastructure/docker-compose.local.yml up -d
	@sleep 8

## down: Set down database temporary for integration testing
.PHONY: down
down: 
	@echo "make down ${IS_IN_PROGRESS}"
	@docker-compose -f ./infrastructure/docker-compose.local.yml down -t 1

## run: run for running app on local
.PHONY: run
run:
	@go run main.go


.PHONY: migration-up
migration-up:
	GOOSE_DRIVER=mysql GOOSE_DBSTRING="root:pwd@tcp(localhost:3307)/golang_clean_architecture?parseTime=true" goose -dir=$(MIGRATION_DIR) up

.PHONY: migration-down
migration-down: 
	GOOSE_DRIVER=mysql GOOSE_DBSTRING="root:pwd@tcp(localhost:3307)/golang_clean_architecture?parseTime=true" goose -dir=$(MIGRATION_DIR) down

## generate-mocks: will generate mock for internal/app, internal/integrations, and pkg/integrations
.PHONY: generate-mocks
generate-mocks:
	@mkdir mocks 2>/dev/null || true
	@domains="$(if $(DOMAIN),$(DOMAIN),$(DOMAINS))"; \
	for domain in $$domains; do \
		for file in repository usecase; do \
			if [ -f internal/app/$$domain/$$file.go ]; then \
				echo "Generating mock-app for $$file: $$domain"; \
				mockgen -source internal/app/$$domain/$$file.go -destination mocks/$$domain/$${file}_mock.go -package=mocks; \
			fi; \
		done; \
	done


## unit-test: will test with unit tags
.PHONY: unit-test
unit-test:
	@echo "make unit-test ${IS_IN_PROGRESS}"
	@go clean -testcache
	@go test \
		--race -count=1 -cpu=1 -parallel=1 -timeout=90s -failfast -vet= \
		-cover -covermode=atomic -coverprofile=./.coverage/unit.out \
		./internal/users/usecase/...

## integration-test: will test with integration tags
.PHONY: integration-test
integration-test:
	@echo "make integration-test ${IS_IN_PROGRESS}"
	@go clean -testcache
	@go test --race -timeout=90s -failfast \
		-vet= -cover -covermode=atomic -coverprofile=./.coverage/integration.out -tags=integration \
		./internal/users/repository/...

## run-integration-test: will run integration test with any dependencies
.PHONY: run-integration-test
run-integration-test:setup migration-up integration-test migration-down down

## tests: run tests(integration, unit & e2e testing) and any dependencies
.PHONY: tests
tests:run-integration-test unit-test

## cover: will report all test coverage
.PHONY: cover
cover:
	@make -s cover-with type=integration
	@make -s cover-with type=unit



## swagger-gen: swagger-gen for generate swagger documentation
.PHONY: swagger-gen
swagger-gen:
	@swag init -md ./docs/ && ./docs/fix.sh