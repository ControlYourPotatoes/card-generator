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

    if s.TargetType != "" {
        validTargets := map[string]bool{
            "Creature": true,
            "Player":   true,
            "Any":      true,
        }
        if !validTargets[s.TargetType] {
            return &validation.ValidationError{
                Type:    validation.ErrorTypeInvalid,
                Message: "invalid target type",
                Field:   "targetType",
            }
        }
    }
    return nil
}

func (s *Spell) ToData() *card.CardData {
    data := s.BaseCard.ToData()
    data.TargetType = s.TargetType
    return data
}
