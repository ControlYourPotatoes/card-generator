package card

import (
    "github.com/yourusername/cardgame/internal/core/card/validation"
)

type Spell struct {
    BaseCard
    TargetType string // "Creature", "Player", "Any"
}

func (s *Spell) Validate() *validation.ValidationError {
    // Validate base fields
    if err := validation.ValidateName(s.Name); err != nil {
        return err
    }
    if err := validation.ValidateCost(s.Cost); err != nil {
        return err
    }
    if err := validation.ValidateEffect(s.Effect); err != nil {
        return err
    }

    // Spell-specific validation
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

func (s *Spell) ToData() *CardData {
    data := s.BaseCard.ToData()
    data.TargetType = s.TargetType
    return data
}
