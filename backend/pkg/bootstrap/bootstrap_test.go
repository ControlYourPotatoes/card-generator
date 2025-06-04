package bootstrap

import (
	"io"
	"strings"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator"
	store "github.com/ControlYourPotatoes/card-generator/backend/internal/storage"
	"github.com/ControlYourPotatoes/card-generator/backend/pkg/config"
	"github.com/ControlYourPotatoes/card-generator/backend/pkg/di"
)

func TestNewApplication_Success(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	if app == nil {
		t.Fatal("Application is nil")
	}

	if app.Config == nil {
		t.Fatal("Application config is nil")
	}

	if app.Container == nil {
		t.Fatal("Application container is nil")
	}

	// Verify container has expected services
	services := app.Container.GetRegisteredServices()
	expectedServices := []string{"config", "cardStore", "cardGenerator", "csvParserFactory"}

	if len(services) != len(expectedServices) {
		t.Errorf("Expected %d services, got %d", len(expectedServices), len(services))
	}

	// Check each expected service exists
	for _, expected := range expectedServices {
		found := false
		for _, service := range services {
			if service == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected service '%s' not found in registered services", expected)
		}
	}
}

func TestNewApplication_ConfigResolution(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	// Resolve config from container
	configInstance, err := app.Container.Resolve("config")
	if err != nil {
		t.Fatalf("Failed to resolve config: %v", err)
	}

	config, ok := configInstance.(*config.Config)
	if !ok {
		t.Fatal("Resolved config is not of type *config.Config")
	}

	// Verify it's the same config as in the application
	if config != app.Config {
		t.Error("Resolved config is not the same as application config")
	}
}

func TestGetCardStore_Success(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	storeInstance, err := app.GetCardStore()
	if err != nil {
		t.Fatalf("Failed to get card store: %v", err)
	}

	if storeInstance == nil {
		t.Fatal("Card store is nil")
	}

	// Verify it implements the Store interface
	_, ok := storeInstance.(store.Store)
	if !ok {
		t.Fatal("Returned object does not implement Store interface")
	}
}

func TestGetCardGenerator_Success(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	gen, err := app.GetCardGenerator()
	if err != nil {
		t.Fatalf("Failed to get card generator: %v", err)
	}

	if gen == nil {
		t.Fatal("Card generator is nil")
	}

	// Verify it implements the CardGenerator interface
	_, ok := gen.(generator.CardGenerator)
	if !ok {
		t.Fatal("Returned object does not implement CardGenerator interface")
	}
}

func TestGetCSVParser_Success(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	// Create a test reader
	reader := strings.NewReader("test,data\n1,2")

	parser, err := app.GetCSVParser(reader)
	if err != nil {
		t.Fatalf("Failed to get CSV parser: %v", err)
	}

	if parser == nil {
		t.Fatal("CSV parser is nil")
	}

	// Verify it's not nil (type is already guaranteed by the function signature)
	// No need for type assertion since GetCSVParser returns *parser.CSVParser
}

func TestGetCSVParser_MultipleInstances(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	// Create two parsers
	reader1 := strings.NewReader("test1")
	reader2 := strings.NewReader("test2")

	parser1, err := app.GetCSVParser(reader1)
	if err != nil {
		t.Fatalf("Failed to get first CSV parser: %v", err)
	}

	parser2, err := app.GetCSVParser(reader2)
	if err != nil {
		t.Fatalf("Failed to get second CSV parser: %v", err)
	}

	// Verify they are different instances (transient registration)
	if parser1 == parser2 {
		t.Error("CSV parsers should be different instances (transient)")
	}
}

func TestSingletonBehavior(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	// Get card store twice
	store1, err := app.GetCardStore()
	if err != nil {
		t.Fatalf("Failed to get first card store: %v", err)
	}

	store2, err := app.GetCardStore()
	if err != nil {
		t.Fatalf("Failed to get second card store: %v", err)
	}

	// Verify they are the same instance (singleton registration)
	if store1 != store2 {
		t.Error("Card stores should be the same instance (singleton)")
	}

	// Same test for card generator
	gen1, err := app.GetCardGenerator()
	if err != nil {
		t.Fatalf("Failed to get first card generator: %v", err)
	}

	gen2, err := app.GetCardGenerator()
	if err != nil {
		t.Fatalf("Failed to get second card generator: %v", err)
	}

	if gen1 != gen2 {
		t.Error("Card generators should be the same instance (singleton)")
	}
}

func TestRegisterStorage_MemoryType(t *testing.T) {
	container := di.NewContainer()
	cfg := &config.Config{
		Storage: config.StorageConfig{
			Type: "memory",
		},
	}

	err := registerStorage(container, cfg)
	if err != nil {
		t.Fatalf("Failed to register memory storage: %v", err)
	}

	// Verify storage is registered
	services := container.GetRegisteredServices()
	found := false
	for _, service := range services {
		if service == "cardStore" {
			found = true
			break
		}
	}

	if !found {
		t.Error("cardStore service not registered")
	}
}

func TestRegisterStorage_FileType(t *testing.T) {
	container := di.NewContainer()
	cfg := &config.Config{
		Storage: config.StorageConfig{
			Type: "file",
		},
	}

	err := registerStorage(container, cfg)
	if err != nil {
		t.Fatalf("Failed to register file storage: %v", err)
	}

	// For now, file type falls back to memory storage
	// This test verifies the registration succeeds
	services := container.GetRegisteredServices()
	found := false
	for _, service := range services {
		if service == "cardStore" {
			found = true
			break
		}
	}

	if !found {
		t.Error("cardStore service not registered")
	}
}

func TestRegisterStorage_UnsupportedType(t *testing.T) {
	container := di.NewContainer()
	cfg := &config.Config{
		Storage: config.StorageConfig{
			Type: "unsupported",
		},
	}

	err := registerStorage(container, cfg)
	if err == nil {
		t.Fatal("Expected error for unsupported storage type")
	}

	expectedError := "unsupported storage type: unsupported"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCreateCSVParserFactory(t *testing.T) {
	factory := createCSVParserFactory()

	if factory == nil {
		t.Fatal("Factory function is nil")
	}

	// Test creating a parser
	reader := strings.NewReader("test,data")
	parser := factory(reader)

	if parser == nil {
		t.Fatal("Parser created by factory is nil")
	}

	// Verify it's not nil (type is already guaranteed by the factory function)
	// No need for type assertion since factory returns *parser.CSVParser
}

func TestShutdown(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	// Verify application has services before shutdown
	services := app.Container.GetRegisteredServices()
	if len(services) == 0 {
		t.Fatal("No services registered before shutdown")
	}

	err = app.Shutdown()
	if err != nil {
		t.Fatalf("Failed to shutdown application: %v", err)
	}

	// Verify container is cleared after shutdown
	services = app.Container.GetRegisteredServices()
	if len(services) != 0 {
		t.Errorf("Expected 0 services after shutdown, got %d", len(services))
	}
}

func TestShutdown_StoreError(t *testing.T) {
	// This test is more about ensuring the shutdown process handles errors gracefully
	// Since we're using memory store which has a simple Close() implementation,
	// this test verifies the shutdown flow works
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	err = app.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown should not fail with memory store: %v", err)
	}
}

// Test helper for creating test readers
func createTestReader(content string) io.Reader {
	return strings.NewReader(content)
}

func TestApplication_ConfigurationInjection(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	// Verify that the injected configuration matches expectations
	if app.Config.Storage.Type != "file" { // default for test environment
		t.Errorf("Expected default storage type 'file', got '%s'", app.Config.Storage.Type)
	}

	if app.Config.Server.Environment != "development" { // default
		t.Errorf("Expected default environment 'development', got '%s'", app.Config.Server.Environment)
	}
}

func TestApplication_DependencyResolutionChain(t *testing.T) {
	app, err := NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}

	// Test that we can resolve dependencies in a chain
	// Config -> Store -> Generator -> Parser

	// 1. Resolve config
	configInstance, err := app.Container.Resolve("config")
	if err != nil {
		t.Fatalf("Failed to resolve config: %v", err)
	}
	config := configInstance.(*config.Config)

	// 2. Resolve store
	store, err := app.GetCardStore()
	if err != nil {
		t.Fatalf("Failed to resolve store: %v", err)
	}

	// 3. Resolve generator
	generator, err := app.GetCardGenerator()
	if err != nil {
		t.Fatalf("Failed to resolve generator: %v", err)
	}

	// 4. Create parser
	parser, err := app.GetCSVParser(strings.NewReader("test"))
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	// Verify all dependencies are non-nil
	if config == nil || store == nil || generator == nil || parser == nil {
		t.Error("One or more dependencies in the chain is nil")
	}
}
