package cardtypes

import (
    "image"
    "path/filepath"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

type IncantationTemplate struct {
    framesPath string
    artBounds  image.Rectangle
}

func NewIncantationTemplate() (*IncantationTemplate, error) {
    return &IncantationTemplate{
        framesPath: filepath.Join("internal", "generator", "templates", "images"),
        artBounds: image.Rect(170, 240, 1330, 1000), // Default art bounds
    }, nil
}

func (t *IncantationTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    framePath := filepath.Join(t.framesPath, "BaseIncantation.png")
    return LoadFrame(framePath)
}

func (t *IncantationTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    return layout.GetDefaultBounds()
}

func (t *IncantationTemplate) GetArtBounds() image.Rectangle {
    return t.artBounds
}

func (t *IncantationTemplate) isSpecialFrame(data *card.CardData) bool {
    // Logic to determine if card should use special frame
    return false
}