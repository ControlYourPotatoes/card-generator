// svg/templates/factory.go
package templates

import (
	"fmt"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/svg/metadata"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/base"
)

// SVGTemplate interface defined locally to avoid import cycle with svg package
type SVGTemplate interface {
	base.Template // Inherit existing interface for compatibility
	GetSVGTemplate() string
	GetInteractiveZones() map[string]metadata.InteractiveZone
	GetAnimationTargets() []metadata.AnimationTarget
}

// Factory creates SVG templates based on card type
type Factory struct {
	templateDir string
	loader      *Loader
}

// NewFactory creates a new SVG template factory
func NewFactory(templateDir string) *Factory {
	return &Factory{
		templateDir: templateDir,
		loader:      NewLoader(templateDir),
	}
}

// CreateTemplate creates an SVG template for the specified card type
func (f *Factory) CreateTemplate(cardType card.CardType) (SVGTemplate, error) {
	switch cardType {
	case card.TypeCreature:
		return f.createCreatureTemplate()
	case card.TypeArtifact:
		return f.createArtifactTemplate()
	case card.TypeSpell:
		return f.createSpellTemplate()
	case card.TypeIncantation:
		return f.createIncantationTemplate()
	case card.TypeAnthem:
		return f.createAnthemTemplate()
	default:
		return nil, fmt.Errorf("unsupported card type: %v", cardType)
	}
}

// createCreatureTemplate creates the creature SVG template (Phase 2 implementation)
func (f *Factory) createCreatureTemplate() (SVGTemplate, error) {
	template := &CreatureTemplate{
		templateDir: f.templateDir,
		loader:      f.loader,
	}
	return template, nil
}

// Template creation methods (to be implemented in Phase 3)
func (f *Factory) createArtifactTemplate() (SVGTemplate, error) {
	// TODO: Implement in Phase 3
	return nil, fmt.Errorf("artifact template not yet implemented")
}

func (f *Factory) createSpellTemplate() (SVGTemplate, error) {
	// TODO: Implement in Phase 3
	return nil, fmt.Errorf("spell template not yet implemented")
}

func (f *Factory) createIncantationTemplate() (SVGTemplate, error) {
	// TODO: Implement in Phase 3
	return nil, fmt.Errorf("incantation template not yet implemented")
}

func (f *Factory) createAnthemTemplate() (SVGTemplate, error) {
	// TODO: Implement in Phase 3
	return nil, fmt.Errorf("anthem template not yet implemented")
}

// NewSVGTemplate is a public function that can be called from outside packages to avoid import cycles
func NewSVGTemplate(cardType card.CardType, templateDir string) (SVGTemplate, error) {
	factory := NewFactory(templateDir)
	return factory.CreateTemplate(cardType)
} 