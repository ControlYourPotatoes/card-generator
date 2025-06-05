// svg/templates/loader.go
package templates

import (
	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// Loader handles loading SVG templates from various sources
type Loader struct {
	templateDir string
	cache       map[card.CardType]SVGTemplate
}

// NewLoader creates a new SVG template loader
func NewLoader(templateDir string) *Loader {
	return &Loader{
		templateDir: templateDir,
		cache:       make(map[card.CardType]SVGTemplate),
	}
}

// LoadTemplate loads an SVG template for the specified card type
func (l *Loader) LoadTemplate(cardType card.CardType) (SVGTemplate, error) {
	// TODO: Implement template loading in Phase 2
	return nil, nil
}

// GetTemplateNames returns a list of available template names
func (l *Loader) GetTemplateNames() []string {
	// TODO: Implement in Phase 2
	return []string{}
}

// ValidateTemplate checks if a template string is valid SVG
func (l *Loader) ValidateTemplate(template string) error {
	// TODO: Implement SVG validation in Phase 2
	return nil
}

// ClearCache clears the template cache
func (l *Loader) ClearCache() {
	l.cache = make(map[card.CardType]SVGTemplate)
} 