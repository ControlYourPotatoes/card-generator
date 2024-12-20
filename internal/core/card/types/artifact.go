package types

import (
	"strings"

	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"
)

type Artifact struct {
    card.BaseCard
    IsEquipment bool
}

func (a *Artifact) Validate() *validation.ValidationError {
    
    baseValidator := validation.BaseValidator{
        Name: a.Name,
        Cost: a.Cost,
        Effect: a.Effect,
    } 
    
    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }

    if a.IsEquipment && !strings.Contains(strings.ToLower(a.Effect), "equip") {
        return &validation.ValidationError{
            Type:    validation.ErrorTypeInvalid,
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