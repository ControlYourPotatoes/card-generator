// templates/base/template.go
package base

import (
    "image"
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

type Template interface {
    GetFrame(data *card.CardData) (image.Image, error)
    GetTextBounds(data *card.CardData) *layout.TextBounds
    GetArtBounds() image.Rectangle
}

// BaseTemplate provides common template functionality
type BaseTemplate struct {
    framesPath string
    artBounds  image.Rectangle
}

func NewBaseTemplate() *BaseTemplate {
    return &BaseTemplate{
        framesPath: getTemplateDir(),
        artBounds:  GetDefaultArtBounds(),
    }
}

// This method must be implemented
func (b *BaseTemplate) GetArtBounds() image.Rectangle {
    return b.artBounds
}