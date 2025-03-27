.PHONY: build-importer run-importer clean clean-db

# Build the importer tool
build-importer:
	@echo "Building card importer..."
	@go build -o bin/importer cmd/importer/main.go

# Build the DB cleaner tool
build-db-cleaner:
	@echo "Building database cleaner..."
	@go build -o bin/db-cleaner cmd/db-cleaner/main.go

# Run the importer with a dry run for testing
test-importer: build-importer
	@echo "Running test import (dry run)..."
	@./bin/importer -type=anthem -file=test/data/anthem_cards.csv -dry-run

# Run the importer with actual database import
run-importer: build-importer
	@echo "Running card import to database..."
	@./bin/importer -type=anthem -file=test/data/anthem_cards.csv

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

# Clean database (remove all seed data)
clean-db: build-db-cleaner
	@echo "WARNING: This will delete ALL cards from the database!"
	@echo "Are you sure? [y/N] " && read ans && [ ${ans:-N} = y ]
	@echo "Cleaning database..."
	@./bin/db-cleaner -confirm

# Create directories needed for the project
init:
	@mkdir -p bin
	@mkdir -p test/data
	@mkdir -p cmd/db-cleaner

# Help command
help:
	@echo "Available commands:"
	@echo "  make build-importer   - Build the card importer tool"
	@echo "  make test-importer    - Test the importer with a dry run (no database writes)"
	@echo "  make run-importer     - Run the importer and save cards to the database"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make clean-db         - Clean all card data from the database"
	@echo "  make init             - Create necessary directories"
	@echo "  make help             - Show this help information"