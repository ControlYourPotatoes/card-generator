package card

import "strings"

// Artifact represents an artifact card
type Artifact struct {
	BaseCard
	IsEquipment bool
}

// Validate performs validation specific to Artifact cards
func (a *Artifact) Validate() error {
	// First validate the base card
	if err := a.BaseCard.Validate(); err != nil {
		return err
	}
	
	// Validate artifact-specific rules
	if a.IsEquipment && !strings.Contains(strings.ToLower(a.Effect), "equip") {
		return NewValidationError("equipment artifact must reference 'equip' in its effect", "effect")
	}
	
	return nil
}

// ToDTO converts Artifact to CardDTO
func (a *Artifact) ToDTO() *CardDTO {
	dto := a.BaseCard.ToDTO()
	dto.IsEquipment = a.IsEquipment
	return dto
}

// ToData is for backward compatibility
func (a *Artifact) ToData() *CardDTO {
	return a.ToDTO()
}

// NewArtifactFromDTO creates an Artifact from a CardDTO
func NewArtifactFromDTO(dto *CardDTO) *Artifact {
	return &Artifact{
		BaseCard:    NewBaseCardFromDTO(dto),
		IsEquipment: dto.IsEquipment,
	}
}