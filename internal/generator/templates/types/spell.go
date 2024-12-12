package types

import (
    "image"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

type SpellTemplate struct {
    *templates.BaseTemplate
}

func NewSpellTemplate() (*SpellTemplate, error) {
    return &SpellTemplate{
        BaseTemplate: templates.NewBaseTemplate(),
    }, nil
}

func (t *SpellTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    return t.LoadFrame("BaseSpell.png")
}

func (t *SpellTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    return &layout.TextBounds{
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
}