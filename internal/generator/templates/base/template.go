// internal/generator/template/base/template.go
package base

import (
    "image"
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

// Template defines what each card template must provide
type Template interface {
    GetFrame(data *card.CardData) (image.Image, error)
    GetTextBounds(data *card.CardData) *layout.TextBounds
    GetArtBounds() image.Rectangle
}