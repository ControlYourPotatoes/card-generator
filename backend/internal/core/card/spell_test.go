package card

import (
	"testing"
	"time"
)

func TestSpellValidation(t *testing.T) {
	tests := []struct {
		name        string
		spell       Spell
		expectError bool
		errorField  string
	}{
		{
			name: "Valid spell with no target type",
			spell: Spell{
				BaseCard: BaseCard{
					Name:   "Global Effect",
					Cost:   3,
					Effect: "All creatures get -1/-1 until end of turn.",
					Type:   TypeSpell,
				},
				TargetType: "",
			},
			expectError: false,
		},
		{
			name: "Valid spell with creature target",
			spell: Spell{
				BaseCard: BaseCard{
					Name:   "Lightning Bolt",
					Cost:   1,
					Effect: "Deal 3 damage to target creature.",
					Type:   TypeSpell,
				},
				TargetType: "Creature",
			},
			expectError: false,
		},
		{
			name: "Valid spell with player target",
			spell: Spell{
				BaseCard: BaseCard{
					Name:   "Mind Blast",
					Cost:   2,
					Effect: "Target player discards a card.",
					Type:   TypeSpell,
				},
				TargetType: "Player",
			},
			expectError: false,
		},
		{
			name: "Valid spell with any target",
			spell: Spell{
				BaseCard: BaseCard{
					Name:   "Arcane Blast",
					Cost:   2,
					Effect: "Deal 2 damage to any target.",
					Type:   TypeSpell,
				},
				TargetType: "Any",
			},
			expectError: false,
		},
		{
			name: "Invalid target type",
			spell: Spell{
				BaseCard: BaseCard{
					Name:   "Invalid Target",
					Cost:   2,
					Effect: "This has an invalid target type.",
					Type:   TypeSpell,
				},
				TargetType: "Invalid", // Invalid target type
			},
			expectError: true,
			errorField:  "targetType",
		},
		{
			name: "Wrong type",
			spell: Spell{
				BaseCard: BaseCard{
					Name:   "Wrong Type",
					Cost:   2,
					Effect: "This has the wrong type.",
					Type:   TypeCreature, // Wrong type
				},
				TargetType: "Creature",
			},
			expectError: true,
			errorField:  "type",
		},
		{
			name: "Base card validation error",
			spell: Spell{
				BaseCard: BaseCard{
					Name:   "", // Invalid: empty name
					Cost:   2,
					Effect: "Valid effect",
					Type:   TypeSpell,
				},
				TargetType: "Creature",
			},
			expectError: true,
			errorField:  "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spell.Validate()

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

func TestSpellToDTO(t *testing.T) {
	// Setup a test spell
	spell := Spell{
		BaseCard: BaseCard{
			ID:        "spell-123",
			Name:      "Test Spell",
			Cost:      2,
			Effect:    "Deal 3 damage to target creature.",
			Type:      TypeSpell,
			Keywords:  []string{"DAMAGE"},
			CreatedAt: time.Now().Truncate(time.Second),
			UpdatedAt: time.Now().Truncate(time.Second),
			Metadata:  map[string]string{"artist": "Test Spell Artist"},
		},
		TargetType: "Creature",
	}

	// Convert to DTO
	dto := spell.ToDTO()

	// Verify base fields
	if dto.ID != spell.ID {
		t.Errorf("Expected ID %s, got %s", spell.ID, dto.ID)
	}

	if dto.Name != spell.Name {
		t.Errorf("Expected Name %s, got %s", spell.Name, dto.Name)
	}

	if dto.Cost != spell.Cost {
		t.Errorf("Expected Cost %d, got %d", spell.Cost, dto.Cost)
	}

	if dto.Effect != spell.Effect {
		t.Errorf("Expected Effect %s, got %s", spell.Effect, dto.Effect)
	}

	// Verify spell-specific fields
	if dto.TargetType != spell.TargetType {
		t.Errorf("Expected TargetType %s, got %s", spell.TargetType, dto.TargetType)
	}
}

func TestNewSpellFromDTO(t *testing.T) {
	// Create a DTO
	dto := &CardDTO{
		ID:         "spell-dto-456",
		Type:       TypeSpell,
		Name:       "Spell From DTO",
		Cost:       4,
		Effect:     "Deal 2 damage to target player.",
		TargetType: "Player",
		Keywords:   []string{"DAMAGE"},
		Metadata:   map[string]string{"set": "Test Set"},
	}

	// Create spell from DTO
	spell := NewSpellFromDTO(dto)

	// Verify base fields
	if spell.ID != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, spell.ID)
	}

	if spell.Name != dto.Name {
		t.Errorf("Expected Name %s, got %s", dto.Name, spell.Name)
	}

	if spell.Cost != dto.Cost {
		t.Errorf("Expected Cost %d, got %d", dto.Cost, spell.Cost)
	}

	if spell.Effect != dto.Effect {
		t.Errorf("Expected Effect %s, got %s", dto.Effect, spell.Effect)
	}

	if len(spell.Keywords) != len(dto.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(dto.Keywords), len(spell.Keywords))
	}

	// Verify spell-specific fields
	if spell.TargetType != dto.TargetType {
		t.Errorf("Expected TargetType %s, got %s", dto.TargetType, spell.TargetType)
	}
}

func TestDetermineTargetType(t *testing.T) {
	tests := []struct {
		effect             string
		expectedTargetType string
	}{
		{
			effect:             "Deal 3 damage to target creature.",
			expectedTargetType: "Creature",
		},
		{
			effect:             "TARGET CREATURE gets -2/-2 until end of turn.",
			expectedTargetType: "Creature",
		},
		{
			effect:             "Target player discards two cards.",
			expectedTargetType: "Player",
		},
		{
			effect:             "Deal 1 damage to TARGET PLAYER.",
			expectedTargetType: "Player",
		},
		{
			effect:             "Return all creatures to their owners' hands.",
			expectedTargetType: "Any",
		},
		{
			effect:             "All players draw a card.",
			expectedTargetType: "Any",
		},
	}

	for _, tt := range tests {
		t.Run(tt.effect, func(t *testing.T) {
			result := DetermineTargetType(tt.effect)
			if result != tt.expectedTargetType {
				t.Errorf("Expected DetermineTargetType to return %v for effect %q, but got %v",
					tt.expectedTargetType, tt.effect, result)
			}
		})
	}
}
