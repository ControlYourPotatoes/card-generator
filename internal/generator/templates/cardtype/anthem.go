package cardtypes

import (
	"path/filepath"
	"image"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
	"github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type NewAnthemTemplate() (*AnthemTemplate, error) {
	return &AnthemTemplate{
		BaseTemplate: templates.BaseTemplate{
			framesPath: filepath.Join("internal", "generator", "templates", "images"),
			artBounds:  templates.GetDefaultArtBounds(),
		},
	}, nil
}

func (t *AnthemTemplate) GetFrame(data *card.CardData) (image.Image, error) {
	return templates.LoadFrame(filepath.Join(t.framesPath, "BaseAnthem.png"))
}

func (t *AnthemTemplate) GetTextBounds(data *card.CardData) *templates.TextBounds {
	return templates.GetDefaultBounds()
}

func (t *AnthemTemplate) GetArtBounds() image.Rectangle {
	return t.artBounds
}

func (t *AnthemTemplate) isSpecialFrame(data *card.CardData) bool {
	// Logic to determine if card should use special frame
	// Could be based on creature type, keywords, etc.
}

