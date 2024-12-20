package card

import (
    "github.com/yourusername/cardgame/internal/core/card/validation"
)

type Incantation struct {
    BaseCard
    Timing string
}

func (i *Incantation) Validate() *validation.ValidationError {
    // Validate base fields
    if err := validation.ValidateName(i.Name); err != nil {
        return err
    }
    if err := validation.ValidateCost(i.Cost); err != nil {
        return err
    }
    if err := validation.ValidateEffect(i.Effect); err != nil {
        return err
    }

    // Incantation-specific validation
    if i.Timing != "" {
        validTimings := map[string]bool{
            "ON ANY CLASH": true,
            "ON ATTACK":    true,
        }
        if !validTimings[i.Timing] {
            return &validation.ValidationError{
                Type:    validation.ErrorTypeInvalid,
                Message: "invalid timing",
                Field:   "timing",
            }
        }
    }

    return nil
}

func (i *Incantation) ToData() *CardData {
    data := i.BaseCard.ToData()
    data.Timing = i.Timing
    return data
}
