package types

import (
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/base"
)

type IncantationTemplate struct {
    *base.BaseTemplate
}

func NewIncantationTemplate() (*IncantationTemplate, error) {
    return &IncantationTemplate{
        BaseTemplate: base.NewBaseTemplate(),
    }, nil
}

func (t *IncantationTemplate) GetFrame(data *card.CardDTO) (image.Image, error) {
    return t.LoadFrame("BaseIncantation.png")
}

func (t *IncantationTemplate) GetTextBounds(data *card.CardDTO) map[string]image.Rectangle {
    bounds := make(map[string]image.Rectangle)
    bounds["name"] = image.Rect(125, 90, 1375, 170)
    bounds["cost"] = image.Rect(125, 90, 1375, 170)
    bounds["type"] = image.Rect(125, 1885, 1375, 1955)
    bounds["effect"] = image.Rect(160, 1250, 1340, 1750)
    bounds["collector"] = image.Rect(110, 2010, 750, 2090)
    
    // Incantations might have different positioning
    if data.Continuous {
        bounds["effect"] = image.Rect(160, 1200, 1340, 1700)
    }
    
    return bounds
} 