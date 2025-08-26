# Include env variables
include .env

# ------------------------------------------------------------------ #
#                               HELPERS                              #
# ------------------------------------------------------------------ #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'


# ------------------------------------------------------------------ #
#                             API SCRIPT                             #
# ------------------------------------------------------------------ #
## dev: run api in development mode
.PHONY: start/dev
start/dev:
	@air -c .air.toml


# ------------------------------------------------------------------ #
#                          Migration Script                          #
# ------------------------------------------------------------------ #
DB_DSN=postgres://gmoapi:${GMOAPI_DB_PASSWORD}@localhost:5433/gmoapi?sslmode=disable

## migrate/new: create new migration
.PHONY: migrate/new
migrate/new:
	@echo -n "Enter migration name: "; \
	read migration_name; \
	migration_name=$$(echo "$$migration_name" | tr ' ' '_'); \
	if [ -z "$$migration_name" ]; then \
		echo "\nError: Migration name cannot be empty." >&2; \
		exit 1; \
	fi; \
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DB_DSN} goose --dir ./migrations create $$migration_name sql

## migrate/up: apply all migration to latest
.PHONY: migrate/up
migrate/up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DB_DSN} goose --dir ./migrations up

## migrate/down: roll back migration by 1
.PHONY: migrate/down
migrate/down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DB_DSN} goose --dir ./migrations down

## migrate/reset: roll back all migration
.PHONY: migrate/reset
migrate/reset:
	@read -p "Are you sure you want to reset the DB? [y/N] " ans; \
	if echo "$$ans" | grep -iq '^y$$'; then \
		GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DB_DSN} goose --dir ./migrations reset; \
	fi

## migrate/version: show current version applied migration
.PHONY: migrate/version
migrate/version:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DB_DSN} goose --dir ./migrations version

## migrate/status: dump the migration status for the current DB
.PHONY: migrate/status
migrate/status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DB_DSN} goose --dir ./migrations status