// templates/factory/factory.go
package factory

import (
	"fmt"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/base"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/types"
)

// OutputFormat defines the output format for card generation
type OutputFormat string

const (
	FormatPNG OutputFormat = "png"
	FormatSVG OutputFormat = "svg"
)

// NewTemplate creates the appropriate template type (PNG format for backward compatibility)
func NewTemplate(cardType card.CardType) (base.Template, error) {
	return NewPNGTemplate(cardType)
}

// NewPNGTemplate creates PNG templates (existing logic, now explicitly named)
func NewPNGTemplate(cardType card.CardType) (base.Template, error) {
	var template base.Template
	var err error

	switch cardType {
	case card.TypeCreature:
		template, err = types.NewCreatureTemplate()
	case card.TypeArtifact:
		template, err = types.NewArtifactTemplate()
	case card.TypeSpell:
		template, err = types.NewSpellTemplate()
	case card.TypeIncantation:
		template, err = types.NewIncantationTemplate()
	case card.TypeAnthem:
		template, err = types.NewAnthemTemplate()
	default:
		return nil, fmt.Errorf("unsupported card type: %s", cardType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create PNG template: %w", err)
	}

	return template, nil
}
