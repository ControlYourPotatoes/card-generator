// boundary_finder.go detects text rendering areas by naming convention
package ingestion

import (
	"fmt"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// BoundaryFinder finds text rendering areas by naming convention
type BoundaryFinder struct {
	textAreaRules map[string]BoundaryType
}

// NewBoundaryFinder creates a new boundary finder with text area rules
func NewBoundaryFinder() *BoundaryFinder {
	return &BoundaryFinder{
		textAreaRules: InkscapeBoundaryMapping,
	}
}

// InkscapeBoundaryMapping maps Inkscape layer names to boundary types
// Only areas that contain text/symbols get boundaries (text rendering guardrails)
var InkscapeBoundaryMapping = map[string]BoundaryType{
	"boundary-name-text":       BoundaryNameText,
	"boundary-effect-text":     BoundaryEffectText,
	"boundary-cost-symbols":    BoundaryCostSymbols,
	"boundary-keyword-symbols": BoundaryKeywordSymbols,
	"boundary-stats-text":      BoundaryStatsText,
	"boundary-set-icon":        BoundarySetIcon,
}

// DetectBoundaries analyzes layers and identifies text boundaries by naming convention
func (bf *BoundaryFinder) DetectBoundaries(layers map[string]*InkscapeLayer) (map[BoundaryType]*TextBoundary, error) {
	boundaries := make(map[BoundaryType]*TextBoundary)
	
	for _, layer := range layers {
		// Try to match layer name/label to boundary type
		boundaryType := bf.detectBoundaryTypeFromLayer(layer)
		if boundaryType == "" {
			// This layer is not a recognized boundary, skip it
			continue
		}
		
		// Create text boundary from layer
		textBoundary := &TextBoundary{
			Type:            boundaryType,
			SafeZone:        layer.VisibleArea,
			PreferredZone:   bf.calculatePreferredZone(layer.VisibleArea),
			FontConstraints: bf.getFontConstraintsForBoundary(boundaryType),
			ContentType:     bf.getContentTypeForBoundary(boundaryType),
			MaxCharacters:   bf.getMaxCharactersForBoundary(boundaryType),
			LineHeight:      bf.getLineHeightForBoundary(boundaryType),
			Alignment:       bf.getAlignmentForBoundary(boundaryType),
		}
		
		boundaries[boundaryType] = textBoundary
	}
	
	return boundaries, nil
}

// detectBoundaryTypeFromLayer determines boundary type from layer naming
func (bf *BoundaryFinder) detectBoundaryTypeFromLayer(layer *InkscapeLayer) BoundaryType {
	// First, try exact match on layer ID
	if boundaryType, exists := bf.textAreaRules[layer.ID]; exists {
		return boundaryType
	}
	
	// Try exact match on layer label
	if boundaryType, exists := bf.textAreaRules[layer.Label]; exists {
		return boundaryType
	}
	
	// Try substring matching (for layers like "boundary-name-text-main" matching "boundary-name-text")
	for namePattern, boundaryType := range bf.textAreaRules {
		if strings.Contains(layer.ID, namePattern) || strings.Contains(layer.Label, namePattern) {
			return boundaryType
		}
	}
	
	// No match found
	return ""
}

// calculatePreferredZone creates a preferred zone that's slightly smaller than safe zone
func (bf *BoundaryFinder) calculatePreferredZone(safeZone Rectangle) Rectangle {
	// Create preferred zone with 10% margin inside safe zone
	marginX := safeZone.Width * 0.05  // 5% margin on each side
	marginY := safeZone.Height * 0.05 // 5% margin on top and bottom
	
	return Rectangle{
		X:      safeZone.X + marginX,
		Y:      safeZone.Y + marginY,
		Width:  safeZone.Width - (marginX * 2),
		Height: safeZone.Height - (marginY * 2),
	}
}

// getFontConstraintsForBoundary returns font constraints based on boundary type
func (bf *BoundaryFinder) getFontConstraintsForBoundary(boundaryType BoundaryType) FontConstraints {
	switch boundaryType {
	case BoundaryNameText:
		return FontConstraints{
			MinSize:       24.0,
			MaxSize:       48.0,
			PreferredSize: 36.0,
			FontFamily:    "serif",
			FontWeight:    "bold",
			AllowBold:     true,
			AllowItalic:   false,
		}
	case BoundaryEffectText:
		return FontConstraints{
			MinSize:       12.0,
			MaxSize:       24.0,
			PreferredSize: 16.0,
			FontFamily:    "sans-serif",
			FontWeight:    "normal",
			AllowBold:     true,
			AllowItalic:   true,
		}
	case BoundaryCostSymbols:
		return FontConstraints{
			MinSize:       20.0,
			MaxSize:       40.0,
			PreferredSize: 32.0,
			FontFamily:    "monospace",
			FontWeight:    "normal",
			AllowBold:     false,
			AllowItalic:   false,
		}
	case BoundaryKeywordSymbols:
		return FontConstraints{
			MinSize:       14.0,
			MaxSize:       24.0,
			PreferredSize: 18.0,
			FontFamily:    "sans-serif",
			FontWeight:    "bold",
			AllowBold:     true,
			AllowItalic:   false,
		}
	case BoundaryStatsText:
		return FontConstraints{
			MinSize:       18.0,
			MaxSize:       36.0,
			PreferredSize: 24.0,
			FontFamily:    "serif",
			FontWeight:    "bold",
			AllowBold:     true,
			AllowItalic:   false,
		}
	case BoundarySetIcon:
		return FontConstraints{
			MinSize:       16.0,
			MaxSize:       32.0,
			PreferredSize: 24.0,
			FontFamily:    "sans-serif",
			FontWeight:    "normal",
			AllowBold:     false,
			AllowItalic:   false,
		}
	default:
		// Default font constraints
		return FontConstraints{
			MinSize:       12.0,
			MaxSize:       24.0,
			PreferredSize: 16.0,
			FontFamily:    "sans-serif",
			FontWeight:    "normal",
			AllowBold:     true,
			AllowItalic:   true,
		}
	}
}

// getContentTypeForBoundary returns the content type for a boundary
func (bf *BoundaryFinder) getContentTypeForBoundary(boundaryType BoundaryType) ContentType {
	switch boundaryType {
	case BoundaryNameText, BoundaryEffectText, BoundaryStatsText:
		return ContentText
	case BoundaryCostSymbols, BoundaryKeywordSymbols, BoundarySetIcon:
		return ContentSymbol
	default:
		return ContentMixed
	}
}

// getMaxCharactersForBoundary returns character limits for boundaries
func (bf *BoundaryFinder) getMaxCharactersForBoundary(boundaryType BoundaryType) int {
	switch boundaryType {
	case BoundaryNameText:
		return 25 // Card names should be concise
	case BoundaryEffectText:
		return 300 // Effect text can be longer
	case BoundaryCostSymbols:
		return 5 // Max 5 cost symbols
	case BoundaryKeywordSymbols:
		return 10 // Max 10 keyword symbols
	case BoundaryStatsText:
		return 10 // "99/99" style stats
	case BoundarySetIcon:
		return 1 // Single set icon
	default:
		return 50 // Default limit
	}
}

// getLineHeightForBoundary returns line height for boundaries
func (bf *BoundaryFinder) getLineHeightForBoundary(boundaryType BoundaryType) float64 {
	switch boundaryType {
	case BoundaryNameText:
		return 1.2 // Tight line height for titles
	case BoundaryEffectText:
		return 1.4 // More readable line height for effect text
	case BoundaryCostSymbols, BoundaryKeywordSymbols, BoundarySetIcon:
		return 1.0 // Symbols don't need line spacing
	case BoundaryStatsText:
		return 1.1 // Tight for stats
	default:
		return 1.3 // Default line height
	}
}

// getAlignmentForBoundary returns text alignment for boundaries
func (bf *BoundaryFinder) getAlignmentForBoundary(boundaryType BoundaryType) TextAlignment {
	switch boundaryType {
	case BoundaryNameText:
		return AlignCenter // Card names are centered
	case BoundaryEffectText:
		return AlignLeft // Effect text is left-aligned for readability
	case BoundaryCostSymbols, BoundarySetIcon:
		return AlignCenter // Symbols are centered
	case BoundaryKeywordSymbols:
		return AlignLeft // Keywords are left-aligned
	case BoundaryStatsText:
		return AlignCenter // Stats are centered
	default:
		return AlignLeft // Default alignment
	}
}

// ValidateBoundaryRequirements checks if boundaries meet card type requirements
func (bf *BoundaryFinder) ValidateBoundaryRequirements(boundaries map[BoundaryType]*TextBoundary, cardType card.CardType) error {
	requiredBoundaries := GetRequiredBoundariesForCardType(cardType)
	
	for _, required := range requiredBoundaries {
		if _, exists := boundaries[required]; !exists {
			return fmt.Errorf("missing required boundary %s for card type %s", required, cardType)
		}
	}
	
	return nil
}

// GetBoundariesForCardType returns boundaries that should be included for a card type
func (bf *BoundaryFinder) GetBoundariesForCardType(cardType card.CardType) []BoundaryType {
	if boundaries, exists := CardTypeBoundaryRequirements[cardType]; exists {
		return boundaries
	}
	
	// Default boundaries for unknown card types
	return []BoundaryType{
		BoundaryNameText,
		BoundaryEffectText,
		BoundaryCostSymbols,
	}
}

// OptimizeBoundaryForContent adjusts boundary constraints based on content requirements
func (bf *BoundaryFinder) OptimizeBoundaryForContent(boundary *TextBoundary, contentLength int, hasSymbols bool) {
	// Adjust font size based on content length
	if contentLength > boundary.MaxCharacters {
		// Content is too long, reduce font size
		currentSize := boundary.FontConstraints.PreferredSize
		reductionFactor := float64(boundary.MaxCharacters) / float64(contentLength)
		newSize := currentSize * reductionFactor
		
		if newSize < boundary.FontConstraints.MinSize {
			newSize = boundary.FontConstraints.MinSize
		}
		
		boundary.FontConstraints.PreferredSize = newSize
	}
	
	// Adjust content type if symbols are present
	if hasSymbols && boundary.ContentType == ContentText {
		boundary.ContentType = ContentMixed
	}
}

// CreateBoundaryExamples creates example boundaries for testing and documentation
func (bf *BoundaryFinder) CreateBoundaryExamples() map[BoundaryType]*TextBoundary {
	examples := make(map[BoundaryType]*TextBoundary)
	
	// Name title area example
	examples[BoundaryNameText] = &TextBoundary{
		Type:      BoundaryNameText,
		SafeZone:  Rectangle{X: 135, Y: 100, Width: 1230, Height: 60},
		PreferredZone: Rectangle{X: 150, Y: 110, Width: 1200, Height: 40},
		FontConstraints: bf.getFontConstraintsForBoundary(BoundaryNameText),
		ContentType:     bf.getContentTypeForBoundary(BoundaryNameText),
		MaxCharacters:   bf.getMaxCharactersForBoundary(BoundaryNameText),
		LineHeight:      bf.getLineHeightForBoundary(BoundaryNameText),
		Alignment:       bf.getAlignmentForBoundary(BoundaryNameText),
	}
	
	// Cost symbol area example
	examples[BoundaryCostSymbols] = &TextBoundary{
		Type:      BoundaryCostSymbols,
		SafeZone:  Rectangle{X: 1320, Y: 90, Width: 100, Height: 100},
		PreferredZone: Rectangle{X: 1330, Y: 100, Width: 80, Height: 80},
		FontConstraints: bf.getFontConstraintsForBoundary(BoundaryCostSymbols),
		ContentType:     bf.getContentTypeForBoundary(BoundaryCostSymbols),
		MaxCharacters:   bf.getMaxCharactersForBoundary(BoundaryCostSymbols),
		LineHeight:      bf.getLineHeightForBoundary(BoundaryCostSymbols),
		Alignment:       bf.getAlignmentForBoundary(BoundaryCostSymbols),
	}
	
	// Effect text area example
	examples[BoundaryEffectText] = &TextBoundary{
		Type:      BoundaryEffectText,
		SafeZone:  Rectangle{X: 135, Y: 1400, Width: 1230, Height: 300},
		PreferredZone: Rectangle{X: 150, Y: 1420, Width: 1200, Height: 260},
		FontConstraints: bf.getFontConstraintsForBoundary(BoundaryEffectText),
		ContentType:     bf.getContentTypeForBoundary(BoundaryEffectText),
		MaxCharacters:   bf.getMaxCharactersForBoundary(BoundaryEffectText),
		LineHeight:      bf.getLineHeightForBoundary(BoundaryEffectText),
		Alignment:       bf.getAlignmentForBoundary(BoundaryEffectText),
	}
	
	// Stats text area example (for creatures)
	examples[BoundaryStatsText] = &TextBoundary{
		Type:      BoundaryStatsText,
		SafeZone:  Rectangle{X: 1200, Y: 1950, Width: 150, Height: 80},
		PreferredZone: Rectangle{X: 1210, Y: 1960, Width: 130, Height: 60},
		FontConstraints: bf.getFontConstraintsForBoundary(BoundaryStatsText),
		ContentType:     bf.getContentTypeForBoundary(BoundaryStatsText),
		MaxCharacters:   bf.getMaxCharactersForBoundary(BoundaryStatsText),
		LineHeight:      bf.getLineHeightForBoundary(BoundaryStatsText),
		Alignment:       bf.getAlignmentForBoundary(BoundaryStatsText),
	}
	
	return examples
} 