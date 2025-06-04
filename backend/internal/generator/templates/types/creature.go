package types

import (
	"image"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/base"
)

type CreatureTemplate struct {
	base.BaseTemplate // Embed the base template directly
}

func NewCreatureTemplate() (*CreatureTemplate, error) {
	baseTemplate := base.NewBaseTemplate()
	return &CreatureTemplate{
		BaseTemplate: *baseTemplate, // Note the dereference here
	}, nil
}

// GetFrame implements the Template interface
func (t *CreatureTemplate) GetFrame(data *card.CardDTO) (image.Image, error) {
	var frameName string
	if t.isSpecialFrame(data) {
		frameName = "SpecialCreatureWithStats.png"
	} else {
		frameName = "BaseCreature.png"
	}
	return t.LoadFrame(frameName)
}

func (t *CreatureTemplate) GetTextBounds(data *card.CardDTO) map[string]image.Rectangle {
	bounds := make(map[string]image.Rectangle)

	// Basic text bounds
	bounds["name"] = image.Rect(125, 90, 1375, 170)
	bounds["cost"] = image.Rect(125, 90, 1375, 170)
	bounds["type"] = image.Rect(125, 1885, 1375, 1955)
	bounds["effect"] = image.Rect(160, 1250, 1340, 1750)
	bounds["collector"] = image.Rect(110, 2010, 750, 2090)

	// Add stats positioning for creatures
	bounds["attack"] = image.Rect(130, 1820, 230, 1900)
	bounds["defense"] = image.Rect(1270, 1820, 1370, 1900)

	return bounds
}

func (t *CreatureTemplate) isSpecialFrame(data *card.CardDTO) bool {
	// Check for special traits or keywords that would trigger special frame
	specialTraits := []string{"Legendary", "Ancient", "Divine"}
	for _, trait := range specialTraits {
		if strings.Contains(data.Trait, trait) {
			return true
		}
	}

	// Check for special keywords in effect text
	specialKeywords := []string{"CRITICAL", "DOUBLE STRIKE", "INDESTRUCTIBLE"}
	for _, keyword := range specialKeywords {
		if strings.Contains(data.Effect, keyword) {
			return true
		}
	}

	return false
}
