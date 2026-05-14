ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

BASE_STACK = docker compose -f docker-compose.yml

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

test-integration: ### Run integration tests
	go test -v ./integration_test/...
.PHONY: test-integration

test-usecase: ### Run usecase tests
	go test -v ./internal/usecase/...
.PHONY: test-usecase

test: ### run test
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/... ./pkg/...
.PHONY: test

compose-up-db: ### Run docker compose db container in background
	$(BASE_STACK) up -d db
	@echo "DB running on localhost:5433"
.PHONY: compose-up-db

compose-stop-db: ### Stop the db container (keeps it for fast restart)
	$(BASE_STACK) stop db
.PHONY: compose-stop-db

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations '$(word 2,$(MAKECMDGOALS))'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(MYSQL_URL)' up
.PHONY: migrate-up

migrate-down: ### rollback the last migration
	migrate -path migrations -database '$(MYSQL_URL)' down 1
.PHONY: migrate-down

mock: ### run mockgen
	go tool mockgen -source ./internal/repo/interfaces.go -package usecase_test > ./internal/usecase/mocks_repo_test.go
	go tool mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock

lint: ### Run golangci-lint
	golangci-lint run
.PHONY: lint

format: ### Run code formatter
	go fix ./...
	go tool gofumpt -l -w .
	go tool gci write . --skip-generated -s standard -s default
.PHONY: format

swag-v1: ### swag init
	go tool swag init --parseDependency -g internal/controller/restapi/router.go
.PHONY: swag-v1

deps: ### deps tidy + verify
	go mod tidy && go mod verify
.PHONY: deps

vul-check: ### run vulnerable check
	go tool govulncheck ./...
.PHONY: vul-check

pre-commit: deps swag-v1 mock format lint test ### run pre-commit
.PHONY: pre-commit

