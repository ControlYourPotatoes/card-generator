package cardtypes

import (
	"image"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
	"github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type IncantationTemplate struct {
	templates.BaseTemplate
}

func NewIncantationTemplate() (*IncantationTemplate, error) {
	return &IncantationTemplate{
		BaseTemplate: templates.BaseTemplate{
			framesPath: filepath.Join("internal", "generator", "templates", "images"),
			artBounds:  templates.GetDefaultArtBounds(),
		},
	}, nil
}

func (t *IncantationTemplate) GetFrame(data *card.CardData) (image.Image, error) {
	return templates.LoadFrame(filepath.Join(t.framesPath, "BaseIncantation.png"))
}

func (t *IncantationTemplate) GetTextBounds(data *card.CardData) *templates.TextBounds {
	return templates.GetDefaultBounds()
}

func (t *IncantationTemplate) GetArtBounds() image.Rectangle {
	return t.artBounds
}

func (t *IncantationTemplate) isSpecialFrame(data *card.CardData) bool {
	// Logic to determine if card should use special frame
	// Could be based on creature type, keywords, etc.
	return false
}
