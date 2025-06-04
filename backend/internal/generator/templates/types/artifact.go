package types

import (
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/base"
)

type ArtifactTemplate struct {
	*base.BaseTemplate
}

func NewArtifactTemplate() (*ArtifactTemplate, error) {
	return &ArtifactTemplate{
		BaseTemplate: base.NewBaseTemplate(),
	}, nil
}

func (t *ArtifactTemplate) GetFrame(data *card.CardDTO) (image.Image, error) {
	var frameName string
	if data.IsEquipment {
		frameName = "BaseArtifactEquipment.png"
	} else {
		frameName = "BaseArtifact.png"
	}
	return t.LoadFrame(frameName)
}

func (t *ArtifactTemplate) GetTextBounds(data *card.CardDTO) map[string]image.Rectangle {
	bounds := make(map[string]image.Rectangle)
	bounds["name"] = image.Rect(125, 90, 1375, 170)
	bounds["cost"] = image.Rect(125, 90, 1375, 170)
	bounds["type"] = image.Rect(125, 1885, 1375, 1955)
	bounds["effect"] = image.Rect(160, 1250, 1340, 1750)
	bounds["collector"] = image.Rect(110, 2010, 750, 2090)
	
	// Equipment artifacts might need special positioning
	if data.IsEquipment {
		bounds["effect"] = image.Rect(160, 1200, 1340, 1700)
	}
	
	return bounds
}