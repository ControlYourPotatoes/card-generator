package card

import (
	"testing"
	"time"
)

func TestCreatureValidation(t *testing.T) {
	tests := []struct {
		name        string
		creature    Creature
		expectError bool
		errorField  string
	}{
		{
			name: "Valid creature",
			creature: Creature{
				BaseCard: BaseCard{
					Name:   "Valid Creature",
					Cost:   3,
					Effect: "Valid effect",
					Type:   TypeCreature,
				},
				Attack:  2,
				Defense: 3,
				Trait:   TraitBeast,
			},
			expectError: false,
		},
		{
			name: "Negative attack",
			creature: Creature{
				BaseCard: BaseCard{
					Name:   "Negative Attack",
					Cost:   3,
					Effect: "Valid effect",
					Type:   TypeCreature,
				},
				Attack:  -1,
				Defense: 3,
				Trait:   TraitBeast,
			},
			expectError: true,
			errorField:  "attack",
		},
		{
			name: "Negative defense",
			creature: Creature{
				BaseCard: BaseCard{
					Name:   "Negative Defense",
					Cost:   3,
					Effect: "Valid effect",
					Type:   TypeCreature,
				},
				Attack:  2,
				Defense: -1,
				Trait:   TraitBeast,
			},
			expectError: true,
			errorField:  "defense",
		},
		{
			name: "Invalid trait",
			creature: Creature{
				BaseCard: BaseCard{
					Name:   "Invalid Trait",
					Cost:   3,
					Effect: "Valid effect",
					Type:   TypeCreature,
				},
				Attack:  2,
				Defense: 3,
				Trait:   "InvalidTrait",
			},
			expectError: true,
			errorField:  "trait",
		},
		{
			name: "Empty trait is valid",
			creature: Creature{
				BaseCard: BaseCard{
					Name:   "No Trait",
					Cost:   3,
					Effect: "Valid effect",
					Type:   TypeCreature,
				},
				Attack:  2,
				Defense: 3,
				Trait:   "",
			},
			expectError: false,
		},
		{
			name: "Base card validation error",
			creature: Creature{
				BaseCard: BaseCard{
					Name:   "", // Invalid: empty name
					Cost:   3,
					Effect: "Valid effect",
					Type:   TypeCreature,
				},
				Attack:  2,
				Defense: 3,
				Trait:   TraitBeast,
			},
			expectError: true,
			errorField:  "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.creature.Validate()

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if tt.expectError {
				valErr, ok := err.(ValidationError)
				if !ok {
					t.Errorf("Expected ValidationError but got different error type: %T", err)
					return
				}

				if valErr.Field != tt.errorField {
					t.Errorf("Expected error on field %s but got error on field %s", tt.errorField, valErr.Field)
				}
			}
		})
	}
}

func TestCreatureToDTO(t *testing.T) {
	// Setup a test creature
	creature := Creature{
		BaseCard: BaseCard{
			ID:        "creature-123",
			Name:      "Test Creature",
			Cost:      5,
			Effect:    "Test creature effect",
			Type:      TypeCreature,
			Keywords:  []string{"CRITICAL"},
			CreatedAt: time.Now().Truncate(time.Second),
			UpdatedAt: time.Now().Truncate(time.Second),
			Metadata:  map[string]string{"artist": "Test Creature Artist"},
		},
		Attack:  3,
		Defense: 4,
		Trait:   TraitDragon,
	}

	// Convert to DTO
	dto := creature.ToDTO()

	// Verify base fields
	if dto.ID != creature.ID {
		t.Errorf("Expected ID %s, got %s", creature.ID, dto.ID)
	}

	if dto.Name != creature.Name {
		t.Errorf("Expected Name %s, got %s", creature.Name, dto.Name)
	}

	if dto.Cost != creature.Cost {
		t.Errorf("Expected Cost %d, got %d", creature.Cost, dto.Cost)
	}

	if dto.Effect != creature.Effect {
		t.Errorf("Expected Effect %s, got %s", creature.Effect, dto.Effect)
	}

	// Verify creature-specific fields
	if dto.Attack != creature.Attack {
		t.Errorf("Expected Attack %d, got %d", creature.Attack, dto.Attack)
	}

	if dto.Defense != creature.Defense {
		t.Errorf("Expected Defense %d, got %d", creature.Defense, dto.Defense)
	}

	if dto.Trait != string(creature.Trait) {
		t.Errorf("Expected Trait %s, got %s", creature.Trait, dto.Trait)
	}
}

func TestNewCreatureFromDTO(t *testing.T) {
	// Create a DTO
	dto := &CardDTO{
		ID:       "creature-dto-456",
		Type:     TypeCreature,
		Name:     "Creature From DTO",
		Cost:     4,
		Effect:   "Created from DTO",
		Attack:   5,
		Defense:  5,
		Trait:    string(TraitAngel),
		Keywords: []string{"HASTE", "DOUBLE STRIKE"},
		Metadata: map[string]string{"set": "Test Set"},
	}

	// Create creature from DTO
	creature := NewCreatureFromDTO(dto)

	// Verify base fields
	if creature.ID != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, creature.ID)
	}

	if creature.Name != dto.Name {
		t.Errorf("Expected Name %s, got %s", dto.Name, creature.Name)
	}

	if creature.Cost != dto.Cost {
		t.Errorf("Expected Cost %d, got %d", dto.Cost, creature.Cost)
	}

	if creature.Effect != dto.Effect {
		t.Errorf("Expected Effect %s, got %s", dto.Effect, creature.Effect)
	}

	if len(creature.Keywords) != len(dto.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(dto.Keywords), len(creature.Keywords))
	}

	// Verify creature-specific fields
	if creature.Attack != dto.Attack {
		t.Errorf("Expected Attack %d, got %d", dto.Attack, creature.Attack)
	}

	if creature.Defense != dto.Defense {
		t.Errorf("Expected Defense %d, got %d", dto.Defense, creature.Defense)
	}

	if string(creature.Trait) != dto.Trait {
		t.Errorf("Expected Trait %s, got %s", dto.Trait, creature.Trait)
	}
}

func TestTraitValidation(t *testing.T) {
	validTraits := []Trait{
		TraitBeast,
		TraitWarrior,
		TraitDragon,
		TraitDemon,
		TraitAngel,
		TraitLegendary,
		TraitAncient,
		TraitDivine,
	}

	for _, trait := range validTraits {
		if !trait.IsValid() {
			t.Errorf("Expected trait %s to be valid", trait)
		}
	}

	invalidTraits := []Trait{
		"Unknown",
		"Monster",
		"Human",
		"Elf",
	}

	for _, trait := range invalidTraits {
		if trait.IsValid() {
			t.Errorf("Expected trait %s to be invalid", trait)
		}
	}
}
