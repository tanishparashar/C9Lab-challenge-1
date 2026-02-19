# Database Migration Makefile
# Requires golang-migrate to be installed

# Database connection string
DB_URL=postgres://devops:devops123@localhost:5432/devops_db?sslmode=disable
MIGRATIONS_PATH=./migrations

# Commands
.PHONY: help migrate-up migrate-down migrate-create migrate-force migrate-version migrate-drop

help:
	@echo "Available commands:"
	@echo "  make migrate-up       - Apply all pending migrations"
	@echo "  make migrate-down     - Rollback the last migration"
	@echo "  make migrate-create NAME=<name> - Create a new migration file"
	@echo "  make migrate-version  - Show current migration version"
	@echo "  make migrate-force VERSION=<n> - Force set migration version"
	@echo "  make migrate-drop     - Drop everything (use with caution!)"

migrate-up:
	@echo "Applying migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

migrate-version:
	@echo "Getting current migration version..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make migrate-force VERSION=<version_number>"; \
		exit 1; \
	fi
	@echo "Forcing migration version to $(VERSION)..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)

migrate-drop:
	@echo "WARNING: This will drop all tables!"
	@read -p "Are you sure? (yes/no): " confirm && [ "$$confirm" = "yes" ] || exit 1
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop
