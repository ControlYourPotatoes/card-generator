package card

import (
    "github.com/yourusername/cardgame/internal/core/card/validation"
)

type Creature struct {
    BaseCard
    Attack  int
    Defense int
    Trait   string
}

func (c *Creature) Validate() *validation.ValidationError {
    // Validate base fields
    if err := validation.ValidateName(c.Name); err != nil {
        return err
    }
    if err := validation.ValidateCost(c.Cost); err != nil {
        return err
    }
    if err := validation.ValidateEffect(c.Effect); err != nil {
        return err
    }

    // Creature-specific validation
    if c.Attack < 0 {
        return &validation.ValidationError{
            Type:    validation.ErrorTypeRange,
            Message: "attack cannot be negative",
            Field:   "attack",
        }
    }
    if c.Defense < 0 {
        return &validation.ValidationError{
            Type:    validation.ErrorTypeRange,
            Message: "defense cannot be negative",
            Field:   "defense",
        }
    }

    return nil
}

func (c *Creature) ToData() *CardData {
    data := c.BaseCard.ToData()
    data.Attack = c.Attack
    data.Defense = c.Defense
    data.Trait = c.Trait
    return data
}
