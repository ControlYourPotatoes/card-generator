package templates

import (
    "fmt"
    "image"
    "image/png"
    "os"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/templates/cardtypes"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

// Card dimensions from CardConjurer template
const (
    CardWidth  = 1500
    CardHeight = 2100
    MarginX    = 66
    MarginY    = 60
)

// Template interface defines what each card template must provide
type Template interface {
    // GetFrame returns the appropriate frame image for the card
    GetFrame(data *card.CardData) (image.Image, error)

    // GetTextBounds returns text positioning for the card
    GetTextBounds(data *card.CardData) *layout.TextBounds

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

// Default art bounds used by most cards
func GetDefaultArtBounds() image.Rectangle {
    return image.Rect(170, 240, 1330, 1000)
}

// NewTemplate creates appropriate template for card type
func NewTemplate(cardType card.CardType) (Template, error) {
    switch cardType {
    case card.TypeCreature:
        return cardtypes.NewCreatureTemplate()
    case card.TypeArtifact:
        return cardtypes.NewArtifactTemplate()
    case card.TypeSpell:
        return cardtypes.NewSpellTemplate()
    case card.TypeIncantation:
        return cardtypes.NewIncantationTemplate()
    case card.TypeAnthem:
        return cardtypes.NewAnthemTemplate()
    default:
        return nil, fmt.Errorf("unsupported card type: %s", cardType)
    }
}