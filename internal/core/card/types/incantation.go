package types

import (

	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"
	"github.com/ControlYourPotatoes/card-generator/internal/core/common"
    
)

type Incantation struct {
    card.BaseCard
    Timing string
}

func (i *Incantation) Validate() *common.ValidationError {
    
    baseValidator := validation.BaseValidator{
        Name: i.Name,
        Cost: i.Cost,
        Effect: i.Effect,
    }

    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }

    // Validate incantation-specific properties
    if err := validation.ValidateIncantation(i.Timing); err != nil {
        return err
    }

    return nil
}

func (i *Incantation) ToData() *card.CardData {
    data := i.BaseCard.ToData()
    data.Timing = i.Timing
    return data
}