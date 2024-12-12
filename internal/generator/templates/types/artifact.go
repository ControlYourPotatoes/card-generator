package types

import (
    "image"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type ArtifactTemplate struct {
    *templates.BaseTemplate
}

func NewArtifactTemplate() (*ArtifactTemplate, error) {
    return &ArtifactTemplate{
        BaseTemplate: templates.NewBaseTemplate(),
    }, nil
}

func (t *ArtifactTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    return t.LoadFrame("BaseArtifact.png")
}

func (t *ArtifactTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    bounds := &layout.TextBounds{
        Name: layout.TextConfig{
            Bounds:    image.Rect(125, 90, 1375, 170),
            FontSize:  72,
            Alignment: "center",
        },
        Effect: layout.TextConfig{
            Bounds:    image.Rect(160, 1250, 1340, 1750),
            FontSize:  48,
            Alignment: "left",
        },
    }

    // Adjust bounds for equipment artifacts
    if data.IsEquipment {
        bounds.Effect.Bounds = image.Rect(160, 1250, 1340, 1750)
    }

    return bounds
}