.PHONY: build-importer run-importer clean clean-db test test-unit test-integration test-e2e test-coverage test-benchmark test-verbose help

# Build the importer tool
build-importer:
	@echo "Building card importer..."
	@go build -o bin/importer cmd/importer/main.go

# Build the DB cleaner tool
build-db-cleaner:
	@echo "Building database cleaner..."
	@go build -o bin/db-cleaner cmd/db-cleaner/main.go

# Build the card generator DI tool
build-cardgen-di:
	@echo "Building card generator with DI..."
	@cd backend && go build -o ../bin/cardgen-di ./cmd/cardgen-di

# Run the importer with a dry run for testing
test-importer: build-importer
	@echo "Running test import (dry run)..."
	@./bin/importer -type=anthem -file=backend/test/testdata/anthem_cards.csv -dry-run

# Run the importer with actual database import
run-importer: build-importer
	@echo "Running card import to database..."
	@./bin/importer -type=anthem -file=backend/test/testdata/anthem_cards.csv

# Clean database (remove all seed data)
clean-db: build-db-cleaner
	@echo "WARNING: This will delete ALL cards from the database!"
	@echo "Running database cleaner with -confirm flag..."
	@.\bin\db-cleaner -confirm

# Test targets

# Run all tests
test: test-unit test-integration
	@echo "All tests completed successfully!"

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@cd backend && go test ./pkg/... -v

# Run unit tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@cd backend && go test ./pkg/... -coverprofile=coverage.out
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

# Run tests with detailed output
test-verbose:
	@echo "Running tests with verbose output..."
	@cd backend && go test ./pkg/... -v -count=1

# Run benchmarks
test-benchmark:
	@echo "Running benchmarks..."
	@cd backend && go test ./pkg/... -bench=. -benchmem

# Run integration tests (Phase 2)
test-integration:
	@echo "Running integration tests..."
	@cd backend && go test ./test/integration/... -v

# Run end-to-end tests (Phase 3)
test-e2e:
	@echo "Running end-to-end tests..."
	@cd backend && go test ./test/e2e/... -v

# Test specific package
test-pkg:
	@if [ -z "$(PKG)" ]; then \
		echo "Usage: make test-pkg PKG=package_name"; \
		echo "Example: make test-pkg PKG=di"; \
		exit 1; \
	fi
	@echo "Testing package: $(PKG)"
	@cd backend && go test ./pkg/$(PKG) -v

# Watch tests (requires entr or similar tool)
test-watch:
	@echo "Watching for changes and running tests..."
	@find backend -name "*.go" | entr -c make test-unit

# Clean test artifacts
clean-test:
	@echo "Cleaning test artifacts..."
	@cd backend && rm -f coverage.out coverage.html
	@cd backend && go clean -testcache

# Format and lint code
fmt:
	@echo "Formatting code..."
	@cd backend && go fmt ./pkg/...

# Vet code for issues
vet:
	@echo "Vetting code..."
	@cd backend && go vet ./pkg/...

# Run all quality checks
quality: fmt vet test-unit
	@echo "All quality checks passed!"

# Development workflow
dev: quality test-coverage
	@echo "Development checks completed!"

# CI workflow
ci: quality test-coverage test-benchmark
	@echo "CI pipeline completed!"

# Format all code (including potentially problematic files)
fmt-all:
	@echo "Formatting all code..."
	@cd backend && go fmt ./...

# Vet all code (including potentially problematic files)
vet-all:
	@echo "Vetting all code..."
	@cd backend && go vet ./...

# Generate test data
generate-testdata:
	@echo "Generating test data..."
	@echo "Copying test data to backend/test/testdata..."
	@cp test/data/*.csv backend/test/testdata/ 2>/dev/null || true

# Create directories needed for the project
init:
	@mkdir -p bin
	@mkdir -p backend/test/testdata
	@mkdir -p backend/test/integration
	@mkdir -p backend/test/e2e
	@mkdir -p cmd/db-cleaner

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@cd backend && go clean

# Help command
help:
	@echo "Available commands:"
	@echo "  make build-importer     - Build the card importer tool"
	@echo "  make build-cardgen-di   - Build the card generator with DI"
	@echo "  make test-importer      - Test the importer with a dry run (no database writes)"
	@echo "  make run-importer       - Run the importer and save cards to the database"
	@echo "  make clean              - Remove build artifacts"
	@echo "  make clean-db           - Clean all card data from the database"
	@echo "  make clean-test         - Clean test artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  make test               - Run all tests (unit + integration)"
	@echo "  make test-unit          - Run unit tests only"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make test-e2e           - Run end-to-end tests"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make test-verbose       - Run tests with detailed output"
	@echo "  make test-benchmark     - Run performance benchmarks"
	@echo "  make test-pkg PKG=name  - Test specific package"
	@echo "  make test-watch         - Watch for changes and run tests"
	@echo ""
	@echo "Quality:"
	@echo "  make fmt                - Format code"
	@echo "  make vet                - Vet code for issues"
	@echo "  make quality            - Run format, vet, and unit tests"
	@echo "  make dev                - Development workflow (quality + coverage)"
	@echo "  make ci                 - CI workflow (quality + coverage + benchmarks)"
	@echo ""
	@echo "  make generate-testdata  - Copy test data to backend/test/testdata"
	@echo "  make init               - Create necessary directories"
	@echo "  make help               - Show this help information"