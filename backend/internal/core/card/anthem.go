package card

// Anthem represents an Anthem card (continuous effect)
type Anthem struct {
	BaseCard
	Continuous bool // Anthems should always be continuous
}

// Validate performs anthem-specific validation in addition to base card validation
func (a *Anthem) Validate() error {
	// First, validate the base card fields
	if err := a.BaseCard.Validate(); err != nil {
		return err
	}

	// Perform anthem-specific validation
	if a.Type != TypeAnthem {
		return NewValidationError("card type must be Anthem", "type")
	}

	// All anthems should be continuous
	if !a.Continuous {
		return NewValidationError("anthem must be continuous", "continuous")
	}

	return nil
}

// ToDTO converts Anthem to CardDTO
func (a *Anthem) ToDTO() *CardDTO {
	dto := a.BaseCard.ToDTO()
	dto.Continuous = a.Continuous
	return dto
}

// ToData maintains backward compatibility with the older interface
func (a *Anthem) ToData() *CardDTO {
	return a.ToDTO()
}

// NewAnthemFromDTO creates a new Anthem from CardDTO
func NewAnthemFromDTO(dto *CardDTO) *Anthem {
	return &Anthem{
		BaseCard:    NewBaseCardFromDTO(dto),
		Continuous:  dto.Continuous,
	}
}