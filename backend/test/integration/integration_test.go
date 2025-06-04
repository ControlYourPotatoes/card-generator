package integration

import (
	"strings"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/backend/pkg/bootstrap"
)

// TestApplicationIntegration tests the full application bootstrap and component integration
func TestApplicationIntegration(t *testing.T) {
	app, err := bootstrap.NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}
	defer app.Shutdown()

	// Test that all components can be resolved and work together
	store, err := app.GetCardStore()
	if err != nil {
		t.Fatalf("Failed to get card store: %v", err)
	}

	generator, err := app.GetCardGenerator()
	if err != nil {
		t.Fatalf("Failed to get card generator: %v", err)
	}

	csvParser, err := app.GetCSVParser(strings.NewReader("test,data"))
	if err != nil {
		t.Fatalf("Failed to get CSV parser: %v", err)
	}

	// Verify all components are not nil
	if store == nil || generator == nil || csvParser == nil {
		t.Error("One or more components are nil")
	}

	// TODO: Add more comprehensive integration tests as the system grows
	// - Test card parsing → storage → generation workflow
	// - Test configuration changes affecting components
	// - Test error propagation between components
}

// TestConfigurationIntegration tests configuration loading and component configuration
func TestConfigurationIntegration(t *testing.T) {
	app, err := bootstrap.NewApplication("test")
	if err != nil {
		t.Fatalf("Failed to create application: %v", err)
	}
	defer app.Shutdown()

	// Verify configuration is properly injected into components
	if app.Config == nil {
		t.Fatal("Application config is nil")
	}

	// Test that storage type from config affects the actual storage implementation
	if app.Config.Storage.Type != "file" {
		t.Errorf("Expected storage type 'file', got %s", app.Config.Storage.Type)
	}

	// TODO: Add tests for:
	// - Different storage configurations
	// - Generator configuration affecting output
	// - Database configuration affecting connections
}

// TestErrorHandlingIntegration tests error propagation across components
func TestErrorHandlingIntegration(t *testing.T) {
	// TODO: Implement tests for:
	// - Invalid configuration handling
	// - Component failure scenarios
	// - Graceful degradation
	// - Error logging and reporting
	t.Skip("Error handling integration tests not yet implemented")
}

// TestPerformanceIntegration tests performance characteristics of integrated components
func TestPerformanceIntegration(t *testing.T) {
	// TODO: Implement tests for:
	// - Memory usage with multiple components
	// - Concurrent access to shared resources
	// - Resource cleanup verification
	t.Skip("Performance integration tests not yet implemented")
}
