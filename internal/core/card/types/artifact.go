package types

import (
    "strings"
    
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
)

type Artifact struct {
    card.BaseCard
    IsEquipment bool
}

func (a *Artifact) Validate() *card.ValidationError {
    if err := a.ValidateBase(); err != nil {
        return err
    }

    if a.IsEquipment && !strings.Contains(strings.ToLower(a.Effect), "equip") {
        return &card.ValidationError{
            Type:    card.ErrorTypeInvalid,
            Message: "equipment artifact must contain equip effect",
            Field:   "effect",
        }
    }
    return nil
}

func (a *Artifact) ToData() *card.CardData {
    data := a.BaseCard.ToData()
    data.IsEquipment = a.IsEquipment
    return data
}