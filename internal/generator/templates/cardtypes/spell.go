package cardtypes

import (
    "image"
    "path/filepath"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

type SpellTemplate struct {
    framesPath string
    artBounds  image.Rectangle
}

func NewSpellTemplate() (*SpellTemplate, error) {
    return &SpellTemplate{
        framesPath: filepath.Join("internal", "generator", "templates", "images"),
        artBounds: image.Rect(170, 240, 1330, 1000), // Default art bounds
    }, nil
}

func (t *SpellTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    framePath := filepath.Join(t.framesPath, "BaseSpell.png")
    return LoadFrame(framePath)
}

func (t *SpellTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    return layout.GetDefaultBounds()
}

func (t *SpellTemplate) GetArtBounds() image.Rectangle {
    return t.artBounds
}

func (t *SpellTemplate) isSpecialFrame(data *card.CardData) bool {
    // Logic to determine if card should use special frame
    return false
}