// File: internal/generator/templates/cardtypes/cardtypes.go
package cardtypes

import (
    "image"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

// CardTemplate defines the interface that all card type templates must implement
type CardTemplate interface {
    // GetFrame returns the appropriate frame image for the card
    GetFrame(data *card.CardData) (image.Image, error)

    // GetTextBounds returns text positioning for the card
    GetTextBounds(data *card.CardData) *layout.TextBounds

    // GetArtBounds returns where card art should be placed
    GetArtBounds() image.Rectangle
}

// BaseCardTemplate provides common functionality for all card templates
type BaseCardTemplate struct {
    FramesPath string
    ArtBounds  image.Rectangle
}

func (b *BaseCardTemplate) GetArtBounds() image.Rectangle {
    return b.ArtBounds
}