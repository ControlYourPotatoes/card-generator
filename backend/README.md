# Card Generator Backend - Dependency Injection Implementation

This document describes the dependency injection foundation that has been implemented for the card generator backend.

## Overview

The dependency injection (DI) system provides a clean separation of concerns and makes the application more testable and maintainable. The implementation follows modern software engineering principles including SOLID principles and Clean Architecture.

## Architecture

### Core Components

1. **DI Container** (`pkg/di/container.go`)

   - Simple dependency injection container with reflection-based service resolution
   - Supports singleton and transient lifetimes
   - Thread-safe implementation with proper locking

2. **Configuration Management** (`pkg/config/config.go`)

   - Centralized configuration loading from multiple sources
   - Environment variables, YAML files, and defaults
   - Comprehensive validation for all configuration sections

3. **Bootstrap Module** (`pkg/bootstrap/bootstrap.go`)
   - Application initialization and dependency wiring
   - Service registration and configuration
   - Graceful shutdown handling

### Project Structure

```
backend/
├── pkg/                        # Shared packages
│   ├── di/                     # Dependency injection container
│   ├── config/                 # Configuration management
│   └── bootstrap/              # Application bootstrap
├── internal/
│   ├── core/
│   │   └── card/               # Core domain models
│   ├── storage/                # Storage implementations
│   │   ├── memory/             # In-memory storage
│   │   └── store.go            # Storage interface
│   ├── generator/              # Card generation
│   │   ├── art/                # Art processing
│   │   ├── text/               # Text rendering
│   │   ├── templates/          # Card templates
│   │   └── generator.go        # Main generator
│   └── parser/                 # CSV parsing
├── cmd/
│   ├── cardgen/                # Original main (legacy)
│   └── cardgen-di/             # New DI-enabled main
└── config/
    └── config.development.yaml # Development configuration
```

## Usage

### Running the Application

```bash
# Build the DI-enabled application
go build ./cmd/cardgen-di

# Run with default settings
./cardgen-di

# Run with custom configuration
./cardgen-di -env production -input cards.csv -output results.json
```

### Configuration

The application loads configuration from multiple sources in order of precedence:

1. Environment variables (highest priority)
2. YAML configuration files
3. Default values (lowest priority)

#### Environment Variables

Key environment variables:

- `APP_ENV`: Environment (development, production, test)
- `SERVER_PORT`: Server port
- `DB_TYPE`: Database type (memory, postgres, sqlite)
- `STORAGE_TYPE`: Storage type (memory, file, s3)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

#### YAML Configuration

Configuration files are loaded based on environment:

- `config/config.development.yaml`
- `config/config.production.yaml`
- `config/config.test.yaml`

### Dependency Registration

Services are registered in the bootstrap module:

```go
// Register singleton services
container.RegisterSingleton("cardStore", func() store.Store {
    return memory.New()
})

// Register transient services
container.RegisterTransient("csvParser", func(reader io.Reader) *parser.CSVParser {
    return parser.NewCSVParser(reader)
})
```

### Service Resolution

Services can be resolved from the application container:

```go
// Get card store
store, err := app.GetCardStore()

// Get card generator
generator, err := app.GetCardGenerator()

// Get CSV parser
parser, err := app.GetCSVParser(file)
```

## Benefits

1. **Loose Coupling**: Dependencies are injected rather than hard-coded
2. **Testability**: Easy to mock dependencies for unit testing
3. **Configuration**: Centralized and validated configuration management
4. **Maintainability**: Clear separation of concerns
5. **Extensibility**: Easy to add new services and implementations

## Future Enhancements

1. **Enhanced DI Container**: Support for dependency resolution in factory functions
2. **Service Discovery**: Automatic service registration via reflection
3. **Health Checks**: Built-in health check endpoints
4. **Metrics**: Prometheus metrics integration
5. **Tracing**: Distributed tracing support

## Migration from Legacy

The original `cmd/cardgen/main.go` demonstrates the old approach with manual dependency wiring. The new `cmd/cardgen-di/main.go` shows the clean DI approach. Both applications are functionally equivalent but the DI version is more maintainable and testable.

## Testing

```bash
# Run all tests
go test ./...

# Test specific packages
go test ./pkg/di
go test ./pkg/config
go test ./pkg/bootstrap
```
