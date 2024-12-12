package types

import (
    "image"
    "strings"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type IncantationTemplate struct {
    *templates.BaseTemplate
}

func NewIncantationTemplate() (*IncantationTemplate, error) {
    return &IncantationTemplate{
        BaseTemplate: templates.NewBaseTemplate(),
    }, nil
}

func (t *IncantationTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    if t.isSpecialIncantation(data) {
        return t.LoadFrame("SpecialIncantation.png")
    }
    return t.LoadFrame("BaseIncantation.png")
}

func (t *IncantationTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
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

    // Adjust bounds if it has timing text
    if data.Timing != "" {
        bounds.Effect.Bounds = image.Rect(160, 1300, 1340, 1750)
    }

    return bounds
}

func (t *IncantationTemplate) isSpecialIncantation(data *card.CardData) bool {
    return strings.Contains(strings.ToUpper(data.Effect), "ON ANY CLASH") ||
           strings.Contains(strings.ToUpper(data.Effect), "ON ATTACK")
}