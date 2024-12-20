package types

import (
	
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
    "github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"	
)

type Spell struct {
    card.BaseCard
    TargetType string
}

func (s *Spell) Validate() *validation.ValidationError {
    
    baseValidator := validation.BaseValidator{
        Name: s.Name,
        Cost: s.Cost,
        Effect: s.Effect,
    }

    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }

    // Validate spell-specific properties
    if err := validation.ValidateSpell(s.TargetType); err != nil {
        return err
    }

    return nil
}

func (s *Spell) ToData() *card.CardData {
    data := s.BaseCard.ToData()
    data.TargetType = s.TargetType
    return data
}
