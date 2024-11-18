package store

import (
	"github.com/ControlYourPotatoes/card-generator/internal/card"
)

// Store defines the interface for card storage operations
type Store interface {
	// Save stores a card and returns its ID
	Save(card card.Card) (string, error)

	// Load retrieves a card by its ID
	Load(id string) (card.Card, error)

	// List returns all stored cards
	List() ([]card.Card, error)

	// Delete removes a card by its ID
	Delete(id string) error

	// Close cleans up any resources
	Close() error
}