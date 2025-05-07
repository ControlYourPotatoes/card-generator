package card

import (
	"strings"
)

// Artifact represents an Artifact card
type Artifact struct {
	BaseCard
	IsEquipment bool
}

// Validate performs artifact-specific validation in addition to base card validation
func (a *Artifact) Validate() error {
	// First, validate the base card fields
	if err := a.BaseCard.Validate(); err != nil {
		return err
	}

	// Perform artifact-specific validation
	if a.Type != TypeArtifact {
		return NewValidationError("card type must be Artifact", "type")
	}

	// If it's marked as equipment, ensure the effect mentions "equip" or "equipped"
	if a.IsEquipment && !hasEquipmentText(a.Effect) {
		return NewValidationError("equipment artifact must contain equip effect", "effect")
	}

	return nil
}

// ToDTO converts Artifact to CardDTO
func (a *Artifact) ToDTO() *CardDTO {
	dto := a.BaseCard.ToDTO()
	dto.IsEquipment = a.IsEquipment
	return dto
}

// ToData maintains backward compatibility with the older interface
func (a *Artifact) ToData() *CardDTO {
	return a.ToDTO()
}

// NewArtifactFromDTO creates a new Artifact from CardDTO
func NewArtifactFromDTO(dto *CardDTO) *Artifact {
	artifact := &Artifact{
		BaseCard:    NewBaseCardFromDTO(dto),
		IsEquipment: dto.IsEquipment,
	}
	
	// If it has equipment-related text but IsEquipment flag is not set,
	// automatically set it
	if !artifact.IsEquipment && hasEquipmentText(artifact.Effect) {
		artifact.IsEquipment = true
	}
	
	// Add EQUIPMENT keyword if missing
	if artifact.IsEquipment {
		hasEquipmentKeyword := false
		for _, keyword := range artifact.Keywords {
			if strings.ToUpper(keyword) == "EQUIPMENT" {
				hasEquipmentKeyword = true
				break
			}
		}
		
		if !hasEquipmentKeyword {
			artifact.Keywords = append(artifact.Keywords, "EQUIPMENT")
		}
	}
	
	return artifact
}

// DetermineIsEquipment checks if an artifact is an equipment type
// based on its effect text (useful for parsing)
func DetermineIsEquipment(effect string) bool {
	return hasEquipmentText(effect)
}

// hasEquipmentText checks if the text contains equipment-related keywords
func hasEquipmentText(text string) bool {
	effectLower := strings.ToLower(text)
	return strings.Contains(effectLower, "equip") || 
	       strings.Contains(effectLower, "equippment") ||
		   strings.Contains(effectLower, "equipped")
}