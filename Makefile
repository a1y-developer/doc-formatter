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

.PHONY: help test cover cover-html lint lint-fix atlas migration migration-status migrate migrate-down migrate-hash
