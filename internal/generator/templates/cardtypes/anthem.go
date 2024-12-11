package cardtypes

import (
    "image"
    "path/filepath"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

type AnthemTemplate struct {
    framesPath string
    artBounds  image.Rectangle
}

func NewAnthemTemplate() (*AnthemTemplate, error) {
    return &AnthemTemplate{
        framesPath: filepath.Join("internal", "generator", "templates", "images"),
        artBounds: image.Rect(170, 240, 1330, 1000), // Default art bounds
    }, nil
}

func (t *AnthemTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    framePath := filepath.Join(t.framesPath, "AnthemBase.png")
    return LoadFrame(framePath)
}

func (t *AnthemTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    bounds := layout.GetDefaultBounds()
    // Adjust text bounds specifically for Anthem cards
    bounds.Effect.Bounds = image.Rect(160, 1200, 1340, 1700) // Slightly larger effect box
    return bounds
}

func (t *AnthemTemplate) GetArtBounds() image.Rectangle {
    return t.artBounds
}