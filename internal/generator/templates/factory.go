package templates

import (
    "fmt"
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/templates/types"
)

// NewTemplate creates the appropriate template type
func NewTemplate(cardType card.CardType) (Template, error) {
    // Factory just handles template creation
    switch cardType {
    case card.TypeCreature:
        return c.NewCreatureTemplate()
    case card.TypeArtifact:
        return types.NewArtifactTemplate()
    case card.TypeSpell:
        return types.NewSpellTemplate()
    case card.TypeIncantation:
        return types.NewIncantationTemplate()
    case card.TypeAnthem:
        return types.NewAnthemTemplate()
    default:
        return nil, fmt.Errorf("unsupported card type: %s", cardType)
    }
}