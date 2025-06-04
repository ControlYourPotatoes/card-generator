package card

import (
	"testing"
	"time"
)

func TestIncantationValidation(t *testing.T) {
	tests := []struct {
		name        string
		incantation Incantation
		expectError bool
		errorField  string
	}{
		{
			name: "Valid incantation with no timing",
			incantation: Incantation{
				BaseCard: BaseCard{
					Name:   "Flexible Incantation",
					Cost:   2,
					Effect: "Counter target spell.",
					Type:   TypeIncantation,
				},
				Timing: "",
			},
			expectError: false,
		},
		{
			name: "Valid incantation with ON ANY CLASH timing",
			incantation: Incantation{
				BaseCard: BaseCard{
					Name:   "Clash Incantation",
					Cost:   1,
					Effect: "ON ANY CLASH: Draw a card.",
					Type:   TypeIncantation,
				},
				Timing: "ON ANY CLASH",
			},
			expectError: false,
		},
		{
			name: "Valid incantation with ON ATTACK timing",
			incantation: Incantation{
				BaseCard: BaseCard{
					Name:   "Attack Incantation",
					Cost:   2,
					Effect: "ON ATTACK: Deal 1 damage to target creature.",
					Type:   TypeIncantation,
				},
				Timing: "ON ATTACK",
			},
			expectError: false,
		},
		{
			name: "Invalid timing",
			incantation: Incantation{
				BaseCard: BaseCard{
					Name:   "Invalid Timing",
					Cost:   2,
					Effect: "This has an invalid timing.",
					Type:   TypeIncantation,
				},
				Timing: "ON DEFENSE", // Invalid timing
			},
			expectError: true,
			errorField:  "timing",
		},
		{
			name: "Wrong type",
			incantation: Incantation{
				BaseCard: BaseCard{
					Name:   "Wrong Type",
					Cost:   2,
					Effect: "This has the wrong type.",
					Type:   TypeCreature, // Wrong type
				},
				Timing: "ON ATTACK",
			},
			expectError: true,
			errorField:  "type",
		},
		{
			name: "Base card validation error",
			incantation: Incantation{
				BaseCard: BaseCard{
					Name:   "", // Invalid: empty name
					Cost:   2,
					Effect: "Valid effect",
					Type:   TypeIncantation,
				},
				Timing: "ON ANY CLASH",
			},
			expectError: true,
			errorField:  "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.incantation.Validate()

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

func TestIncantationToDTO(t *testing.T) {
	// Setup a test incantation
	incantation := Incantation{
		BaseCard: BaseCard{
			ID:        "incantation-123",
			Name:      "Test Incantation",
			Cost:      2,
			Effect:    "ON ATTACK: Deal 2 damage to target creature.",
			Type:      TypeIncantation,
			Keywords:  []string{"DAMAGE"},
			CreatedAt: time.Now().Truncate(time.Second),
			UpdatedAt: time.Now().Truncate(time.Second),
			Metadata:  map[string]string{"artist": "Test Incantation Artist"},
		},
		Timing: "ON ATTACK",
	}

	// Convert to DTO
	dto := incantation.ToDTO()

	// Verify base fields
	if dto.ID != incantation.ID {
		t.Errorf("Expected ID %s, got %s", incantation.ID, dto.ID)
	}

	if dto.Name != incantation.Name {
		t.Errorf("Expected Name %s, got %s", incantation.Name, dto.Name)
	}

	if dto.Cost != incantation.Cost {
		t.Errorf("Expected Cost %d, got %d", incantation.Cost, dto.Cost)
	}

	if dto.Effect != incantation.Effect {
		t.Errorf("Expected Effect %s, got %s", incantation.Effect, dto.Effect)
	}

	// Verify incantation-specific fields
	if dto.Timing != incantation.Timing {
		t.Errorf("Expected Timing %s, got %s", incantation.Timing, dto.Timing)
	}
}

func TestNewIncantationFromDTO(t *testing.T) {
	// Create a DTO
	dto := &CardDTO{
		ID:       "incantation-dto-456",
		Type:     TypeIncantation,
		Name:     "Incantation From DTO",
		Cost:     4,
		Effect:   "ON ANY CLASH: Counter target spell.",
		Timing:   "ON ANY CLASH",
		Keywords: []string{"COUNTER"},
		Metadata: map[string]string{"set": "Test Set"},
	}

	// Create incantation from DTO
	incantation := NewIncantationFromDTO(dto)

	// Verify base fields
	if incantation.ID != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, incantation.ID)
	}

	if incantation.Name != dto.Name {
		t.Errorf("Expected Name %s, got %s", dto.Name, incantation.Name)
	}

	if incantation.Cost != dto.Cost {
		t.Errorf("Expected Cost %d, got %d", dto.Cost, incantation.Cost)
	}

	if incantation.Effect != dto.Effect {
		t.Errorf("Expected Effect %s, got %s", dto.Effect, incantation.Effect)
	}

	if len(incantation.Keywords) != len(dto.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(dto.Keywords), len(incantation.Keywords))
	}

	// Verify incantation-specific fields
	if incantation.Timing != dto.Timing {
		t.Errorf("Expected Timing %s, got %s", dto.Timing, incantation.Timing)
	}
}

func TestDetermineTiming(t *testing.T) {
	tests := []struct {
		effect         string
		expectedTiming string
	}{
		{
			effect:         "ON ANY CLASH: Draw a card.",
			expectedTiming: "ON ANY CLASH",
		},
		{
			effect:         "When your opponent declares attackers, ON ANY CLASH: Counter target spell.",
			expectedTiming: "ON ANY CLASH",
		},
		{
			effect:         "ON ATTACK: Deal 2 damage to target creature.",
			expectedTiming: "ON ATTACK",
		},
		{
			effect:         "ON ATTACK: Your creatures get +1/+1 until end of turn.",
			expectedTiming: "ON ATTACK",
		},
		{
			effect:         "Counter target spell.",
			expectedTiming: "",
		},
		{
			effect:         "Draw two cards.",
			expectedTiming: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.effect, func(t *testing.T) {
			result := DetermineTiming(tt.effect)
			if result != tt.expectedTiming {
				t.Errorf("Expected DetermineTiming to return %v for effect %q, but got %v",
					tt.expectedTiming, tt.effect, result)
			}
		})
	}
}
