// internal/generator/template/templates/factory.go
package templates

import (
    "fmt"
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/template/base"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/template/types"
)

// NewTemplate creates the appropriate template type
func NewTemplate(cardType card.CardType) (base.Template, error) {
    switch cardType {
    case card.TypeCreature:
        return types.NewCreatureTemplate()
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