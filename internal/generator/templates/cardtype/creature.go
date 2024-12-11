package cardtypes

import (
	"image"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
	"github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type CreatureTemplate struct {
	templates.BaseTemplate
}

func NewCreatureTemplate() (*CreatureTemplate, error) {
	return &CreatureTemplate{
		BaseTemplate: templates.BaseTemplate{
			framesPath: filepath.Join("internal", "generator", "templates", "images"),
			artBounds:  templates.GetDefaultArtBounds(),
		},
	}, nil
}

func (t *CreatureTemplate) GetFrame(data *card.CardData) (image.Image, error) {
	// Choose between BaseCreature.png and SpecialCreature.png based on data
	framePath := filepath.Join(t.framesPath, "BaseCreature.png")
	if t.isSpecialFrame(data) {
		framePath = filepath.Join(t.framesPath, "SpecialCreature.png")
	}
	return templates.LoadFrame(framePath)
}

func (t *CreatureTemplate) GetTextBounds(data *card.CardData) *templates.TextBounds {
	bounds := templates.GetDefaultBounds()

	// Add creature-specific stats positioning
	bounds.Stats = &templates.StatsConfig{
		Left: templates.TextConfig{
			Bounds:    image.Rect(115, 1885, 280, 1955),
			FontSize:  70,
			Alignment: "center",
		},
		Right: templates.TextConfig{
			Bounds:    image.Rect(1220, 1885, 1385, 1955),
			FontSize:  70,
			Alignment: "center",
		},
	}

	return bounds
}

func (t *CreatureTemplate) GetArtBounds() image.Rectangle {
	return t.artBounds
}

func (t *CreatureTemplate) isSpecialFrame(data *card.CardData) bool {
	// Logic to determine if card should use special frame
	// Could be based on creature type, keywords, etc.
	return false
}
