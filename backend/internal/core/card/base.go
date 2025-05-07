// Package card defines the core domain model for cards
package card

import (
	"time"
)

// CardType represents the type of a card
type CardType string

const (
	TypeCreature    CardType = "Creature"
	TypeArtifact    CardType = "Artifact"
	TypeSpell       CardType = "Spell"
	TypeIncantation CardType = "Incantation"
	TypeAnthem      CardType = "Anthem"
)

// For backward compatibility - we'll use strings for keywords initially
// Later we'll migrate to a proper Keyword type
type Keyword = string

// Card defines the core interface that all cards must implement
type Card interface {
	GetID() string
	GetName() string
	GetCost() int
	GetEffect() string
	GetType() CardType
	GetKeywords() []string  
	GetMetadata() map[string]string
	Validate() error
	ToDTO() *CardDTO
	ToData() *CardDTO // For backward compatibility
}

// BaseCard provides common functionality for all card types
type BaseCard struct {
	ID        string            // Database ID
	Name      string
	Cost      int
	Effect    string
	Type      CardType
	Keywords  []string          
	CreatedAt time.Time
	UpdatedAt time.Time
	Metadata  map[string]string
}

// Base getters for BaseCard
func (b BaseCard) GetID() string               { return b.ID }
func (b BaseCard) GetName() string             { return b.Name }
func (b BaseCard) GetCost() int                { return b.Cost }
func (b BaseCard) GetEffect() string           { return b.Effect }
func (b BaseCard) GetType() CardType           { return b.Type }
func (b BaseCard) GetKeywords() []string       { return b.Keywords }
func (b BaseCard) GetMetadata() map[string]string { return b.Metadata }

// Validate performs basic validation on the card
// This will be overridden with a full implementation once we have the validation package
func (b BaseCard) Validate() error {
	// Basic validation - will be replaced with complete implementation
	if b.Name == "" {
		return ErrEmptyName
	}
	if b.Effect == "" {
		return ErrEmptyEffect
	}
	if b.Cost < -1 { // -1 is allowed for X costs
		return ErrInvalidCost
	}
	return nil
}

// Temporary error definitions - these will move to a validation package
var (
	ErrEmptyName    = NewValidationError("name cannot be empty", "name")
	ErrEmptyEffect  = NewValidationError("effect cannot be empty", "effect")
	ErrInvalidCost  = NewValidationError("cost cannot be negative (except -1 for X costs)", "cost")
)

// ValidationError defines a simple validation error
// This will be replaced by the validation package implementation
type ValidationError struct {
	Message string
	Field   string
}

func (e ValidationError) Error() string {
	return e.Message + " (field: " + e.Field + ")"
}

// NewValidationError creates a new validation error
func NewValidationError(message, field string) ValidationError {
	return ValidationError{
		Message: message,
		Field:   field,
	}
}