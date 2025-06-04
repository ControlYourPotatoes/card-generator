package card

import (
	"encoding/json"
	"testing"
	"time"
)

func TestToDTO(t *testing.T) {
	// Setup test time
	now := time.Now().Truncate(time.Second)

	// Setup a test card
	baseCard := BaseCard{
		ID:        "123",
		Name:      "Test DTO Card",
		Cost:      4,
		Effect:    "Test DTO effect",
		Type:      TypeSpell,
		Keywords:  []string{"CRITICAL", "HASTE"},
		CreatedAt: now,
		UpdatedAt: now,
		Metadata: map[string]string{
			"artist": "Test Artist",
			"rarity": "Rare",
		},
	}

	// Convert to DTO
	dto := baseCard.ToDTO()

	// Verify DTO fields
	if dto.ID != baseCard.ID {
		t.Errorf("Expected ID %s, got %s", baseCard.ID, dto.ID)
	}

	if dto.Name != baseCard.Name {
		t.Errorf("Expected Name %s, got %s", baseCard.Name, dto.Name)
	}

	if dto.Cost != baseCard.Cost {
		t.Errorf("Expected Cost %d, got %d", baseCard.Cost, dto.Cost)
	}

	if dto.Type != baseCard.Type {
		t.Errorf("Expected Type %s, got %s", baseCard.Type, dto.Type)
	}

	if dto.Effect != baseCard.Effect {
		t.Errorf("Expected Effect %s, got %s", baseCard.Effect, dto.Effect)
	}

	if len(dto.Keywords) != len(baseCard.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(baseCard.Keywords), len(dto.Keywords))
	} else {
		for i, k := range baseCard.Keywords {
			if dto.Keywords[i] != k {
				t.Errorf("Expected keyword %s, got %s", k, dto.Keywords[i])
			}
		}
	}

	if !dto.CreatedAt.Equal(baseCard.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", baseCard.CreatedAt, dto.CreatedAt)
	}

	if !dto.UpdatedAt.Equal(baseCard.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", baseCard.UpdatedAt, dto.UpdatedAt)
	}

	for k, v := range baseCard.Metadata {
		if dto.Metadata[k] != v {
			t.Errorf("Expected metadata %s=%s, got %s", k, v, dto.Metadata[k])
		}
	}
}

func TestDTORoundTrip(t *testing.T) {
	// Setup test time
	now := time.Now().Truncate(time.Second)

	// Create a DTO
	originalDTO := &CardDTO{
		ID:        "456",
		Type:      TypeCreature,
		Name:      "Creature Card",
		Cost:      2,
		Effect:    "Creature effect",
		Attack:    3,
		Defense:   4,
		Trait:     "Beast",
		Keywords:  []string{"FACEOFF"},
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  map[string]string{"set": "Core"},
	}

	// Convert to BaseCard
	baseCard := NewBaseCardFromDTO(originalDTO)

	// Convert back to DTO
	roundTripDTO := baseCard.ToDTO()

	// Verify ID, Name, Cost, Effect, Type
	if roundTripDTO.ID != originalDTO.ID {
		t.Errorf("Expected ID %s, got %s", originalDTO.ID, roundTripDTO.ID)
	}

	if roundTripDTO.Name != originalDTO.Name {
		t.Errorf("Expected Name %s, got %s", originalDTO.Name, roundTripDTO.Name)
	}

	if roundTripDTO.Cost != originalDTO.Cost {
		t.Errorf("Expected Cost %d, got %d", originalDTO.Cost, roundTripDTO.Cost)
	}

	if roundTripDTO.Effect != originalDTO.Effect {
		t.Errorf("Expected Effect %s, got %s", originalDTO.Effect, roundTripDTO.Effect)
	}

	if roundTripDTO.Type != originalDTO.Type {
		t.Errorf("Expected Type %s, got %s", originalDTO.Type, roundTripDTO.Type)
	}

	// Verify Keywords
	if len(roundTripDTO.Keywords) != len(originalDTO.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(originalDTO.Keywords), len(roundTripDTO.Keywords))
	} else {
		for i, k := range originalDTO.Keywords {
			if roundTripDTO.Keywords[i] != k {
				t.Errorf("Expected keyword %s, got %s", k, roundTripDTO.Keywords[i])
			}
		}
	}

	// Verify timestamps
	if !roundTripDTO.CreatedAt.Equal(originalDTO.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", originalDTO.CreatedAt, roundTripDTO.CreatedAt)
	}

	if !roundTripDTO.UpdatedAt.Equal(originalDTO.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", originalDTO.UpdatedAt, roundTripDTO.UpdatedAt)
	}

	// Verify metadata
	for k, v := range originalDTO.Metadata {
		if roundTripDTO.Metadata[k] != v {
			t.Errorf("Expected metadata %s=%s, got %s", k, v, roundTripDTO.Metadata[k])
		}
	}

	// Note: Type-specific fields like Attack, Defense, etc. are lost
	// in BaseCard conversion - this is expected and will be handled
	// by specific card type implementations
}

func TestDTOJsonSerialization(t *testing.T) {
	// Create a DTO with various fields
	dto := &CardDTO{
		ID:          "789",
		Type:        TypeArtifact,
		Name:        "Magic Artifact",
		Cost:        5,
		Effect:      "Do something magical",
		IsEquipment: true,
		Keywords:    []string{"INDESTRUCTIBLE"},
		CreatedAt:   time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
		Metadata:    map[string]string{"artist": "Leonardo", "set": "Magic Set"},
	}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Failed to marshal DTO to JSON: %v", err)
	}

	// Deserialize back to a DTO
	var deserializedDTO CardDTO
	err = json.Unmarshal(jsonBytes, &deserializedDTO)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to DTO: %v", err)
	}

	// Verify fields
	if deserializedDTO.ID != dto.ID {
		t.Errorf("Expected ID %s, got %s", dto.ID, deserializedDTO.ID)
	}

	if deserializedDTO.Name != dto.Name {
		t.Errorf("Expected Name %s, got %s", dto.Name, deserializedDTO.Name)
	}

	if deserializedDTO.Cost != dto.Cost {
		t.Errorf("Expected Cost %d, got %d", dto.Cost, deserializedDTO.Cost)
	}

	if deserializedDTO.Effect != dto.Effect {
		t.Errorf("Expected Effect %s, got %s", dto.Effect, deserializedDTO.Effect)
	}

	if deserializedDTO.Type != dto.Type {
		t.Errorf("Expected Type %s, got %s", dto.Type, deserializedDTO.Type)
	}

	if deserializedDTO.IsEquipment != dto.IsEquipment {
		t.Errorf("Expected IsEquipment %v, got %v", dto.IsEquipment, deserializedDTO.IsEquipment)
	}

	if len(deserializedDTO.Keywords) != len(dto.Keywords) {
		t.Errorf("Expected %d keywords, got %d", len(dto.Keywords), len(deserializedDTO.Keywords))
	} else {
		for i, k := range dto.Keywords {
			if deserializedDTO.Keywords[i] != k {
				t.Errorf("Expected keyword %s, got %s", k, deserializedDTO.Keywords[i])
			}
		}
	}

	// Time comparison needs to handle JSON serialization peculiarities
	if !deserializedDTO.CreatedAt.Equal(dto.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", dto.CreatedAt, deserializedDTO.CreatedAt)
	}

	if !deserializedDTO.UpdatedAt.Equal(dto.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", dto.UpdatedAt, deserializedDTO.UpdatedAt)
	}

	// Verify metadata fields
	for k, v := range dto.Metadata {
		if deserializedDTO.Metadata[k] != v {
			t.Errorf("Expected metadata %s=%s, got %s", k, v, deserializedDTO.Metadata[k])
		}
	}
}

func TestToDataCompatibility(t *testing.T) {
	// Setup a test card
	baseCard := BaseCard{
		ID:       "compatibility-123",
		Name:     "Compatibility Test Card",
		Cost:     1,
		Effect:   "Test backward compatibility",
		Type:     TypeIncantation,
		Keywords: []string{"HASTE"},
	}

	// Test that ToData() and ToDTO() return equivalent results
	dataResult := baseCard.ToData()
	dtoResult := baseCard.ToDTO()

	// Verify they have the same values
	if dataResult.ID != dtoResult.ID {
		t.Errorf("ToData.ID != ToDTO.ID: %s vs %s", dataResult.ID, dtoResult.ID)
	}

	if dataResult.Name != dtoResult.Name {
		t.Errorf("ToData.Name != ToDTO.Name: %s vs %s", dataResult.Name, dtoResult.Name)
	}

	if dataResult.Cost != dtoResult.Cost {
		t.Errorf("ToData.Cost != ToDTO.Cost: %d vs %d", dataResult.Cost, dtoResult.Cost)
	}

	if dataResult.Effect != dtoResult.Effect {
		t.Errorf("ToData.Effect != ToDTO.Effect: %s vs %s", dataResult.Effect, dtoResult.Effect)
	}

	if dataResult.Type != dtoResult.Type {
		t.Errorf("ToData.Type != ToDTO.Type: %s vs %s", dataResult.Type, dtoResult.Type)
	}
}
