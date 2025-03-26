package card

import (
	"testing"
	"time"
)

func TestArtifactValidation(t *testing.T) {
	tests := []struct {
		name        string
		artifact    Artifact
		expectError bool
		errorField  string
	}{
		{
			name: "Valid non-equipment artifact",
			artifact: Artifact{
				BaseCard: BaseCard{
					Name:   "Valid Artifact",
					Cost:   2,
					Effect: "When this enters play, draw a card.",
					Type:   TypeArtifact,
				},
				IsEquipment: false,
			},
			expectError: false,
		},
		{
			name: "Valid equipment artifact",
			artifact: Artifact{
				BaseCard: BaseCard{
					Name:   "Sword of Testing",
					Cost:   3,
					Effect: "Equip to a Creature. Equipped creature gets +2/+2.",
					Type:   TypeArtifact,
				},
				IsEquipment: true,
			},
			expectError: false,
		},
		{
			name: "Wrong type",
			artifact: Artifact{
				BaseCard: BaseCard{
					Name:   "Wrong Type",
					Cost:   2,
					Effect: "This has the wrong type.",
					Type:   TypeCreature, // Wrong type
				},
				IsEquipment: false,
			},
			expectError: true,
			errorField:  "type",
		},
		{
			name: "Base card validation error",
			artifact: Artifact{
				BaseCard: BaseCard{
					Name:   "", // Invalid: empty name
					Cost:   2,
					Effect: "Valid effect",
					Type:   TypeArtifact,
				},
				IsEquipment: false,
			},
			expectError: true,
			errorField:  "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.artifact.Validate()
			
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

func TestArtifactToDTO(t *testing.T) {
	// Setup a test artifact
	artifact := Artifact{
		BaseCard: BaseCard{
			ID:        "artifact-123",
			Name:      "Test Artifact",
			Cost:      2,
			Effect:    "Equip to a Creature. Equipped creature gains +1/+1.",
			Type:      TypeArtifact,
			Keywords:  []string{"EQUIPMENT"},
			CreatedAt: time.Now().Truncate(time.Second),
			UpdatedAt: time.Now().Truncate(time.Second),
			Metadata:  map[string]string{"artist": "Test Artifact Artist"},
		},
		IsEquipment: true,
	}
	
	// Convert to DTO
	dto := artifact.ToDTO()
	
	// Verify base fields
	if dto.ID != artifact.ID {
		t.Errorf("Expected ID %s, got %s", artifact.ID, dto.ID)
	}
	
	if dto.Name != artifact.Name {
		t.Errorf("Expected Name %s, got %s", artifact.Name, dto.Name)
	}
	
	if dto.Cost != artifact.Cost {
		t.Errorf("Expected Cost %d, got %d", artifact.Cost, dto.Cost)
	}
	
	if dto.Effect != artifact.Effect {
		t.Errorf("Expected Effect %s, got %s", artifact.Effect, dto.Effect)
	}
	
	// Verify artifact-specific fields
	if dto.IsEquipment != artifact.IsEquipment {
		t.Errorf("Expected IsEquipment %v, got %v", artifact.IsEquipment, dto.IsEquipment)
	}
}

func TestNewArtifactFromDTO(t *testing.T) {
	// Create a DTO
	dto := &CardDTO{
		ID:          "artifact-dto-456",
		Type:        TypeArtifact,
		Name:        "Artifact From DTO",
		Cost:        4,
		Effect:      "Equip :2. Equipped creature gets +2/+0.",
		IsEquipment: true,
		Keywords:    []string{"EQUIPMENT"},
		Metadata:    map[string]string{"set": "Test Set"},
	}
	
	// Create artifact from DTO
	artifact := NewArtifactFromDTO(dto)
	
	// Verify base fields
	if artifact.ID != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, artifact.ID)
	}
	
	if artifact.Name != dto.Name {
		t.Errorf("Expected Name %s, got %s", dto.Name, artifact.Name)
	}
	
	if artifact.Cost != dto.Cost {
		t.Errorf("Expected Cost %d, got %d", dto.Cost, artifact.Cost)
	}
	
	if artifact.Effect != dto.Effect {
		t.Errorf("Expected Effect %s, got %s", dto.Effect, artifact.Effect)
	}
	
	if len(artifact.Keywords) != len(dto.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(dto.Keywords), len(artifact.Keywords))
	}
	
	// Verify artifact-specific fields
	if artifact.IsEquipment != dto.IsEquipment {
		t.Errorf("Expected IsEquipment %v, got %v", dto.IsEquipment, artifact.IsEquipment)
	}
}

func TestDetermineIsEquipment(t *testing.T) {
	tests := []struct {
		effect          string
		expectedIsEquip bool
	}{
		{
			effect:          "Equip to a creature. Equipped creature gains +1/+1.",
			expectedIsEquip: true,
		},
		{
			effect:          "EQUIP this to a warrior for extra damage.",
			expectedIsEquip: true,
		},
		{
			effect:          "This artifact can be equipped to any beast.",
			expectedIsEquip: true,
		},
		{
			effect:          "When this enters play, draw a card.",
			expectedIsEquip: false,
		},
		{
			effect:          "Sacrifice this: Deal 2 damage to target creature.",
			expectedIsEquip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.effect, func(t *testing.T) {
			result := DetermineIsEquipment(tt.effect)
			if result != tt.expectedIsEquip {
				t.Errorf("Expected DetermineIsEquipment to return %v for effect %q, but got %v", 
					tt.expectedIsEquip, tt.effect, result)
			}
		})
	}
}

func TestAutosetEquipment(t *testing.T) {
	// Test that artifacts with "equip" or "equipped" in their effect 
	// automatically get IsEquipment set to true
	dto := &CardDTO{
		Type:   TypeArtifact,
		Name:   "Auto Equipment",
		Cost:   3,
		Effect: "Equip: 2",
		// IsEquipment not explicitly set
	}
	
	artifact := NewArtifactFromDTO(dto)
	
	if !artifact.IsEquipment {
		t.Errorf("Expected IsEquipment to be automatically set to true for effect containing 'equipped'")
	}
	
	// Also check that the EQUIPMENT keyword was added
	hasEquipmentKeyword := false
	for _, keyword := range artifact.Keywords {
		if keyword == "EQUIPMENT" {
			hasEquipmentKeyword = true
			break
		}
	}
	
	if !hasEquipmentKeyword {
		t.Errorf("Expected EQUIPMENT keyword to be automatically added")
	}
}