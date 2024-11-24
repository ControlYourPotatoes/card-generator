package templates

import (
    "image"
    

    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

// Template defines the interface for all card templates
type Template interface {
    // Apply applies the template to the given image
    Apply(img *image.RGBA) error
    
    // GetArtBounds returns the boundaries where card art should be placed
    GetArtBounds() image.Rectangle
    
    // GetTextBounds returns the boundaries for different text elements
    GetTextBounds() TextBounds
}

// TextBounds holds the positions for all text elements
type TextBounds struct {
    Name   image.Rectangle
    Effect image.Rectangle
    Stats  StatsPosition // Only used by creatures
}

type StatsPosition struct {
    Attack  image.Rectangle
    Defense image.Rectangle
}

// BaseTemplate contains common template functionality
type BaseTemplate struct {
    frame    image.Image
    artFrame image.Rectangle
    textPos  TextBounds
}

func LoadTemplateForType(cardType card.CardType) (Template, error) {
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