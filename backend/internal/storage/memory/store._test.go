package memory

import (
	"testing"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

func TestMemoryStore(t *testing.T) {
	store := New()

	// Create a test card
	testCard := &card.Creature{
		BaseCard: card.BaseCard{
			Name:   "Test Creature",
			Cost:   1,
			Effect: "Test Effect",
			Type:   card.TypeCreature,
		},
		Attack:  1,
		Defense: 1,
		Trait:   "Test",
	}

	// Test Save
	t.Run("Save", func(t *testing.T) {
		id, err := store.Save(testCard)
		if err != nil {
			t.Errorf("Failed to save card: %v", err)
		}
		if id == "" {
			t.Error("Expected non-empty ID")
		}
	})

	// Test Load
	t.Run("Load", func(t *testing.T) {
		id := generateID(testCard)
		loaded, err := store.Load(id)
		if err != nil {
			t.Errorf("Failed to load card: %v", err)
		}

		if loaded.GetName() != testCard.GetName() {
			t.Errorf("Loaded card name = %v, want %v", loaded.GetName(), testCard.GetName())
		}
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		cards, err := store.List()
		if err != nil {
			t.Errorf("Failed to list cards: %v", err)
		}
		if len(cards) != 1 {
			t.Errorf("Expected 1 card, got %d", len(cards))
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		id := generateID(testCard)
		err := store.Delete(id)
		if err != nil {
			t.Errorf("Failed to delete card: %v", err)
		}

		// Verify card is deleted
		_, err = store.Load(id)
		if err == nil {
			t.Error("Expected error loading deleted card")
		}
	})

	// Clean up
	if err := store.Close(); err != nil {
		t.Errorf("Failed to close store: %v", err)
	}
}
