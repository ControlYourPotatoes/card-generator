package templates

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
)

// Card dimensions from CardConjurer template
const (
	CardWidth  = 1500
	CardHeight = 2100
	MarginX    = 66
	MarginY    = 60
)

// TextBounds defines positions for text elements
type TextBounds struct {
	Name          TextConfig
	ManaCost      TextConfig
	Type          TextConfig
	Effect        TextConfig
	Stats         *StatsConfig // Only for creatures
	CollectorInfo TextConfig
}

// TextConfig holds positioning and styling data
type TextConfig struct {
	Bounds    image.Rectangle
	FontSize  float64
	Alignment string // "left", "center", "right"
}

// StatsConfig holds creature stats positioning
type StatsConfig struct {
	Left  TextConfig // Power/Attack
	Right TextConfig // Defense/Toughness
}

// Template interface defines what each card template must provide
type Template interface {
	// GetFrame returns the appropriate frame image for the card
	GetFrame(data *card.CardData) (image.Image, error)

	// GetTextBounds returns text positioning for the card
	GetTextBounds(data *card.CardData) *TextBounds

	// GetArtBounds returns where card art should be placed
	GetArtBounds() image.Rectangle
}

// BaseTemplate provides common template functionality
type BaseTemplate struct {
	framesPath string
	artBounds  image.Rectangle
}

// Load a frame image from the templates directory
func LoadFrame(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open frame: %w", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode frame: %w", err)
	}

	return img, nil
}

// Get default text bounds used by most cards
func GetDefaultBounds() *TextBounds {
	return &TextBounds{
		Name: TextConfig{
			Bounds:    image.Rect(125, 90, 1375, 170),
			FontSize:  80,
			Alignment: "left",
		},
		ManaCost: TextConfig{
			Bounds:    image.Rect(125, 90, 1375, 170),
			FontSize:  80,
			Alignment: "right",
		},
		Type: TextConfig{
			Bounds:    image.Rect(125, 1885, 1375, 1955),
			FontSize:  70,
			Alignment: "center",
		},
		Effect: TextConfig{
			Bounds:    image.Rect(160, 1300, 1340, 1800),
			FontSize:  60,
			Alignment: "left",
		},
		CollectorInfo: TextConfig{
			Bounds:    image.Rect(110, 2010, 750, 2090),
			FontSize:  32,
			Alignment: "left",
		},
	}
}

// Default art bounds used by most cards
func GetDefaultArtBounds() image.Rectangle {
	return image.Rect(170, 240, 1330, 1000)
}

// NewTemplate creates appropriate template for card type
func NewTemplate(cardType card.CardType) (Template, error) {
	switch cardType {
	case card.TypeCreature:
		return NewCreatureTemplate()
	case card.TypeArtifact:
		return NewArtifactTemplate()
	case card.TypeSpell:
		return NewSpellTemplate()
	case card.TypeIncantation:
		return NewIncantationTemplate()
	case card.TypeAnthem:
		return NewAnthemTemplate()
	default:
		return nil, fmt.Errorf("unsupported card type: %s", cardType)
	}
}
