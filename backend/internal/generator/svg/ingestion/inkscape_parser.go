// Package ingestion handles parsing and processing of raw Inkscape SVG files
// into structured template components for the card generator.
//
// This package implements the ingestion pipeline as outlined in Phase 3a:
// - Parse raw Inkscape SVG files by naming convention
// - Extract objects (visual elements) and boundaries (text areas) separately
// - Create transparency positioning metadata
// - Generate clean, structured template data
package ingestion

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// InkscapeParser handles parsing of Inkscape SVG files into structured components
type InkscapeParser struct {
	layerExtractor   *LayerExtractor
	objectDetector   *ObjectDetector
	boundaryFinder   *BoundaryFinder
	metadataBuilder  *MetadataBuilder
}

// NewInkscapeParser creates a new parser with all required components
func NewInkscapeParser() *InkscapeParser {
	return &InkscapeParser{
		layerExtractor:  NewLayerExtractor(),
		objectDetector:  NewObjectDetector(),
		boundaryFinder:  NewBoundaryFinder(),
		metadataBuilder: NewMetadataBuilder(),
	}
}

// ParsedTemplate represents the structured result of parsing an Inkscape SVG file
type ParsedTemplate struct {
	Objects     map[ObjectType]*CardObject       // Visual elements (frames, backgrounds, etc.)
	Boundaries  map[BoundaryType]*TextBoundary   // Text rendering areas only
	Positioning *TransparencyPositioning         // Opacity-based positioning system
	Metadata    *TemplateMetadata                // Template metadata and validation info
}

// ParseSVGFile processes an Inkscape SVG file and returns structured template data
func (p *InkscapeParser) ParseSVGFile(filepath string) (*ParsedTemplate, error) {
	// Open and parse SVG file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SVG file %s: %w", filepath, err)
	}
	defer file.Close()

	// Parse SVG DOM
	svgDoc, err := p.parseSVGDOM(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG DOM: %w", err)
	}

	// Extract layers by naming convention
	layers, err := p.layerExtractor.ExtractLayers(svgDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract layers: %w", err)
	}

	// Detect objects vs boundaries by Inkscape naming convention
	objects, err := p.objectDetector.DetectObjects(layers)
	if err != nil {
		return nil, fmt.Errorf("failed to detect objects: %w", err)
	}

	boundaries, err := p.boundaryFinder.DetectBoundaries(layers)
	if err != nil {
		return nil, fmt.Errorf("failed to detect boundaries: %w", err)
	}

	// Build transparency positioning metadata
	positioning, err := p.buildTransparencyPositioning(layers, objects, boundaries)
	if err != nil {
		return nil, fmt.Errorf("failed to build positioning: %w", err)
	}

	// Generate template metadata
	metadata := p.metadataBuilder.BuildMetadata(filepath, objects, boundaries, positioning)

	return &ParsedTemplate{
		Objects:     objects,
		Boundaries:  boundaries,
		Positioning: positioning,
		Metadata:    metadata,
	}, nil
}

// ParseSVGContent processes SVG content directly (for testing or in-memory processing)
func (p *InkscapeParser) ParseSVGContent(content string) (*ParsedTemplate, error) {
	reader := strings.NewReader(content)
	
	// Parse SVG DOM
	svgDoc, err := p.parseSVGDOM(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SVG content: %w", err)
	}

	// Extract layers by naming convention
	layers, err := p.layerExtractor.ExtractLayers(svgDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract layers: %w", err)
	}

	// Process similar to ParseSVGFile
	objects, err := p.objectDetector.DetectObjects(layers)
	if err != nil {
		return nil, fmt.Errorf("failed to detect objects: %w", err)
	}

	boundaries, err := p.boundaryFinder.DetectBoundaries(layers)
	if err != nil {
		return nil, fmt.Errorf("failed to detect boundaries: %w", err)
	}

	positioning, err := p.buildTransparencyPositioning(layers, objects, boundaries)
	if err != nil {
		return nil, fmt.Errorf("failed to build positioning: %w", err)
	}

	metadata := p.metadataBuilder.BuildMetadata("in-memory", objects, boundaries, positioning)

	return &ParsedTemplate{
		Objects:     objects,
		Boundaries:  boundaries,
		Positioning: positioning,
		Metadata:    metadata,
	}, nil
}

// parseSVGDOM parses the SVG XML structure from a reader
func (p *InkscapeParser) parseSVGDOM(reader io.Reader) (*SVGDocument, error) {
	var svgDoc SVGDocument
	
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&svgDoc); err != nil {
		return nil, fmt.Errorf("failed to decode SVG XML: %w", err)
	}

	return &svgDoc, nil
}

// buildTransparencyPositioning creates opacity-based positioning from layers and objects
func (p *InkscapeParser) buildTransparencyPositioning(
	layers map[string]*InkscapeLayer,
	objects map[ObjectType]*CardObject,
	boundaries map[BoundaryType]*TextBoundary,
) (*TransparencyPositioning, error) {
	
	positioning := &TransparencyPositioning{
		Layers:      make(map[string]*TransparencyLayer),
		OpacityMaps: make(map[string]*OpacityMap),
		BlendModes:  make(map[string]BlendMode),
	}

	// Process each layer for transparency positioning
	for layerID, layer := range layers {
		// Create transparency layer
		transparencyLayer := &TransparencyLayer{
			ID:         layerID,
			ObjectType: layer.ObjectType,
			ZIndex:     layer.ZIndex,
			BlendMode:  layer.BlendMode,
		}

		// Generate opacity map for positioning
		opacityMap := &OpacityMap{
			PrimaryZone:      layer.VisibleArea,
			FadeZones:        layer.FadeZones,
			FullyVisible:     []Rectangle{layer.VisibleArea},
			PartiallyVisible: make([]GradientZone, 0),
			FullyHidden:      make([]Rectangle, 0),
		}

		positioning.Layers[layerID] = transparencyLayer
		positioning.OpacityMaps[layerID] = opacityMap
		positioning.BlendModes[layerID] = layer.BlendMode
	}

	return positioning, nil
}

// ValidateCardTypeCompatibility checks if parsed template is compatible with card type
func (p *InkscapeParser) ValidateCardTypeCompatibility(template *ParsedTemplate, cardType card.CardType) error {
	// Check if template has required objects for card type
	requiredObjects := GetRequiredObjectsForCardType(cardType)
	for _, required := range requiredObjects {
		if _, exists := template.Objects[required]; !exists {
			return fmt.Errorf("template missing required object %s for card type %s", required, cardType)
		}
	}

	// Check if template has required boundaries for card type
	requiredBoundaries := GetRequiredBoundariesForCardType(cardType)
	for _, required := range requiredBoundaries {
		if _, exists := template.Boundaries[required]; !exists {
			return fmt.Errorf("template missing required boundary %s for card type %s", required, cardType)
		}
	}

	return nil
}

// GetSupportedCardTypes returns the card types this template can support
func (p *InkscapeParser) GetSupportedCardTypes(template *ParsedTemplate) []card.CardType {
	var supportedTypes []card.CardType

	// Check each card type for compatibility
	for _, cardType := range []card.CardType{
		card.TypeCreature,
		card.TypeArtifact, 
		card.TypeSpell,
		card.TypeAnthem,
	} {
		if p.ValidateCardTypeCompatibility(template, cardType) == nil {
			supportedTypes = append(supportedTypes, cardType)
		}
	}

	return supportedTypes
} 