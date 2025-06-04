package art

import (
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// ArtProcessor defines how card art should be processed
type ArtProcessor interface {
    // ProcessArt handles retrieving and processing art for a card
    ProcessArt(data *card.CardDTO, bounds image.Rectangle) (image.Image, error)
}

// ArtSource represents where art can come from
type ArtSource interface {
    // GetArt retrieves art for a given card
    GetArt(cardID string) (image.Image, error)
}

// ArtProvider is a factory for creating art processors
type ArtProvider interface {
    // GetProcessor returns an appropriate art processor
    GetProcessor() (ArtProcessor, error)
}