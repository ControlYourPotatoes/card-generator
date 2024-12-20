package types

import (


    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"
)

type Creature struct {
    card.BaseCard
    Attack  int
    Defense int
    Trait   string
}

func (c *Creature) Validate() *validation.ValidationError {
    
	baseValidator := validation.BaseValidator{
		Name: c.Name,
		Cost: c.Cost,
		Effect: c.Effect,
	}
    
    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }
    
    // Validate creature-specific properties
    if err := validation.ValidateCreature(c.Attack, c.Defense); err != nil {
        return err
    }
    
    return nil
}

func (c *Creature) ToData() *card.CardData {
    data := c.BaseCard.ToData()
    data.Attack = c.Attack
    data.Defense = c.Defense
    data.Trait = c.Trait
    return data
}