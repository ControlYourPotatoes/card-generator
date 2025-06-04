package text

import (
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// TextProcessor defines how card text should be rendered
type TextProcessor interface {
	// RenderText renders text onto the given image within the specified bounds
	RenderText(img *image.RGBA, data *card.CardDTO, bounds map[string]image.Rectangle) error
}

// basicTextProcessor implements TextProcessor with basic text rendering
type basicTextProcessor struct{}

// NewTextProcessor creates a new basic text processor
func NewTextProcessor() (TextProcessor, error) {
	return &basicTextProcessor{}, nil
}

// RenderText implements TextProcessor interface
func (t *basicTextProcessor) RenderText(img *image.RGBA, data *card.CardDTO, bounds map[string]image.Rectangle) error {
	// For now, this is a placeholder implementation
	// In a real implementation, this would render the card's text fields
	// using a font library like freetype or golang.org/x/image/font
	return nil
}
