package art

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// placeholderProcessor implements ArtProcessor using placeholder images
type placeholderProcessor struct {
	client *http.Client
}

// NewPlaceholderProcessor creates a new placeholder art processor
func NewPlaceholderProcessor() ArtProcessor {
	return &placeholderProcessor{
		client: &http.Client{},
	}
}

// ProcessArt implements ArtProcessor interface
func (p *placeholderProcessor) ProcessArt(data *card.CardDTO, bounds image.Rectangle) (image.Image, error) {
	width := bounds.Dx()
	height := bounds.Dy()

	// Create URL for placeholder image
	url := fmt.Sprintf("https://placeholderimage.dev/%dx%d", width, height)

	// Fetch the placeholder image
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch placeholder image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image, status: %d", resp.StatusCode)
	}

	// Decode the image
	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Create a new image with the correct bounds
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, image.Point{}, draw.Over)

	return dst, nil
}
