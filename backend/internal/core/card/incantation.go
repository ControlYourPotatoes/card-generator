package card

import (
	"strings"
)

// Incantation represents an Incantation card
type Incantation struct {
	BaseCard
	Timing string // "ON ANY CLASH", "ON ATTACK", etc.
}

// Validate performs incantation-specific validation in addition to base card validation
func (i *Incantation) Validate() error {
	// First, validate the base card fields
	if err := i.BaseCard.Validate(); err != nil {
		return err
	}

	// Perform incantation-specific validation
	if i.Type != TypeIncantation {
		return NewValidationError("card type must be Incantation", "type")
	}

	// Validate timing if specified
	if i.Timing != "" {
		validTimings := map[string]bool{
			"ON ANY CLASH": true,
			"ON ATTACK":    true,
		}
		if !validTimings[i.Timing] {
			return NewValidationError("invalid timing", "timing")
		}
	}

	return nil
}

// ToDTO converts Incantation to CardDTO
func (i *Incantation) ToDTO() *CardDTO {
	dto := i.BaseCard.ToDTO()
	dto.Timing = i.Timing
	return dto
}

// ToData maintains backward compatibility with the older interface
func (i *Incantation) ToData() *CardDTO {
	return i.ToDTO()
}

// NewIncantationFromDTO creates a new Incantation from CardDTO
func NewIncantationFromDTO(dto *CardDTO) *Incantation {
	return &Incantation{
		BaseCard: NewBaseCardFromDTO(dto),
		Timing:   dto.Timing,
	}
}

// DetermineTiming analyzes the effect text to determine the timing
func DetermineTiming(effect string) string {
	if strings.Contains(effect, "ON ANY CLASH") {
		return "ON ANY CLASH"
	}
	if strings.Contains(effect, "ON ATTACK") {
		return "ON ATTACK"
	}
	return ""
}
