package types

// We'll just add any additional tagging-specific functionality here
// while using the core card types directly

import (
	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// We can add helper functions specific to tagging if needed
func ExtractTaggableData(c *card.CardDTO) map[string]interface{} {
	return map[string]interface{}{
		"type":     c.Type,
		"effect":   c.Effect,
		"keywords": c.Keywords,
		"cost":     c.Cost,
		"trait":    c.Trait,
		"attack":   c.Attack,
		"defense":  c.Defense,
	}
}
