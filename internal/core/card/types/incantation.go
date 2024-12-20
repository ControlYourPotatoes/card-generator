package types

import (
	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"
)

type Incantation struct {
    card.BaseCard
    Timing string
}

func (i *Incantation) Validate() *validation.ValidationError {
    
    baseValidator := validation.BaseValidator{
        Name: i.Name,
        Cost: i.Cost,
        Effect: i.Effect,
    }

    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }

    return nil
}

func (i *Incantation) ToData() *card.CardData {
    data := i.BaseCard.ToData()
    data.Timing = i.Timing
    return data
}