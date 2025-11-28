SHELL = /bin/bash
PWD:=$(shell pwd)

GOLINTER			?= golangci-lint
GOLINTER_VERSION	?= v2.6.2
COVER_FILE			?= coverage.out
SOURCE_PATHS		?= ./...
BUILD_PATH 			?= ./_build/bin
ATLAS				?= atlas


help:  ## This help message :)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test:  ## Run the tests
	go test -gcflags="all=-l -N" -timeout=10m `go list $(SOURCE_PATHS) | grep -v /api` ${TEST_FLAGS}

cover:  ## Generates coverage report
	go test -gcflags="all=-l -N" -timeout=10m `go list $(SOURCE_PATHS) | grep -v /api` -coverprofile $(COVER_FILE) ${TEST_FLAGS}

cover-html:  ## Generates coverage report and displays it in the browser
	go tool cover -html=$(COVER_FILE)

lint:  ## Lint, will not fix but sets exit code on error
	@which $(GOLINTER) > /dev/null || (echo "Installing $(GOLINTER)@$(GOLINTER_VERSION) ..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLINTER_VERSION) && echo -e "Installation complete!\n")
	$(GOLINTER) run $(SOURCE_PATHS)

lint-fix:  ## Lint, will try to fix errors and modify code
	@which $(GOLINTER) > /dev/null || (echo "Installing $(GOLINTER)@$(GOLINTER_VERSION) ..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLINTER_VERSION) && echo -e "Installation complete!\n")
	$(GOLINTER) run --fix $(SOURCE_PATHS)

doc:  ## Start the documentation server with godoc
	@which godoc > /dev/null || (echo "Installing godoc@latest ..."; go install golang.org/x/tools/cmd/godoc@latest && echo -e "Installation complete!\n")
	godoc -http=:6060

atlas: ## Install Atlas CLI
	@which $(ATLAS) > /dev/null || (echo "Installing $(ATLAS) ..."; curl -sSf https://atlasgo.sh | sh && echo -e "Installation complete!\n")

migration-status: atlas ## Show current migration status for a service
	@{ \
		read -p "Enter service name (e.g. auth): " SERVICE_NAME; \
		$(ATLAS) migrate status --env "$${SERVICE_NAME}"; \
	}

migration: atlas ## Generate new database migration file for service
	@read -p "Enter service name (e.g. auth): " SERVICE_NAME; \
	$(ATLAS) migrate diff --env $${SERVICE_NAME}; \
	$(ATLAS) migrate status --env $${SERVICE_NAME}

migrate-dry-run: atlas ## Show pending migrations (dry-run)
	@{ \
		read -p "Enter service name (e.g. auth): " SERVICE_NAME; \
		$(ATLAS) migrate apply --env "$${SERVICE_NAME}" --dry-run; \
	}


migrate: atlas ## Apply pending migrations for a service
	@{ \
		read -p "Enter service name (e.g. auth): " SERVICE_NAME; \
		$(ATLAS) migrate apply --env "$${SERVICE_NAME}" --dry-run; \
		echo "--- Above is the dry-run preview ---"; \
		read -p "Apply these migrations? (y/n): " CONFIRM; \
		if [ "$$CONFIRM" != "y" ]; then \
			echo "Migration cancelled."; \
			exit 0; \
		fi; \
		$(ATLAS) migrate apply --env "$${SERVICE_NAME}"; \
		$(ATLAS) migrate status --env "$${SERVICE_NAME}"; \
	}

migrate-down: atlas ## Migrate service database down to specific version
	@{ \
		read -p "Enter service name (e.g. auth): " SERVICE_NAME; \
		read -p "Enter version to migrate DOWN to (e.g. 20230915120000): " VERSION; \
		echo "You are about to migrate DOWN '$${SERVICE_NAME}' to version '$${VERSION}'"; \
		read -p "Migrate down? (y/n): " CONFIRM; \
		if [ "$$CONFIRM" != "y" ]; then \
			echo "Migration cancelled."; \
			exit 0; \
		fi; \
		$(ATLAS) migrate down --env "$${SERVICE_NAME}" --to-version "$${VERSION}"; \
		$(ATLAS) migrate status --env "$${SERVICE_NAME}"; \
	}

migrate-hash: atlas ## Re-hash service database migration files
	@{ \
		read -p "Enter service name (e.g. auth): " SERVICE_NAME; \
		$(ATLAS) migrate hash --env "$${SERVICE_NAME}"; \
		$(ATLAS) migrate status --env "$${SERVICE_NAME}"; \
	}

migrate-validate: atlas ## Validate service database migration files
	@{ \
		read -p "Enter service name (e.g. auth): " SERVICE_NAME; \
		$(ATLAS) migrate validate --env "$${SERVICE_NAME}"; \
	}

gen-api-spec: ## Generate API Specification with OpenAPI format
	@echo "Checking swag installation..."
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "Installing swag@v1.16.4 ..."; \
		go install github.com/swaggo/swag/cmd/swag@v1.16.4 || { echo 'Failed to install swag'; exit 1; }; \
		echo "swag installed successfully!"; \
	else \
		echo "swag is already installed"; \
	fi
	@echo "Running swag init..."
	@swag init --parseInternal -g cmd/gateway/main.go -o api/http/gateway/v1/ || { echo 'swag init failed!'; exit 1; }

	@echo "Running swag fmt..."
	@swag fmt --dir internal/gateway/ || { echo 'swag fmt failed!'; exit 1; }
	@echo "API Spec generated successfully without errors!"

gen-api-doc: ## Generate API Documentation by API Specification
	@echo "Checking swagger installation..."
	@if ! command -v swagger >/dev/null 2>&1; then \
		echo "Installing swagger@v0.33.1 ..."; \
		go install github.com/go-swagger/go-swagger/cmd/swagger@v0.33.1 || { echo 'Failed to install swagger'; exit 1; }; \
		echo "swagger installed successfully!"; \
	else \
		echo "swagger is already installed"; \
	fi
	@echo "Generating API Markdown Documentation..."
	@swagger generate markdown -f ./api/http/gateway/v1/swagger.json --output=docs/api.md || { echo 'swagger generate markdown failed!'; exit 1; }
	@echo "API documentation generated successfully!"

.PHONY: help test cover cover-html lint lint-fix doc atlas migration migration-status migrate migrate-down migrate-hash migrate-validate gen-api-spec gen-api-doc
