package types

// We'll just add any additional tagging-specific functionality here
// while using the core card types directly

import (
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
)

// We can add helper functions specific to tagging if needed
func ExtractTaggableData(c *card.CardData) map[string]interface{} {
    return map[string]interface{}{
        "type":      c.Type,
        "effect":    c.Effect,
        "keywords":  c.Keywords,
        "trait":     c.Trait,
        "cost":      c.Cost,
    }
}