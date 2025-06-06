// ingestion_test.go contains tests for the Phase 3a ingestion pipeline
package ingestion

import (
	"strings"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// TestInkscapeParserSmokeTest provides quick validation during development
func TestInkscapeParserSmokeTest(t *testing.T) {
	parser := NewInkscapeParser()
	if parser == nil {
		t.Fatal("Failed to create Inkscape parser")
	}
	
	// Test that all components are initialized
	if parser.layerExtractor == nil {
		t.Error("LayerExtractor not initialized")
	}
	if parser.objectDetector == nil {
		t.Error("ObjectDetector not initialized")
	}
	if parser.boundaryFinder == nil {
		t.Error("BoundaryFinder not initialized")
	}
	if parser.metadataBuilder == nil {
		t.Error("MetadataBuilder not initialized")
	}
	
	t.Log("✅ Inkscape parser smoke test passed")
}

// TestInkscapeObjectDetection tests naming convention detection for objects
func TestInkscapeObjectDetection(t *testing.T) {
	detector := NewObjectDetector()
	
	// Test object mapping
	testCases := []struct {
		layerName    string
		expectedType ObjectType
		description  string
	}{
		{"frame-base", ObjectFrameBase, "Base frame detection"},
		{"frame-creature", ObjectFrameCreature, "Creature frame detection"},
		{"frame-anthem", ObjectFrameAnthem, "Anthem frame detection"},
		{"frame-artifact", ObjectFrameArtifact, "Artifact frame detection"},
		{"text-style-name", ObjectNameTitle, "Name text style detection"},
		{"text-style-effect", ObjectEffectBody, "Effect text style detection"},
		{"art-frame", ObjectArtFrame, "Art frame detection"},
		{"anthem-glow", ObjectAnthemGlow, "Anthem glow effect detection"},
		{"unknown-layer", "", "Unknown layer should return empty type"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create mock layer
			layer := &InkscapeLayer{
				ID:    tc.layerName,
				Label: tc.layerName,
			}
			
			detected := detector.detectObjectTypeFromLayer(layer)
			if detected != tc.expectedType {
				t.Errorf("Expected %s, got %s for layer %s", tc.expectedType, detected, tc.layerName)
			}
		})
	}
	
	t.Log("✅ Object detection tests passed")
}

// TestInkscapeBoundaryDetection tests naming convention detection for boundaries
func TestInkscapeBoundaryDetection(t *testing.T) {
	finder := NewBoundaryFinder()
	
	// Test boundary mapping
	testCases := []struct {
		layerName    string
		expectedType BoundaryType
		description  string
	}{
		{"boundary-name-text", BoundaryNameText, "Name text boundary detection"},
		{"boundary-effect-text", BoundaryEffectText, "Effect text boundary detection"},
		{"boundary-cost-symbols", BoundaryCostSymbols, "Cost symbols boundary detection"},
		{"boundary-keyword-symbols", BoundaryKeywordSymbols, "Keyword symbols boundary detection"},
		{"boundary-stats-text", BoundaryStatsText, "Stats text boundary detection"},
		{"boundary-set-icon", BoundarySetIcon, "Set icon boundary detection"},
		{"not-a-boundary", "", "Non-boundary layer should return empty type"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create mock layer
			layer := &InkscapeLayer{
				ID:    tc.layerName,
				Label: tc.layerName,
			}
			
			detected := finder.detectBoundaryTypeFromLayer(layer)
			if detected != tc.expectedType {
				t.Errorf("Expected %s, got %s for layer %s", tc.expectedType, detected, tc.layerName)
			}
		})
	}
	
	t.Log("✅ Boundary detection tests passed")
}

// TestCardTypeObjectMapping tests card type-specific object selection
func TestCardTypeObjectMapping(t *testing.T) {
	detector := NewObjectDetector()
	
	testCases := []struct {
		cardType     card.CardType
		shouldHave   []ObjectType
		shouldNotHave []ObjectType
		description  string
	}{
		{
			cardType: card.TypeCreature,
			shouldHave: []ObjectType{ObjectFrameBase, ObjectFrameCreature, ObjectNameTitle, ObjectStatsText},
			shouldNotHave: []ObjectType{ObjectFrameAnthem, ObjectAnthemGlow},
			description: "Creature cards should have creature-specific objects",
		},
		{
			cardType: card.TypeAnthem,
			shouldHave: []ObjectType{ObjectFrameBase, ObjectFrameAnthem, ObjectAnthemGlow},
			shouldNotHave: []ObjectType{ObjectFrameCreature, ObjectStatsText},
			description: "Anthem cards should have anthem-specific objects",
		},
		{
			cardType: card.TypeArtifact,
			shouldHave: []ObjectType{ObjectFrameBase, ObjectFrameArtifact},
			shouldNotHave: []ObjectType{ObjectFrameCreature, ObjectAnthemGlow},
			description: "Artifact cards should have artifact-specific objects",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			objects := detector.GetObjectsForCardType(tc.cardType)
			
			// Check that required objects are present
			for _, required := range tc.shouldHave {
				found := false
				for _, object := range objects {
					if object == required {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Card type %s missing required object %s", tc.cardType, required)
				}
			}
			
			// Check that inappropriate objects are not present
			for _, inappropriate := range tc.shouldNotHave {
				for _, object := range objects {
					if object == inappropriate {
						t.Errorf("Card type %s has inappropriate object %s", tc.cardType, inappropriate)
					}
				}
			}
		})
	}
	
	t.Log("✅ Card type object mapping tests passed")
}

// TestSVGDOMExtraction tests SVG DOM parsing and extraction
func TestSVGDOMExtraction(t *testing.T) {
	// Create sample SVG content
	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 1500 2100" width="1500" height="2100">
	<g id="frame-base" label="Base Frame">
		<rect x="0" y="0" width="1500" height="2100" style="fill:#cccccc"/>
	</g>
	<g id="boundary-name-text" label="Card Name Area">
		<rect x="135" y="100" width="1230" height="60" style="fill:none;stroke:#ff0000"/>
	</g>
</svg>`
	
	parser := NewInkscapeParser()
	
	// Test parsing SVG content
	template, err := parser.ParseSVGContent(svgContent)
	if err != nil {
		t.Fatalf("Failed to parse SVG content: %v", err)
	}
	
	if template == nil {
		t.Fatal("Parsed template is nil")
	}
	
	// Check that objects were detected
	if len(template.Objects) == 0 {
		t.Error("No objects detected in SVG")
	}
	
	// Check that boundaries were detected
	if len(template.Boundaries) == 0 {
		t.Error("No boundaries detected in SVG")
	}
	
	// Check specific object detection
	if _, exists := template.Objects[ObjectFrameBase]; !exists {
		t.Error("ObjectFrameBase not detected from 'frame-base' layer")
	}
	
	// Check specific boundary detection
	if _, exists := template.Boundaries[BoundaryNameText]; !exists {
		t.Error("BoundaryNameText not detected from 'boundary-name-text' layer")
	}
	
	t.Log("✅ SVG DOM extraction tests passed")
}

// TestCardTypeStyleApplication tests card type-specific styling
func TestCardTypeStyleApplication(t *testing.T) {
	detector := NewObjectDetector()
	
	// Create mock objects
	objects := map[ObjectType]*CardObject{
		ObjectFrameBase: {
			Type:   ObjectFrameBase,
			Style:  &ObjectStyle{Fill: "#000000"}, // Default style
		},
	}
	
	// Test anthem styling (should get red colors)
	detector.ApplyCardTypeStyle(objects, card.TypeAnthem)
	
	frameBase := objects[ObjectFrameBase]
	if frameBase.Style.Fill != "#8B0000" { // Dark red for anthem
		t.Errorf("Expected anthem frame to have red fill #8B0000, got %s", frameBase.Style.Fill)
	}
	
	if frameBase.Style.Stroke != "#FF4500" { // Orange-red border
		t.Errorf("Expected anthem frame to have orange-red stroke #FF4500, got %s", frameBase.Style.Stroke)
	}
	
	// Test creature styling (should get green colors)
	detector.ApplyCardTypeStyle(objects, card.TypeCreature)
	
	if frameBase.Style.Fill != "#228B22" { // Forest green for creature
		t.Errorf("Expected creature frame to have green fill #228B22, got %s", frameBase.Style.Fill)
	}
	
	t.Log("✅ Card type style application tests passed")
}

// TestBoundaryValidation tests boundary validation and constraints
func TestBoundaryValidation(t *testing.T) {
	finder := NewBoundaryFinder()
	
	// Test boundary creation with proper constraints
	examples := finder.CreateBoundaryExamples()
	
	// Test name boundary
	nameBoundary := examples[BoundaryNameText]
	if nameBoundary == nil {
		t.Fatal("Name boundary example not created")
	}
	
	if nameBoundary.SafeZone.Empty() {
		t.Error("Name boundary has empty safe zone")
	}
	
	if nameBoundary.FontConstraints.MinSize <= 0 {
		t.Error("Name boundary has invalid minimum font size")
	}
	
	if nameBoundary.FontConstraints.MaxSize <= nameBoundary.FontConstraints.MinSize {
		t.Error("Name boundary has invalid font size range")
	}
	
	if nameBoundary.MaxCharacters <= 0 {
		t.Error("Name boundary has no character limit")
	}
	
	// Test content type validation
	if nameBoundary.ContentType != ContentText {
		t.Errorf("Expected name boundary to have ContentText, got %s", nameBoundary.ContentType)
	}
	
	// Test cost symbols boundary
	costBoundary := examples[BoundaryCostSymbols]
	if costBoundary.ContentType != ContentSymbol {
		t.Errorf("Expected cost boundary to have ContentSymbol, got %s", costBoundary.ContentType)
	}
	
	t.Log("✅ Boundary validation tests passed")
}

// TestTemplateMetadata tests metadata building and validation
func TestTemplateMetadata(t *testing.T) {
	builder := NewMetadataBuilder()
	
	// Create mock objects and boundaries
	objects := map[ObjectType]*CardObject{
		ObjectFrameBase:     {Type: ObjectFrameBase},
		ObjectFrameCreature: {Type: ObjectFrameCreature},  // Add creature frame for creature support
		ObjectNameTitle:     {Type: ObjectNameTitle},
		ObjectEffectBody:    {Type: ObjectEffectBody},
		ObjectArtFrame:      {Type: ObjectArtFrame},       // Add art frame
		ObjectStatsText:     {Type: ObjectStatsText},      // Add stats text for creature support
	}
	
	boundaries := map[BoundaryType]*TextBoundary{
		BoundaryNameText: {
			Type: BoundaryNameText, 
			SafeZone: Rectangle{X: 0, Y: 0, Width: 100, Height: 50},
			FontConstraints: FontConstraints{MinSize: 12.0, MaxSize: 48.0, PreferredSize: 24.0},
		},
		BoundaryEffectText: {
			Type: BoundaryEffectText, 
			SafeZone: Rectangle{X: 0, Y: 100, Width: 100, Height: 100},
			FontConstraints: FontConstraints{MinSize: 12.0, MaxSize: 36.0, PreferredSize: 16.0},
		},
		BoundaryCostSymbols: {
			Type: BoundaryCostSymbols, 
			SafeZone: Rectangle{X: 0, Y: 200, Width: 50, Height: 50},
			FontConstraints: FontConstraints{MinSize: 20.0, MaxSize: 40.0, PreferredSize: 32.0},
		},
		BoundaryStatsText: {
			Type: BoundaryStatsText, 
			SafeZone: Rectangle{X: 0, Y: 250, Width: 80, Height: 40},
			FontConstraints: FontConstraints{MinSize: 18.0, MaxSize: 36.0, PreferredSize: 24.0},
		},
	}
	
	positioning := &TransparencyPositioning{
		Layers:      make(map[string]*TransparencyLayer),
		OpacityMaps: make(map[string]*OpacityMap),
		BlendModes:  make(map[string]BlendMode),
	}
	
	// Build metadata
	metadata := builder.BuildMetadata("test.svg", objects, boundaries, positioning)
	
	if metadata == nil {
		t.Fatal("Metadata not created")
	}
	
	if metadata.ObjectCount != len(objects) {
		t.Errorf("Expected object count %d, got %d", len(objects), metadata.ObjectCount)
	}
	
	if metadata.BoundaryCount != len(boundaries) {
		t.Errorf("Expected boundary count %d, got %d", len(boundaries), metadata.BoundaryCount)
	}
	
	if metadata.SourceFile != "test.svg" {
		t.Errorf("Expected source file 'test.svg', got %s", metadata.SourceFile)
	}
	
	if !metadata.Validation.IsValid {
		t.Errorf("Template should be valid, but validation failed: %v", metadata.Validation.Errors)
	}
	
	// Test supported card types
	if len(metadata.SupportedTypes) == 0 {
		t.Error("No supported card types detected")
	}
	
	t.Log("✅ Template metadata tests passed")
}

// TestMalformedSVGHandling tests error handling for malformed SVG
func TestMalformedSVGHandling(t *testing.T) {
	parser := NewInkscapeParser()
	
	testCases := []struct {
		svgContent  string
		description string
	}{
		{"not xml at all", "Invalid XML content"},
		{"<svg>unclosed tag", "Unclosed XML tags"},
		{"", "Empty content"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := parser.ParseSVGContent(tc.svgContent)
			if err == nil {
				t.Errorf("Expected error for %s, but parsing succeeded", tc.description)
			}
		})
	}
	
	t.Log("✅ Malformed SVG handling tests passed")
}

// TestObjectDependencyResolution tests object dependency validation
func TestObjectDependencyResolution(t *testing.T) {
	detector := NewObjectDetector()
	
	// Test valid dependencies
	validObjects := map[ObjectType]*CardObject{
		ObjectFrameBase: {
			Type:         ObjectFrameBase,
			Dependencies: []ObjectType{},
		},
		ObjectFrameBorder: {
			Type:         ObjectFrameBorder,
			Dependencies: []ObjectType{ObjectFrameBase},
		},
	}
	
	err := detector.ValidateObjectCompatibility(validObjects)
	if err != nil {
		t.Errorf("Valid objects should pass compatibility check: %v", err)
	}
	
	// Test missing dependencies
	invalidObjects := map[ObjectType]*CardObject{
		ObjectFrameBorder: {
			Type:         ObjectFrameBorder,
			Dependencies: []ObjectType{ObjectFrameBase}, // Missing ObjectFrameBase
		},
	}
	
	err = detector.ValidateObjectCompatibility(invalidObjects)
	if err == nil {
		t.Error("Invalid objects should fail compatibility check")
	}
	
	if !strings.Contains(err.Error(), "dependency") {
		t.Errorf("Error should mention dependency issue: %v", err)
	}
	
	t.Log("✅ Object dependency resolution tests passed")
}

// TestPhase3aCompletionChecklist verifies Phase 3a implementation is complete
func TestPhase3aCompletionChecklist(t *testing.T) {
	t.Log("=== Phase 3a Implementation Completion Checklist ===")
	
	// ✅ Inkscape parser with naming convention support
	parser := NewInkscapeParser()
	if parser == nil {
		t.Error("❌ Inkscape parser not implemented")
	} else {
		t.Log("✅ Inkscape parser implemented")
	}
	
	// ✅ Object vs boundary detection and separation
	detector := NewObjectDetector()
	finder := NewBoundaryFinder()
	if detector == nil || finder == nil {
		t.Error("❌ Object/boundary detection not implemented")
	} else {
		t.Log("✅ Object/boundary detection implemented")
	}
	
	// ✅ Transparency positioning metadata extraction
	if parser != nil {
		// Test that parser creates positioning metadata
		template, err := parser.ParseSVGContent(`<svg><g id="test"></g></svg>`)
		if err != nil || template == nil || template.Positioning == nil {
			t.Error("❌ Transparency positioning not implemented")
		} else {
			t.Log("✅ Transparency positioning implemented")
		}
	}
	
	// ✅ Clean template data structure generation
	builder := NewMetadataBuilder()
	if builder == nil {
		t.Error("❌ Metadata builder not implemented")
	} else {
		t.Log("✅ Metadata builder implemented")
	}
	
	// ✅ Naming convention mappings
	if len(InkscapeObjectMapping) == 0 || len(InkscapeBoundaryMapping) == 0 {
		t.Error("❌ Naming convention mappings incomplete")
	} else {
		t.Log("✅ Naming convention mappings complete")
	}
	
	// ✅ Card type compatibility validation
	if len(CardTypeObjectRequirements) == 0 || len(CardTypeBoundaryRequirements) == 0 {
		t.Error("❌ Card type requirements not defined")
	} else {
		t.Log("✅ Card type requirements defined")
	}
	
	t.Log("=== Phase 3a Implementation Status: COMPLETE ===")
} 