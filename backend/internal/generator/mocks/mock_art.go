// internal/generator/mocks/mock_art.go
package mocks

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// MockArtProcessor simulates fetching art with configurable behavior
type MockArtProcessor struct {
	// Simulate fetch delays
	FetchDelay time.Duration
	// Map of card names to error conditions
	ErrorCards map[string]error
	// Simulate network conditions
	ShouldTimeout bool
}

func NewMockArtProcessor() *MockArtProcessor {
	return &MockArtProcessor{
		FetchDelay: 100 * time.Millisecond, // Default small delay
		ErrorCards: make(map[string]error),
	}
}

// SimulateNetworkError adds an error condition for specific cards
func (m *MockArtProcessor) SimulateNetworkError(cardName string, err error) {
	m.ErrorCards[cardName] = err
}

func (m *MockArtProcessor) ProcessArt(data *card.CardDTO, bounds image.Rectangle) (image.Image, error) {
	// Simulate network delay
	time.Sleep(m.FetchDelay)

	// Check for configured errors
	if err, exists := m.ErrorCards[data.Name]; exists {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}

	// Simulate timeout
	if m.ShouldTimeout {
		return nil, fmt.Errorf("failed to fetch image: timeout after 5s")
	}

	// Create a placeholder image with card info
	img := image.NewRGBA(bounds)

	// Fill with a color based on card type
	bgColor := m.getColorForCardType(data.Type)
	draw.Draw(img, bounds, &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw a border
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		img.Set(x, bounds.Min.Y, color.Black)
		img.Set(x, bounds.Max.Y-1, color.Black)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		img.Set(bounds.Min.X, y, color.Black)
		img.Set(bounds.Max.X-1, y, color.Black)
	}

	return img, nil
}

func (m *MockArtProcessor) getColorForCardType(cardType card.CardType) color.Color {
	switch cardType {
	case card.TypeCreature:
		return color.RGBA{150, 200, 150, 255} // Green tint
	case card.TypeSpell:
		return color.RGBA{200, 150, 150, 255} // Red tint
	case card.TypeArtifact:
		return color.RGBA{150, 150, 200, 255} // Blue tint
	default:
		return color.RGBA{200, 200, 200, 255} // Gray
	}
}
