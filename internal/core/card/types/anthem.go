package types

import (
	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"
)

type Anthem struct {
    card.BaseCard
    Continuous bool
}

func (c *Anthem) Validate() *validation.ValidationError {
    
    
    baseValidator := validation.BaseValidator{
		Name: c.Name,
		Cost: c.Cost,
		Effect: c.Effect,
	}
    
    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }

    return nil
}

func (a *Anthem) ToData() *card.CardData {
    data := a.BaseCard.ToData()
    data.Continuous = a.Continuous
    return data
}