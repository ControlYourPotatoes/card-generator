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

// Card defines the core interface that all cards must implement
type Card interface {
	GetID() string
	GetName() string
	GetCost() int
	GetEffect() string
	GetType() CardType
	GetKeywords() []Keyword
	GetMetadata() map[string]string
	Validate() error
	ToDTO() *CardDTO
}

// BaseCard provides common functionality for all card types
type BaseCard struct {
	ID        string            // Database ID
	Name      string
	Cost      int
	Effect    string
	Type      CardType
	Keywords  []Keyword
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
func (b BaseCard) GetKeywords() []Keyword      { return b.Keywords }
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

// CardDTO struct will be defined in dto.go
// Temporary placeholder to make things compile
type CardDTO struct {
	ID          string            `json:"id,omitempty"`
	Type        CardType          `json:"type"`
	Name        string            `json:"name"`
	Cost        int               `json:"cost"`
	Effect      string            `json:"effect"`
	Keywords    []string          `json:"keywords,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Temporary implementation of ToDTO
// Will be replaced with proper implementation in dto.go
func (b BaseCard) ToDTO() *CardDTO {
	// Placeholder implementation
	keywordStrings := make([]string, len(b.Keywords))
	for i, k := range b.Keywords {
		keywordStrings[i] = string(k)
	}
	
	return &CardDTO{
		ID:        b.ID,
		Type:      b.Type,
		Name:      b.Name,
		Cost:      b.Cost,
		Effect:    b.Effect,
		Keywords:  keywordStrings,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		Metadata:  b.Metadata,
	}
}

// Keyword type will be defined in types.go
// Temporary definition to make things compile
type Keyword string