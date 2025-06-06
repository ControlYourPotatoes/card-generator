// layer_extractor.go extracts layers from SVG documents
package ingestion

import (
	"fmt"
	"strconv"
	"strings"
)

// LayerExtractor extracts layers from SVG documents
type LayerExtractor struct {
	// Configuration for layer extraction
}

// NewLayerExtractor creates a new layer extractor
func NewLayerExtractor() *LayerExtractor {
	return &LayerExtractor{}
}

// ExtractLayers processes an SVG document and extracts layer information
func (le *LayerExtractor) ExtractLayers(svgDoc *SVGDocument) (map[string]*InkscapeLayer, error) {
	layers := make(map[string]*InkscapeLayer)
	
	// Process each group (Inkscape layers are typically SVG groups)
	for i, group := range svgDoc.Groups {
		layer := le.extractLayerFromGroup(&group, i)
		if layer != nil {
			layers[layer.ID] = layer
		}
	}
	
	return layers, nil
}

// extractLayerFromGroup converts an SVG group to an Inkscape layer
func (le *LayerExtractor) extractLayerFromGroup(group *SVGGroup, index int) *InkscapeLayer {
	// Create layer with basic information
	layer := &InkscapeLayer{
		ID:     group.ID,
		Label:  group.Label,
		ZIndex: index, // Use group order as initial Z-index
	}
	
	// If no ID, generate one
	if layer.ID == "" {
		layer.ID = fmt.Sprintf("layer-%d", index)
	}
	
	// If no label, use ID
	if layer.Label == "" {
		layer.Label = layer.ID
	}
	
	// Parse visible area from group style or default to full canvas
	layer.VisibleArea = le.parseVisibleAreaFromGroup(group)
	
	// Parse blend mode from style
	layer.BlendMode = le.parseBlendModeFromGroup(group)
	
	// Extract content (simplified - just store the group info for now)
	layer.Content = le.extractContentFromGroup(group)
	
	// Initialize fade zones (empty for now, to be populated by transparency analysis)
	layer.FadeZones = make([]FadeZone, 0)
	
	return layer
}

// parseVisibleAreaFromGroup extracts the visible area from a group
func (le *LayerExtractor) parseVisibleAreaFromGroup(group *SVGGroup) Rectangle {
	// For now, return default card size (1500x2100)
	// In a real implementation, this would analyze the group's bounding box
	return Rectangle{
		X:      0,
		Y:      0,
		Width:  1500,
		Height: 2100,
	}
}

// parseBlendModeFromGroup extracts blend mode from group style
func (le *LayerExtractor) parseBlendModeFromGroup(group *SVGGroup) BlendMode {
	// Parse style attribute for blend mode
	style := strings.ToLower(group.Style)
	
	if strings.Contains(style, "mix-blend-mode:multiply") {
		return BlendModeMultiply
	} else if strings.Contains(style, "mix-blend-mode:screen") {
		return BlendModeScreen
	} else if strings.Contains(style, "mix-blend-mode:overlay") {
		return BlendModeOverlay
	}
	
	// Default to normal blending
	return BlendModeNormal
}

// extractContentFromGroup extracts SVG content from a group
func (le *LayerExtractor) extractContentFromGroup(group *SVGGroup) string {
	// For now, return a placeholder
	// In a real implementation, this would serialize the group's elements
	return fmt.Sprintf(`<g id="%s" style="%s"><!-- Layer content --></g>`, group.ID, group.Style)
}

// AnalyzeLayerHierarchy analyzes the hierarchical structure of layers
func (le *LayerExtractor) AnalyzeLayerHierarchy(layers map[string]*InkscapeLayer) error {
	// Sort layers by Z-index for proper ordering
	le.sortLayersByZIndex(layers)
	
	// Analyze dependencies between layers
	return le.analyzeDependencies(layers)
}

// sortLayersByZIndex adjusts Z-index values based on layer hierarchy
func (le *LayerExtractor) sortLayersByZIndex(layers map[string]*InkscapeLayer) {
	// Create a slice for sorting
	layerSlice := make([]*InkscapeLayer, 0, len(layers))
	for _, layer := range layers {
		layerSlice = append(layerSlice, layer)
	}
	
	// Sort by current Z-index (already set from group order)
	// Update Z-index values to ensure proper spacing
	for i, layer := range layerSlice {
		layer.ZIndex = i * 10 // Use increments of 10 for easy insertion
	}
}

// analyzeDependencies determines which layers depend on others
func (le *LayerExtractor) analyzeDependencies(layers map[string]*InkscapeLayer) error {
	// For now, implement basic dependency analysis
	// In a real implementation, this would analyze layer relationships
	
	for _, layer := range layers {
		// Determine object/boundary type based on layer name
		layer.ObjectType = le.determineObjectType(layer)
		layer.BoundaryType = le.determineBoundaryType(layer)
	}
	
	return nil
}

// determineObjectType determines if a layer represents an object
func (le *LayerExtractor) determineObjectType(layer *InkscapeLayer) ObjectType {
	// Check against object naming patterns
	for namePattern, objectType := range InkscapeObjectMapping {
		if strings.Contains(layer.ID, namePattern) || strings.Contains(layer.Label, namePattern) {
			return objectType
		}
	}
	
	// No object type detected
	return ""
}

// determineBoundaryType determines if a layer represents a boundary
func (le *LayerExtractor) determineBoundaryType(layer *InkscapeLayer) BoundaryType {
	// Check against boundary naming patterns
	for namePattern, boundaryType := range InkscapeBoundaryMapping {
		if strings.Contains(layer.ID, namePattern) || strings.Contains(layer.Label, namePattern) {
			return boundaryType
		}
	}
	
	// No boundary type detected
	return ""
}

// ExtractViewBox extracts the viewBox from the SVG document
func (le *LayerExtractor) ExtractViewBox(svgDoc *SVGDocument) (Rectangle, error) {
	if svgDoc.ViewBox == "" {
		// Default to standard card size if no viewBox
		return Rectangle{X: 0, Y: 0, Width: 1500, Height: 2100}, nil
	}
	
	// Parse viewBox string "minX minY width height"
	parts := strings.Fields(svgDoc.ViewBox)
	if len(parts) != 4 {
		return Rectangle{}, fmt.Errorf("invalid viewBox format: %s", svgDoc.ViewBox)
	}
	
	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("invalid viewBox X: %s", parts[0])
	}
	
	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("invalid viewBox Y: %s", parts[1])
	}
	
	width, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("invalid viewBox width: %s", parts[2])
	}
	
	height, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Rectangle{}, fmt.Errorf("invalid viewBox height: %s", parts[3])
	}
	
	return Rectangle{X: x, Y: y, Width: width, Height: height}, nil
}

// ValidateLayerStructure validates that extracted layers form a valid template
func (le *LayerExtractor) ValidateLayerStructure(layers map[string]*InkscapeLayer) error {
	if len(layers) == 0 {
		return fmt.Errorf("no layers found in SVG document")
	}
	
	// Check for required layer types
	hasObjects := false
	hasBoundaries := false
	
	for _, layer := range layers {
		if layer.ObjectType != "" {
			hasObjects = true
		}
		if layer.BoundaryType != "" {
			hasBoundaries = true
		}
	}
	
	if !hasObjects {
		return fmt.Errorf("no object layers found - template must have visual elements")
	}
	
	if !hasBoundaries {
		return fmt.Errorf("no boundary layers found - template must have text areas")
	}
	
	return nil
}

// CreateLayerSummary creates a summary of extracted layers for debugging
func (le *LayerExtractor) CreateLayerSummary(layers map[string]*InkscapeLayer) map[string]interface{} {
	summary := map[string]interface{}{
		"total_layers": len(layers),
		"objects":      make([]string, 0),
		"boundaries":   make([]string, 0),
		"unclassified": make([]string, 0),
	}
	
	for _, layer := range layers {
		if layer.ObjectType != "" {
			summary["objects"] = append(summary["objects"].([]string), string(layer.ObjectType))
		} else if layer.BoundaryType != "" {
			summary["boundaries"] = append(summary["boundaries"].([]string), string(layer.BoundaryType))
		} else {
			summary["unclassified"] = append(summary["unclassified"].([]string), layer.ID)
		}
	}
	
	return summary
} 