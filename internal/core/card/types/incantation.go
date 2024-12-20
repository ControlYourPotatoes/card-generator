package types

import (
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
)

type Incantation struct {
    card.BaseCard
    Timing string
}

func (i *Incantation) Validate() *card.ValidationError {
    if err := i.ValidateBase(); err != nil {
        return err
    }

    if i.Timing != "" {
        validTimings := map[string]bool{
            "ON ANY CLASH": true,
            "ON ATTACK":    true,
        }
        if !validTimings[i.Timing] {
            return &card.ValidationError{
                Type:    card.ErrorTypeInvalid,
                Message: "invalid timing",
                Field:   "timing",
            }
        }
    }
    return nil
}

func (i *Incantation) ToData() *card.CardData {
    data := i.BaseCard.ToData()
    data.Timing = i.Timing
    return data
}