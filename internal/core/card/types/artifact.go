package card

import (
    "strings"
    "github.com/yourusername/cardgame/internal/core/card/validation"
)

type Artifact struct {
    BaseCard
    IsEquipment bool
}

func (a *Artifact) Validate() *validation.ValidationError {
    // Validate base fields
    if err := validation.ValidateName(a.Name); err != nil {
        return err
    }
    if err := validation.ValidateCost(a.Cost); err != nil {
        return err
    }
    if err := validation.ValidateEffect(a.Effect); err != nil {
        return err
    }

    // Artifact-specific validation
    if a.IsEquipment && !strings.Contains(strings.ToLower(a.Effect), "equip") {
        return &validation.ValidationError{
            Type:    validation.ErrorTypeInvalid,
            Message: "equipment artifact must contain equip effect",
            Field:   "effect",
        }
    }

    return nil
}

func (a *Artifact) ToData() *CardData {
    data := a.BaseCard.ToData()
    data.IsEquipment = a.IsEquipment
    return data
}
