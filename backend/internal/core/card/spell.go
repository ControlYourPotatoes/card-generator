package card

import (
	"strings"
)

// Spell represents a Spell card
type Spell struct {
	BaseCard
	TargetType string // "Creature", "Player", or "Any"
}

// Validate performs spell-specific validation in addition to base card validation
func (s *Spell) Validate() error {
	// First, validate the base card fields
	if err := s.BaseCard.Validate(); err != nil {
		return err
	}

	// Perform spell-specific validation
	if s.Type != TypeSpell {
		return NewValidationError("card type must be Spell", "type")
	}

	// Validate target type if specified
	if s.TargetType != "" {
		validTargets := map[string]bool{
			"Creature": true,
			"Player":   true,
			"Any":      true,
		}
		if !validTargets[s.TargetType] {
			return NewValidationError("invalid target type", "targetType")
		}
	}

	return nil
}

// ToDTO converts Spell to CardDTO
func (s *Spell) ToDTO() *CardDTO {
	dto := s.BaseCard.ToDTO()
	dto.TargetType = s.TargetType
	return dto
}

// ToData maintains backward compatibility with the older interface
func (s *Spell) ToData() *CardDTO {
	return s.ToDTO()
}

// NewSpellFromDTO creates a new Spell from CardDTO
func NewSpellFromDTO(dto *CardDTO) *Spell {
	return &Spell{
		BaseCard:   NewBaseCardFromDTO(dto),
		TargetType: dto.TargetType,
	}
}

// DetermineTargetType analyzes the effect text to determine the target type
func DetermineTargetType(effect string) string {
	effect = strings.ToLower(effect)
	if strings.Contains(effect, "target creature") {
		return "Creature"
	}
	if strings.Contains(effect, "target player") {
		return "Player"
	}
	return "Any"
}
