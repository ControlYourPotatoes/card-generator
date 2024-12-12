package types

import (
    "image"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/templates/base"
)

type AnthemTemplate struct {
    *base.BaseTemplate
}

func NewAnthemTemplate() (*AnthemTemplate, error) {
    return &AnthemTemplate{
        BaseTemplate: base.NewBaseTemplate(),
    }, nil
}

func (t *AnthemTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    return t.LoadFrame("BaseAnthem.png")
}

func (t *AnthemTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    return &layout.TextBounds{
        Name: layout.TextConfig{
            Bounds:    image.Rect(125, 90, 1375, 170),
            FontSize:  72,
            Alignment: "center",
        },
        Effect: layout.TextConfig{
            Bounds:    image.Rect(160, 1200, 1340, 1700),
            FontSize:  48,
            Alignment: "left",
        },
    }
}