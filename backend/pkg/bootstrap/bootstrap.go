package bootstrap

import (
	"fmt"
	"io"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/parser"
	store "github.com/ControlYourPotatoes/card-generator/backend/internal/storage"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/storage/memory"
	"github.com/ControlYourPotatoes/card-generator/backend/pkg/config"
	"github.com/ControlYourPotatoes/card-generator/backend/pkg/di"
)

// Application holds the main application dependencies
type Application struct {
	Config    *config.Config
	Container di.Container
}

// NewApplication creates and configures the application with all dependencies
func NewApplication(env string) (*Application, error) {
	// Load configuration
	cfg, err := config.LoadConfig(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create DI container
	container := di.NewContainer()

	// Register configuration as singleton
	if err := container.RegisterSingleton("config", func() *config.Config {
		return cfg
	}); err != nil {
		return nil, fmt.Errorf("failed to register config: %w", err)
	}

	// Register storage based on configuration
	if err := registerStorage(container, cfg); err != nil {
		return nil, fmt.Errorf("failed to register storage: %w", err)
	}

	// Register card generator
	if err := container.RegisterSingleton("cardGenerator", func() (generator.CardGenerator, error) {
		return generator.NewCardGenerator()
	}); err != nil {
		return nil, fmt.Errorf("failed to register card generator: %w", err)
	}

	// Register CSV parser factory (transient since parsers are typically per-operation)
	if err := container.RegisterTransient("csvParserFactory", createCSVParserFactory); err != nil {
		return nil, fmt.Errorf("failed to register CSV parser factory: %w", err)
	}

	return &Application{
		Config:    cfg,
		Container: container,
	}, nil
}

// registerStorage registers the appropriate storage implementation based on configuration
func registerStorage(container di.Container, cfg *config.Config) error {
	switch cfg.Storage.Type {
	case "memory":
		return container.RegisterSingleton("cardStore", func() store.Store {
			return memory.New()
		})
	case "file":
		// TODO: Implement file storage
		return container.RegisterSingleton("cardStore", func() store.Store {
			return memory.New() // Fallback to memory for now
		})
	case "database":
		// TODO: Implement database storage
		return container.RegisterSingleton("cardStore", func() store.Store {
			return memory.New() // Fallback to memory for now
		})
	default:
		return fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}
}

// createCSVParserFactory creates a CSV parser factory function
func createCSVParserFactory() func(reader io.Reader) *parser.CSVParser {
	return func(reader io.Reader) *parser.CSVParser {
		return parser.NewCSVParser(reader)
	}
}

// GetCardStore resolves the card store from the container
func (app *Application) GetCardStore() (store.Store, error) {
	instance, err := app.Container.Resolve("cardStore")
	if err != nil {
		return nil, err
	}
	return instance.(store.Store), nil
}

// GetCardGenerator resolves the card generator from the container
func (app *Application) GetCardGenerator() (generator.CardGenerator, error) {
	instance, err := app.Container.Resolve("cardGenerator")
	if err != nil {
		return nil, err
	}
	return instance.(generator.CardGenerator), nil
}

// GetCSVParser creates a new CSV parser instance
func (app *Application) GetCSVParser(reader io.Reader) (*parser.CSVParser, error) {
	factoryInstance, err := app.Container.Resolve("csvParserFactory")
	if err != nil {
		return nil, err
	}
	
	factory := factoryInstance.(func(io.Reader) *parser.CSVParser)
	return factory(reader), nil
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	// Get card store and close if it implements a Close method
	if store, err := app.GetCardStore(); err == nil {
		if err := store.Close(); err != nil {
			return fmt.Errorf("failed to close card store: %w", err)
		}
	}

	// Clean the container
	if containerImpl, ok := app.Container.(interface{ Clear() }); ok {
		containerImpl.Clear()
	}

	return nil
} 