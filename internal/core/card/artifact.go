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

	// If it's marked as equipment, ensure the effect mentions "equip" as a word
	if a.IsEquipment {
		effectLower := strings.ToLower(a.Effect)
		
		// Check for word boundaries around "equip" but not "equipment"
		hasEquip := false
		
		// Use a more precise approach with regex
		if strings.Contains(effectLower, "equip") {
			// Check if it's not part of "equipment"
			if !strings.Contains(effectLower, "equipment") || 
			   (strings.Contains(effectLower, "equip ") || 
			   strings.Contains(effectLower, " equip") || 
			   strings.HasPrefix(effectLower, "equip") || 
			   strings.HasSuffix(effectLower, "equip") ||
			   strings.Contains(effectLower, "equip.") ||
			   strings.Contains(effectLower, "equip,") ||
			   strings.Contains(effectLower, "equip:") ||
			   strings.Contains(effectLower, "equip;")) {
				hasEquip = true
			}
		}
		
		if !hasEquip {
			return NewValidationError("equipment artifact must contain equip effect", "effect")
		}
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
	return &Artifact{
		BaseCard:    NewBaseCardFromDTO(dto),
		IsEquipment: dto.IsEquipment,
	}
}

// DetermineIsEquipment checks if an artifact is an equipment type
// based on its effect text (useful for parsing)
func DetermineIsEquipment(effect string) bool {
	effectLower := strings.ToLower(effect)
	
	// Check for word boundaries around "equip" but not "equipment"
	if strings.Contains(effectLower, "equip") {
		// Check if it's not just part of "equipment" or has clear word boundaries
		if !strings.Contains(effectLower, "equipment") || 
		   strings.Contains(effectLower, "equip ") || 
		   strings.Contains(effectLower, " equip") || 
		   strings.HasPrefix(effectLower, "equip") || 
		   strings.HasSuffix(effectLower, "equip") ||
		   strings.Contains(effectLower, "equip.") ||
		   strings.Contains(effectLower, "equip,") ||
		   strings.Contains(effectLower, "equip:") ||
		   strings.Contains(effectLower, "equip;") {
			return true
		}
	}
	
	return false
}