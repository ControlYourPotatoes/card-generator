package cardtypes

import (
	"image"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
	"github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type SpellTemplate struct {
	templates.BaseTemplate
}

func NewSpellTemplate() (*SpellTemplate, error) {
	return &SpellTemplate{
		BaseTemplate: templates.BaseTemplate{
			framesPath: filepath.Join("internal", "generator", "templates", "images"),
			artBounds:  templates.GetDefaultArtBounds(),
		},
	}, nil
}

func (t *SpellTemplate) GetFrame(data *card.CardData) (image.Image, error) {
	return templates.LoadFrame(filepath.Join(t.framesPath, "BaseSpell.png"))
}

func (t *SpellTemplate) GetTextBounds(data *card.CardData) *templates.TextBounds {
	return templates.GetDefaultBounds()
}

func (t *SpellTemplate) GetArtBounds() image.Rectangle {
	return t.artBounds
}

func (t *SpellTemplate) isSpecialFrame(data *card.CardData) bool {
	// Logic to determine if card should use special frame
	// Could be based on creature type, keywords, etc.
	return false
}
