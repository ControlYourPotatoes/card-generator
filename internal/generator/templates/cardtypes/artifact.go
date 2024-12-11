package cardtypes

import (
    "image"
    "path/filepath"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

type ArtifactTemplate struct {
    framesPath string
    artBounds  image.Rectangle
}

func NewArtifactTemplate() (*ArtifactTemplate, error) {
    return &ArtifactTemplate{
        framesPath: filepath.Join("internal", "generator", "templates", "images"),
        artBounds: image.Rect(170, 240, 1330, 1000), // Default art bounds
    }, nil
}

func (t *ArtifactTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    framePath := filepath.Join(t.framesPath, "BaseArtifact.png")
    return LoadFrame(framePath)
}

func (t *ArtifactTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    bounds := layout.GetDefaultBounds()
    // Adjust text bounds specifically for Artifact cards
    if data.IsEquipment {
        bounds.Effect.Bounds = image.Rect(160, 1250, 1340, 1750) // Adjusted for equipment text
    }
    return bounds
}

func (t *ArtifactTemplate) GetArtBounds() image.Rectangle {
    return t.artBounds
}