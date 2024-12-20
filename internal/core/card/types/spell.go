package types

import (
    "github.com/yourusername/cardgame/internal/core/card"
)

type Spell struct {
    card.BaseCard
    TargetType string
}

func (s *Spell) Validate() *card.ValidationError {
    if err := s.ValidateBase(); err != nil {
        return err
    }

    if s.TargetType != "" {
        validTargets := map[string]bool{
            "Creature": true,
            "Player":   true,
            "Any":      true,
        }
        if !validTargets[s.TargetType] {
            return &card.ValidationError{
                Type:    card.ErrorTypeInvalid,
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
