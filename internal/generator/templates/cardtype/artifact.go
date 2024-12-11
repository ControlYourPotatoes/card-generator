package cardtypes

import (
	"image"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
	"github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type ArtifactTemplate struct {
	templates.BaseTemplate
}

func NewArtifactTemplate() (*ArtifactTemplate, error) {
	return &ArtifactTemplate{
		BaseTemplate: templates.BaseTemplate{
			framesPath: filepath.Join("internal", "generator", "templates", "images"),
			artBounds:  templates.GetDefaultArtBounds(),
		},
	}, nil
}

func (t *ArtifactTemplate) GetFrame(data *card.CardData) (image.Image, error) {
	return templates.LoadFrame(filepath.Join(t.framesPath, "BaseArtifact.png"))
}

func (t *ArtifactTemplate) GetTextBounds(data *card.CardData) *templates.TextBounds {
	return templates.GetDefaultBounds()
}

func (t *ArtifactTemplate) GetArtBounds() image.Rectangle {
	return t.artBounds
}

func (t *ArtifactTemplate) isSpecialFrame(data *card.CardData) bool {
	// Logic to determine if card should use special frame
	// Could be based on creature type, keywords, etc.
	return false
}
