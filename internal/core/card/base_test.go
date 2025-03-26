package card

import (
	"testing"
	"time"
)

func TestBaseCardGetters(t *testing.T) {
	// Setup a test card
	baseCard := BaseCard{
		ID:        "1",
		Name:      "Test Card",
		Cost:      3,
		Effect:    "Test effect",
		Type:      TypeCreature,
		Keywords:  []Keyword{"HASTE", "CRITICAL"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]string{
			"artist": "Test Artist",
			"set":    "Core Set",
		},
	}

	// Test getters
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"GetID", baseCard.GetID(), "1"},
		{"GetName", baseCard.GetName(), "Test Card"},
		{"GetCost", baseCard.GetCost(), 3},
		{"GetEffect", baseCard.GetEffect(), "Test effect"},
		{"GetType", baseCard.GetType(), TypeCreature},
		{"GetKeywords length", len(baseCard.GetKeywords()), 2},
		{"GetKeywords[0]", baseCard.GetKeywords()[0], Keyword("HASTE")},
		{"GetMetadata length", len(baseCard.GetMetadata()), 2},
		{"GetMetadata[artist]", baseCard.GetMetadata()["artist"], "Test Artist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.got)
			}
		})
	}
}

func TestBaseCardValidate(t *testing.T) {
	tests := []struct {
		name        string
		card        BaseCard
		expectError bool
		errorType   error
	}{
		{
			name: "Valid card",
			card: BaseCard{
				Name:   "Valid Card",
				Cost:   3,
				Effect: "Valid effect",
				Type:   TypeCreature,
			},
			expectError: false,
		},
		{
			name: "Empty name",
			card: BaseCard{
				Name:   "",
				Cost:   3,
				Effect: "Valid effect",
				Type:   TypeCreature,
			},
			expectError: true,
			errorType:   ErrEmptyName,
		},
		{
			name: "Empty effect",
			card: BaseCard{
				Name:   "Valid Card",
				Cost:   3,
				Effect: "",
				Type:   TypeCreature,
			},
			expectError: true,
			errorType:   ErrEmptyEffect,
		},
		{
			name: "Negative cost",
			card: BaseCard{
				Name:   "Valid Card",
				Cost:   -2,
				Effect: "Valid effect",
				Type:   TypeCreature,
			},
			expectError: true,
			errorType:   ErrInvalidCost,
		},
		{
			name: "X cost (-1) is valid",
			card: BaseCard{
				Name:   "X Cost Card",
				Cost:   -1,
				Effect: "Valid effect",
				Type:   TypeSpell,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate()
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
				return
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}
			
			if tt.expectError {
				if err.Error() != tt.errorType.Error() {
					t.Errorf("Expected error %v but got %v", tt.errorType, err)
				}
			}
		})
	}
}

func TestBaseCardToDTO(t *testing.T) {
	// Setup test time
	now := time.Now().Truncate(time.Second)
	
	// Setup a test card
	baseCard := BaseCard{
		ID:        "1",
		Name:      "Test Card",
		Cost:      3,
		Effect:    "Test effect",
		Type:      TypeCreature,
		Keywords:  []Keyword{"HASTE", "CRITICAL"},
		CreatedAt: now,
		UpdatedAt: now,
		Metadata: map[string]string{
			"artist": "Test Artist",
			"set":    "Core Set",
		},
	}
	
	// Convert to DTO
	dto := baseCard.ToDTO()
	
	// Test DTO fields
	if dto.ID != baseCard.ID {
		t.Errorf("Expected ID %s, got %s", baseCard.ID, dto.ID)
	}
	
	if dto.Name != baseCard.Name {
		t.Errorf("Expected Name %s, got %s", baseCard.Name, dto.Name)
	}
	
	if dto.Cost != baseCard.Cost {
		t.Errorf("Expected Cost %d, got %d", baseCard.Cost, dto.Cost)
	}
	
	if dto.Effect != baseCard.Effect {
		t.Errorf("Expected Effect %s, got %s", baseCard.Effect, dto.Effect)
	}
	
	if dto.Type != baseCard.Type {
		t.Errorf("Expected Type %s, got %s", baseCard.Type, dto.Type)
	}
	
	if len(dto.Keywords) != len(baseCard.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(baseCard.Keywords), len(dto.Keywords))
	} else {
		for i, k := range baseCard.Keywords {
			if dto.Keywords[i] != string(k) {
				t.Errorf("Expected keyword %s, got %s", string(k), dto.Keywords[i])
			}
		}
	}
	
	if !dto.CreatedAt.Equal(baseCard.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", baseCard.CreatedAt, dto.CreatedAt)
	}
	
	if !dto.UpdatedAt.Equal(baseCard.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", baseCard.UpdatedAt, dto.UpdatedAt)
	}
	
	if len(dto.Metadata) != len(baseCard.Metadata) {
		t.Errorf("Expected %d metadata entries, got %d", len(baseCard.Metadata), len(dto.Metadata))
	} else {
		for k, v := range baseCard.Metadata {
			if dto.Metadata[k] != v {
				t.Errorf("Expected metadata %s=%s, got %s", k, v, dto.Metadata[k])
			}
		}
	}
}