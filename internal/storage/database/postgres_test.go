package database_test

import (
	"os"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/internal/storage/database"
	"github.com/joho/godotenv"
)

func setupTestDB(t *testing.T) *database.PostgresStore {
	// Load environment variables
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Logf("Warning: .env file not found, using environment variables")
	}

	// Create database manager
	manager, err := database.NewManager("")
	if err != nil {
		t.Fatalf("Failed to create database manager: %v", err)
	}

	// Connect to database
	if err := manager.Connect(); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	
	// Initialize test data
	if err := manager.InitWithTestData(); err != nil {
		t.Fatalf("Failed to initialize test data: %v", err)
	}
	
	// Get store
	store, err := manager.GetStore()
	if err != nil {
		t.Fatalf("Failed to get store: %v", err)
	}

	return store
}

func TestPostgresStoreBasicOperations(t *testing.T) {
	// Skip if we're not in test mode with a database
	if os.Getenv("RUN_DB_TESTS") != "true" {
		t.Skip("Skipping database tests. Set RUN_DB_TESTS=true to run")
	}

	store := setupTestDB(t)
	defer store.Close()

	// Test card creation
	t.Run("CreateCard", func(t *testing.T) {
		testCard := &card.Spell{
			BaseCard: card.BaseCard{
				Name:     "Test Spell",
				Cost:     2,
				Effect:   "Deal 2 damage to target creature.",
				Type:     card.TypeSpell,
				Keywords: []string{"DAMAGE"},
			},
			TargetType: "Creature",
		}

		id, err := store.Save(testCard)
		if err != nil {
			t.Fatalf("Failed to save card: %v", err)
		}

		if id == "" {
			t.Fatal("Expected non-empty ID")
		}

		// Test loading the card
		loadedCard, err := store.Load(id)
		if err != nil {
			t.Fatalf("Failed to load card: %v", err)
		}

		if loadedCard.GetName() != testCard.GetName() {
			t.Errorf("Expected name %s, got %s", testCard.GetName(), loadedCard.GetName())
		}

		if loadedCard.GetType() != testCard.GetType() {
			t.Errorf("Expected type %s, got %s", testCard.GetType(), loadedCard.GetType())
		}

		// Cast to spell and check specific field
		spellCard, ok := loadedCard.(*card.Spell)
		if !ok {
			t.Fatalf("Expected *card.Spell, got %T", loadedCard)
		}

		if spellCard.TargetType != testCard.TargetType {
			t.Errorf("Expected target type %s, got %s", testCard.TargetType, spellCard.TargetType)
		}

		// Test deleting the card
		err = store.Delete(id)
		if err != nil {
			t.Fatalf("Failed to delete card: %v", err)
		}

		// Verify it's deleted
		_, err = store.Load(id)
		if err == nil {
			t.Error("Expected error loading deleted card")
		}
	})

	// Test listing cards
	t.Run("ListCards", func(t *testing.T) {
		cards, err := store.List()
		if err != nil {
			t.Fatalf("Failed to list cards: %v", err)
		}

		// Should have at least the seed cards
		if len(cards) < 3 {
			t.Errorf("Expected at least 3 cards, got %d", len(cards))
		}
	})

	// Test different card types
	t.Run("CreateDifferentCardTypes", func(t *testing.T) {
		// Create a creature
		creature := &card.Creature{
			BaseCard: card.BaseCard{
				Name:     "Test Creature",
				Cost:     3,
				Effect:   "Test effect",
				Type:     card.TypeCreature,
				Keywords: []string{"HASTE"},
			},
			Attack:  2,
			Defense: 2,
			Trait:   card.TraitBeast,
		}

		creatureID, err := store.Save(creature)
		if err != nil {
			t.Fatalf("Failed to save creature: %v", err)
		}

		// Create an artifact
		artifact := &card.Artifact{
			BaseCard: card.BaseCard{
				Name:     "Test Artifact",
				Cost:     2,
				Effect:   "Equip to a creature. Equipped creature gets +1/+1.",
				Type:     card.TypeArtifact,
				Keywords: []string{"EQUIPMENT"},
			},
			IsEquipment: true,
		}

		artifactID, err := store.Save(artifact)
		if err != nil {
			t.Fatalf("Failed to save artifact: %v", err)
		}

		// Load and verify
		loadedCreature, err := store.Load(creatureID)
		if err != nil {
			t.Fatalf("Failed to load creature: %v", err)
		}

		loadedArtifact, err := store.Load(artifactID)
		if err != nil {
			t.Fatalf("Failed to load artifact: %v", err)
		}

		// Verify creature
		creatureObj, ok := loadedCreature.(*card.Creature)
		if !ok {
			t.Fatalf("Expected *card.Creature, got %T", loadedCreature)
		}

		if creatureObj.Attack != creature.Attack {
			t.Errorf("Expected attack %d, got %d", creature.Attack, creatureObj.Attack)
		}

		// Verify artifact
		artifactObj, ok := loadedArtifact.(*card.Artifact)
		if !ok {
			t.Fatalf("Expected *card.Artifact, got %T", loadedArtifact)
		}

		if artifactObj.IsEquipment != artifact.IsEquipment {
			t.Errorf("Expected isEquipment %v, got %v", artifact.IsEquipment, artifactObj.IsEquipment)
		}

		// Clean up
		_ = store.Delete(creatureID)
		_ = store.Delete(artifactID)
	})
}