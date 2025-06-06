// object_detector.go detects visual elements by Inkscape naming convention
package ingestion

import (
	"fmt"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// ObjectDetector detects visual elements by Inkscape naming convention
type ObjectDetector struct {
	namingRules map[string]ObjectType
}

// NewObjectDetector creates a new object detector with naming rules
func NewObjectDetector() *ObjectDetector {
	return &ObjectDetector{
		namingRules: InkscapeObjectMapping,
	}
}

// InkscapeObjectMapping maps Inkscape layer names to object types
// Based on the naming convention outlined in SVG_TEMPLATE_SYSTEM_PLAN.md
var InkscapeObjectMapping = map[string]ObjectType{
	// Frame objects (card structure)
	"frame-base":           ObjectFrameBase,
	"frame-border":         ObjectFrameBorder,
	"frame-creature":       ObjectFrameCreature,
	"frame-anthem":         ObjectFrameAnthem,
	"frame-artifact":       ObjectFrameArtifact,
	"frame-spell":          ObjectFrameSpell,
	
	// Text styling objects
	"text-style-name":      ObjectNameTitle,
	"text-style-effect":    ObjectEffectBody,
	"text-style-stats":     ObjectStatsText,
	
	// Special objects
	"art-frame":           ObjectArtFrame,
	"anthem-glow":         ObjectAnthemGlow,
	"set-icon":            ObjectSetIcon,
}

// Card type-specific object selection as outlined in the plan
var CardTypeObjects = map[card.CardType][]ObjectType{
	card.TypeCreature: {
		ObjectFrameBase,
		ObjectFrameCreature,  // Standard creature frame
		ObjectNameTitle,
		ObjectEffectBody,
		ObjectArtFrame,
		ObjectStatsText,
	},
	card.TypeAnthem: {
		ObjectFrameBase,
		ObjectFrameAnthem,    // Red anthem frame variant
		ObjectNameTitle,
		ObjectEffectBody,
		ObjectArtFrame,
		ObjectAnthemGlow,     // Special anthem-only effects
	},
	card.TypeArtifact: {
		ObjectFrameBase,
		ObjectFrameArtifact,  // Artifact-specific frame
		ObjectNameTitle,
		ObjectEffectBody,
		ObjectArtFrame,
	},
	card.TypeSpell: {
		ObjectFrameBase,
		ObjectFrameSpell,
		ObjectNameTitle,
		ObjectEffectBody,
		ObjectArtFrame,
	},
}

// Style variations per card type as outlined in the plan
var CardTypeStyles = map[card.CardType]map[ObjectType]*ObjectStyle{
	card.TypeAnthem: {
		ObjectFrameBase: &ObjectStyle{
			Fill:        "#8B0000",  // Dark red for anthem
			Stroke:      "#FF4500",  // Orange-red border
			StrokeWidth: 3.0,
			Opacity:     1.0,
		},
		ObjectAnthemGlow: &ObjectStyle{
			Fill:        "#FF6347",  // Tomato red glow
			Opacity:     0.7,
		},
	},
	card.TypeCreature: {
		ObjectFrameBase: &ObjectStyle{
			Fill:        "#228B22",  // Forest green for creatures
			Stroke:      "#006400",  // Dark green border
			StrokeWidth: 2.0,
			Opacity:     1.0,
		},
	},
	card.TypeArtifact: {
		ObjectFrameBase: &ObjectStyle{
			Fill:        "#C0C0C0",  // Silver for artifacts
			Stroke:      "#808080",  // Gray border
			StrokeWidth: 2.0,
			Opacity:     1.0,
		},
	},
	card.TypeSpell: {
		ObjectFrameBase: &ObjectStyle{
			Fill:        "#4169E1",  // Royal blue for spells
			Stroke:      "#191970",  // Midnight blue border
			StrokeWidth: 2.0,
			Opacity:     1.0,
		},
	},
}

// DetectObjects analyzes layers and identifies visual objects by naming convention
func (od *ObjectDetector) DetectObjects(layers map[string]*InkscapeLayer) (map[ObjectType]*CardObject, error) {
	objects := make(map[ObjectType]*CardObject)
	
	for layerID, layer := range layers {
		// Try to match layer name/label to object type
		objectType := od.detectObjectTypeFromLayer(layer)
		if objectType == "" {
			// This layer is not a recognized object, skip it
			continue
		}
		
		// Create card object from layer
		cardObject := &CardObject{
			Type:            objectType,
			InkscapeID:      layerID,
			SVGContent:      layer.Content,
			Style:           od.getDefaultStyleForObject(objectType),
			Dependencies:    od.getDependenciesForObject(objectType),
			ZIndex:          layer.ZIndex,
			BlendMode:       layer.BlendMode,
			TransparencyMap: od.createOpacityMapFromLayer(layer),
		}
		
		objects[objectType] = cardObject
	}
	
	return objects, nil
}

// detectObjectTypeFromLayer determines object type from layer naming
func (od *ObjectDetector) detectObjectTypeFromLayer(layer *InkscapeLayer) ObjectType {
	// First, try exact match on layer ID
	if objectType, exists := od.namingRules[layer.ID]; exists {
		return objectType
	}
	
	// Try exact match on layer label
	if objectType, exists := od.namingRules[layer.Label]; exists {
		return objectType
	}
	
	// Try substring matching (for layers like "frame-creature-main" matching "frame-creature")
	for namePattern, objectType := range od.namingRules {
		if strings.Contains(layer.ID, namePattern) || strings.Contains(layer.Label, namePattern) {
			return objectType
		}
	}
	
	// No match found
	return ""
}

// getDefaultStyleForObject returns default styling for an object type
func (od *ObjectDetector) getDefaultStyleForObject(objectType ObjectType) *ObjectStyle {
	// Default style that works for most objects
	return &ObjectStyle{
		Fill:        "#000000",  // Black fill
		Stroke:      "#FFFFFF",  // White stroke
		StrokeWidth: 1.0,
		Opacity:     1.0,
		Transform:   "",
		CSSClasses:  []string{},
	}
}

// getDependenciesForObject returns objects that this object depends on
func (od *ObjectDetector) getDependenciesForObject(objectType ObjectType) []ObjectType {
	dependencies := make([]ObjectType, 0)
	
	// Define object dependencies
	switch objectType {
	case ObjectFrameBorder:
		dependencies = append(dependencies, ObjectFrameBase)
	case ObjectFrameCreature, ObjectFrameAnthem, ObjectFrameArtifact, ObjectFrameSpell:
		dependencies = append(dependencies, ObjectFrameBase)
	case ObjectAnthemGlow:
		dependencies = append(dependencies, ObjectFrameAnthem)
	case ObjectArtFrame:
		dependencies = append(dependencies, ObjectFrameBase)
	}
	
	return dependencies
}

// createOpacityMapFromLayer creates an opacity map for positioning
func (od *ObjectDetector) createOpacityMapFromLayer(layer *InkscapeLayer) *OpacityMap {
	return &OpacityMap{
		PrimaryZone:      layer.VisibleArea,
		FadeZones:        layer.FadeZones,
		FullyVisible:     []Rectangle{layer.VisibleArea},
		PartiallyVisible: make([]GradientZone, 0),
		FullyHidden:      make([]Rectangle, 0),
	}
}

// GetObjectsForCardType returns objects that should be included for a card type
func (od *ObjectDetector) GetObjectsForCardType(cardType card.CardType) []ObjectType {
	if objects, exists := CardTypeObjects[cardType]; exists {
		return objects
	}
	
	// Default objects for unknown card types
	return []ObjectType{
		ObjectFrameBase,
		ObjectNameTitle,
		ObjectEffectBody,
		ObjectArtFrame,
	}
}

// GetStyleForCardType returns card type-specific styling
func (od *ObjectDetector) GetStyleForCardType(cardType card.CardType, objectType ObjectType) *ObjectStyle {
	if cardStyles, exists := CardTypeStyles[cardType]; exists {
		if style, exists := cardStyles[objectType]; exists {
			return style
		}
	}
	
	// Return default style if no card-type-specific style exists
	return od.getDefaultStyleForObject(objectType)
}

// ValidateObjectCompatibility checks if objects are compatible with each other
func (od *ObjectDetector) ValidateObjectCompatibility(objects map[ObjectType]*CardObject) error {
	// Check that dependencies are satisfied
	for objectType, cardObject := range objects {
		for _, dependency := range cardObject.Dependencies {
			if _, exists := objects[dependency]; !exists {
				return fmt.Errorf("object %s requires dependency %s which is missing", objectType, dependency)
			}
		}
	}
	
	// Check for conflicting objects (e.g., multiple frame types)
	frameTypes := []ObjectType{ObjectFrameCreature, ObjectFrameAnthem, ObjectFrameArtifact, ObjectFrameSpell}
	frameCount := 0
	for _, frameType := range frameTypes {
		if _, exists := objects[frameType]; exists {
			frameCount++
		}
	}
	
	if frameCount > 1 {
		return fmt.Errorf("multiple frame types detected, only one frame type allowed per template")
	}
	
	return nil
}

// ApplyCardTypeStyle applies card type-specific styling to objects
func (od *ObjectDetector) ApplyCardTypeStyle(objects map[ObjectType]*CardObject, cardType card.CardType) {
	for objectType, cardObject := range objects {
		if style := od.GetStyleForCardType(cardType, objectType); style != nil {
			// Apply card type-specific style
			cardObject.Style = style
		}
	}
} 