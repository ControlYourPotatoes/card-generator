package types

import (
    "github.com/yourusername/cardgame/internal/core/card"
)

type Anthem struct {
    card.BaseCard
    Continuous bool
}

func (a *Anthem) Validate() *card.ValidationError {
    if err := a.ValidateBase(); err != nil {
        return err
    }

    if !a.Continuous {
        return &card.ValidationError{
            Type:    card.ErrorTypeInvalid,
            Message: "anthem must be continuous",
            Field:   "continuous",
        }
    }
    return nil
}

func (a *Anthem) ToData() *card.CardData {
    data := a.BaseCard.ToData()
    data.Continuous = a.Continuous
    return data
}