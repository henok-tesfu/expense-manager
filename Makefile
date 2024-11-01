
ENVIRONMENT ?= development
# Check if we are in development or production
ifeq ($(ENVIRONMENT),development)
    # Include .env in development and fail if it doesn’t exist
    ifneq (,$(wildcard .env))
        include .env
        export
    else
        $(error ".env file is missing in development environment")
    endif
endif

# Variables for migration and database URL
MIGRATE=migrate
DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

# Ensure the migrations directory exists
migrations:
	@mkdir -p migrations

# Create a new migration (usage: make migrate-create)
migrate-create: migrations
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir migrations -seq $$name

# Run all migrations (up)
migrate-up: migrations
	$(MIGRATE) -database "$(DB_URL)" -path migrations up

# Rollback the last migration
migrate-down: migrations
	$(MIGRATE) -database "$(DB_URL)" -path migrations down 1

# Reset the database (drop all tables and re-run migrations)
migrate-reset: migrations
	$(MIGRATE) -database "$(DB_URL)" -path migrations drop -f && $(MIGRATE) -database "$(DB_URL)" -path migrations up

# Run the Go application
run:
	go run cmd/api/main.go
