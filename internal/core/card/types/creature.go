package types

import (
    "github.com/yourusername/cardgame/internal/core/card"
)

type Creature struct {
    card.BaseCard
    Attack  int
    Defense int
    Trait   string
}

func (c *Creature) Validate() *card.ValidationError {
    if err := c.ValidateBase(); err != nil {
        return err
    }
    
    if c.Attack < 0 {
        return &card.ValidationError{
            Type:    card.ErrorTypeRange,
            Message: "attack cannot be negative",
            Field:   "attack",
        }
    }
    if c.Defense < 0 {
        return &card.ValidationError{
            Type:    card.ErrorTypeRange,
            Message: "defense cannot be negative",
            Field:   "defense",
        }
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