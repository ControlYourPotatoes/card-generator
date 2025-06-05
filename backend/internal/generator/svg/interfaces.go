// svg/interfaces.go
package svg

import (
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/svg/metadata"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/base"
)

// SVGTemplate extends the base Template interface with SVG-specific capabilities
type SVGTemplate interface {
	base.Template // Inherit existing interface for compatibility
	GetSVGTemplate() string
	GetInteractiveZones() map[string]metadata.InteractiveZone
	GetAnimationTargets() []metadata.AnimationTarget
}

// SVGGenerator extends the existing CardGenerator interface with SVG-specific generation
type SVGGenerator interface {
	generator.CardGenerator // Inherit existing interface for compatibility
	GenerateSVG(data *card.CardDTO, outputPath string) error
	GenerateWithMetadata(data *card.CardDTO, metadata metadata.SVGMetadata) (string, error)
}

// SVGTemplateLoader handles loading and parsing of SVG templates
type SVGTemplateLoader interface {
	LoadTemplate(cardType card.CardType) (SVGTemplate, error)
	GetTemplateNames() []string
	ValidateTemplate(template string) error
}

// SVGTextRenderer handles text rendering within SVG elements
type SVGTextRenderer interface {
	RenderTextToSVG(data *card.CardDTO, textBounds map[string]image.Rectangle) (string, error)
	GetTextDimensions(text string, fontSize int) (width, height int)
	FormatTextForSVG(text string) string
}

// MetadataBuilder creates game-ready metadata for SVG cards
type MetadataBuilder interface {
	BuildMetadata(data *card.CardDTO, interactiveZones map[string]metadata.InteractiveZone) metadata.SVGMetadata
	GetDefaultZones(cardType card.CardType) map[string]metadata.InteractiveZone
	ValidateMetadata(metadata metadata.SVGMetadata) error
} 