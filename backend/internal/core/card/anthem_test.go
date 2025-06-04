package card

import (
	"testing"
	"time"
)

func TestAnthemValidation(t *testing.T) {
	tests := []struct {
		name        string
		anthem      Anthem
		expectError bool
		errorField  string
	}{
		{
			name: "Valid anthem",
			anthem: Anthem{
				BaseCard: BaseCard{
					Name:   "Glorious Anthem",
					Cost:   3,
					Effect: "All creatures you control get +1/+1.",
					Type:   TypeAnthem,
				},
				Continuous: true,
			},
			expectError: false,
		},
		{
			name: "Non-continuous anthem",
			anthem: Anthem{
				BaseCard: BaseCard{
					Name:   "Invalid Anthem",
					Cost:   3,
					Effect: "All creatures you control get +1/+1.",
					Type:   TypeAnthem,
				},
				Continuous: false, // Invalid: anthems must be continuous
			},
			expectError: true,
			errorField:  "continuous",
		},
		{
			name: "Wrong type",
			anthem: Anthem{
				BaseCard: BaseCard{
					Name:   "Wrong Type",
					Cost:   2,
					Effect: "This has the wrong type.",
					Type:   TypeCreature, // Wrong type
				},
				Continuous: true,
			},
			expectError: true,
			errorField:  "type",
		},
		{
			name: "Base card validation error",
			anthem: Anthem{
				BaseCard: BaseCard{
					Name:   "", // Invalid: empty name
					Cost:   2,
					Effect: "Valid effect",
					Type:   TypeAnthem,
				},
				Continuous: true,
			},
			expectError: true,
			errorField:  "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.anthem.Validate()

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

func TestAnthemToDTO(t *testing.T) {
	// Setup a test anthem
	anthem := Anthem{
		BaseCard: BaseCard{
			ID:        "anthem-123",
			Name:      "Test Anthem",
			Cost:      3,
			Effect:    "All creatures you control get +1/+1.",
			Type:      TypeAnthem,
			Keywords:  []string{"BUFF"},
			CreatedAt: time.Now().Truncate(time.Second),
			UpdatedAt: time.Now().Truncate(time.Second),
			Metadata:  map[string]string{"artist": "Test Anthem Artist"},
		},
		Continuous: true,
	}

	// Convert to DTO
	dto := anthem.ToDTO()

	// Verify base fields
	if dto.ID != anthem.ID {
		t.Errorf("Expected ID %s, got %s", anthem.ID, dto.ID)
	}

	if dto.Name != anthem.Name {
		t.Errorf("Expected Name %s, got %s", anthem.Name, dto.Name)
	}

	if dto.Cost != anthem.Cost {
		t.Errorf("Expected Cost %d, got %d", anthem.Cost, dto.Cost)
	}

	if dto.Effect != anthem.Effect {
		t.Errorf("Expected Effect %s, got %s", anthem.Effect, dto.Effect)
	}

	// Verify anthem-specific fields
	if dto.Continuous != anthem.Continuous {
		t.Errorf("Expected Continuous %v, got %v", anthem.Continuous, dto.Continuous)
	}
}

func TestNewAnthemFromDTO(t *testing.T) {
	// Create a DTO
	dto := &CardDTO{
		ID:         "anthem-dto-456",
		Type:       TypeAnthem,
		Name:       "Anthem From DTO",
		Cost:       4,
		Effect:     "All creatures get +2/+2.",
		Continuous: true,
		Keywords:   []string{"BUFF"},
		Metadata:   map[string]string{"set": "Test Set"},
	}

	// Create anthem from DTO
	anthem := NewAnthemFromDTO(dto)

	// Verify base fields
	if anthem.ID != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, anthem.ID)
	}

	if anthem.Name != dto.Name {
		t.Errorf("Expected Name %s, got %s", dto.Name, anthem.Name)
	}

	if anthem.Cost != dto.Cost {
		t.Errorf("Expected Cost %d, got %d", dto.Cost, anthem.Cost)
	}

	if anthem.Effect != dto.Effect {
		t.Errorf("Expected Effect %s, got %s", dto.Effect, anthem.Effect)
	}

	if len(anthem.Keywords) != len(dto.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(dto.Keywords), len(anthem.Keywords))
	}

	// Verify anthem-specific fields
	if anthem.Continuous != dto.Continuous {
		t.Errorf("Expected Continuous %v, got %v", dto.Continuous, anthem.Continuous)
	}
}
