package database

import (
    "fmt"
    "os"
    "testing"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

// getTestConnection returns a connection string for testing
func getTestConnection() string {
    // Default test connection, override with environment variables if needed
    host := getEnvOrDefault("TEST_DB_HOST", "localhost")
    port := getEnvOrDefault("TEST_DB_PORT", "5432")
    user := getEnvOrDefault("TEST_DB_USER", "postgres")
    password := getEnvOrDefault("TEST_DB_PASSWORD", "postgres")
    dbname := getEnvOrDefault("TEST_DB_NAME", "cardgame_test")
    
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname,
    )
}


// TestPostgresStore tests the PostgreSQL store implementation
func TestPostgresStore(t *testing.T) {
    // Skip if we're not running database tests
    if os.Getenv("SKIP_DB_TESTS") == "true" {
        t.Skip("Skipping database tests")
    }
    
    // Get test connection
    connString := getTestConnection()
    
    // Create store
    store, err := NewPostgresStore(connString)
    if err != nil {
        t.Fatalf("Failed to create store: %v", err)
    }
    defer store.Close()
    
    // Initialize schema
    if err := store.InitSchema(); err != nil {
        t.Fatalf("Failed to initialize schema: %v", err)
    }
    
    // Create a test card
    testCard := &card.Creature{
        BaseCard: card.BaseCard{
            Name:   "Test Creature",
            Cost:   3,
            Effect: "Test effect",
            Type:   card.TypeCreature,
        },
        Attack:  2,
        Defense: 2,
    }
    
    // Save the card
    id, err := store.Save(testCard)
    if err != nil {
        t.Fatalf("Failed to save card: %v", err)
    }
    
    // Load the card
    loadedCard, err := store.Load(id)
    if err != nil {
        t.Fatalf("Failed to load card: %v", err)
    }
    
    // Verify card data
    if loadedCard.GetName() != testCard.GetName() {
        t.Errorf("Expected name %s, got %s", testCard.GetName(), loadedCard.GetName())
    }
    
    if loadedCard.GetCost() != testCard.GetCost() {
        t.Errorf("Expected cost %d, got %d", testCard.GetCost(), loadedCard.GetCost())
    }
    
    if loadedCard.GetEffect() != testCard.GetEffect() {
        t.Errorf("Expected effect %s, got %s", testCard.GetEffect(), loadedCard.GetEffect())
    }
    
    // Test listing cards
    cards, err := store.List()
    if err != nil {
        t.Fatalf("Failed to list cards: %v", err)
    }
    
    if len(cards) < 1 {
        t.Error("Expected at least one card in list")
    }
    
    // Test deleting a card
    if err := store.Delete(id); err != nil {
        t.Fatalf("Failed to delete card: %v", err)
    }
    
    // Verify card is deleted
    _, err = store.Load(id)
    if err == nil {
        t.Error("Expected error when loading deleted card")
    }
}