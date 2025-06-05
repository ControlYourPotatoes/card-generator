// svg/generator.go
package svg

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/svg/metadata"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/svg/templates"
)

// svgGenerator implements the SVGGenerator interface
type svgGenerator struct {
	templateLoader  SVGTemplateLoader
	textRenderer    SVGTextRenderer
	metadataBuilder MetadataBuilder
	templateDir     string
}

// NewSVGGenerator creates a new SVG generator instance
func NewSVGGenerator() (SVGGenerator, error) {
	return NewSVGGeneratorWithConfig("./templates/svg")
}

// NewSVGGeneratorWithConfig creates a new SVG generator with custom template directory
func NewSVGGeneratorWithConfig(templateDir string) (SVGGenerator, error) {
	return &svgGenerator{
		templateDir: templateDir,
		// TODO: Initialize other components in later phases
	}, nil
}

// GenerateCard implements the CardGenerator interface for backward compatibility
func (g *svgGenerator) GenerateCard(data *card.CardDTO, outputPath string) error {
	return g.GenerateSVG(data, outputPath)
}

// ValidateCard implements the CardGenerator interface
func (g *svgGenerator) ValidateCard(data *card.CardDTO) error {
	if data == nil {
		return fmt.Errorf("card data cannot be nil")
	}
	if data.Name == "" {
		return fmt.Errorf("card name cannot be empty")
	}
	if data.Type == "" {
		return fmt.Errorf("card type cannot be empty")
	}
	return nil
}

// Close implements the CardGenerator interface
func (g *svgGenerator) Close() error {
	// TODO: Cleanup resources if needed
	return nil
}

// GenerateSVG generates an SVG card file
func (g *svgGenerator) GenerateSVG(data *card.CardDTO, outputPath string) error {
	// Validate card data first
	if err := g.ValidateCard(data); err != nil {
		return fmt.Errorf("invalid card data: %w", err)
	}

	// Get SVG template for card type
	svgTemplate, err := templates.NewSVGTemplate(data.Type, g.templateDir)
	if err != nil {
		return fmt.Errorf("failed to load SVG template: %w", err)
	}

	// Get SVG template content
	templateContent := svgTemplate.GetSVGTemplate()

	// Process template with card data
	processedSVG, err := g.processTemplate(templateContent, data)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write SVG to file
	if err := os.WriteFile(outputPath, []byte(processedSVG), 0644); err != nil {
		return fmt.Errorf("failed to write SVG file: %w", err)
	}

	return nil
}

// GenerateWithMetadata generates SVG with game metadata
func (g *svgGenerator) GenerateWithMetadata(data *card.CardDTO, meta metadata.SVGMetadata) (string, error) {
	// TODO: Implement in Phase 3
	return "", fmt.Errorf("metadata generation not yet implemented")
}

// processTemplate processes the SVG template with card data
func (g *svgGenerator) processTemplate(templateContent string, data *card.CardDTO) (string, error) {
	// Create template
	tmpl, err := template.New("card").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	templateData := g.prepareTemplateData(data)

	// Execute template
	var result strings.Builder
	if err := tmpl.Execute(&result, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// prepareTemplateData prepares card data for template execution
func (g *svgGenerator) prepareTemplateData(data *card.CardDTO) map[string]interface{} {
	templateData := map[string]interface{}{
		"Name":    data.Name,
		"Cost":    data.Cost,
		"Effect":  data.Effect,
		"Attack":  data.Attack,
		"Defense": data.Defense,
		"Trait":   data.Trait,
		"Type":    string(data.Type),
	}

	// Add formatted type information
	if data.Type == card.TypeCreature {
		if data.Trait != "" {
			templateData["TypeLine"] = fmt.Sprintf("Creature - %s", data.Trait)
		} else {
			templateData["TypeLine"] = "Creature"
		}
	} else {
		templateData["TypeLine"] = string(data.Type)
	}

	// Format effect text with proper line breaks
	if data.Effect != "" {
		// Basic text formatting - replace newlines with HTML breaks for foreignObject
		formattedEffect := strings.ReplaceAll(data.Effect, "\n", "<br/>")
		templateData["FormattedEffect"] = formattedEffect
	}

	return templateData
} 