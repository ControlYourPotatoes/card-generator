// types.go defines core types for the SVG ingestion pipeline
package ingestion

import (
	"encoding/xml"
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// ObjectType represents different visual elements that can be placed on cards
type ObjectType string

const (
	// Frame objects (card structure)
	ObjectFrameBase      ObjectType = "frame-base"      // Base card frame
	ObjectFrameBorder    ObjectType = "frame-border"    // Frame border decoration
	ObjectFrameCreature  ObjectType = "frame-creature"  // Standard creature frame
	ObjectFrameAnthem    ObjectType = "frame-anthem"    // Red anthem frame
	ObjectFrameArtifact  ObjectType = "frame-artifact"  // Metallic artifact frame
	ObjectFrameSpell     ObjectType = "frame-spell"     // Spell frame

	// Text styling objects
	ObjectNameTitle   ObjectType = "text-style-name"   // Name text styling
	ObjectEffectBody  ObjectType = "text-style-effect" // Effect text styling
	ObjectStatsText   ObjectType = "text-style-stats"  // Stats text styling

	// Special objects
	ObjectArtFrame   ObjectType = "art-frame"     // Art placement frame
	ObjectAnthemGlow ObjectType = "anthem-glow"   // Special anthem glow effect
	ObjectSetIcon    ObjectType = "set-icon"      // Set symbol placement
)

// BoundaryType represents text rendering areas (guardrails for text placement)
type BoundaryType string

const (
	// Text boundaries (areas where text can be placed)
	BoundaryNameText      BoundaryType = "boundary-name-text"      // Card name text area
	BoundaryEffectText    BoundaryType = "boundary-effect-text"    // Effect text area
	BoundaryCostSymbols   BoundaryType = "boundary-cost-symbols"   // Mana cost symbols area
	BoundaryKeywordSymbols BoundaryType = "boundary-keyword-symbols" // Keyword symbols area
	BoundaryStatsText     BoundaryType = "boundary-stats-text"     // Power/toughness text area
	BoundarySetIcon       BoundaryType = "boundary-set-icon"       // Set icon placement area
)

// SymbolType represents special symbols that can be placed in text
type SymbolType string

const (
	SymbolManaCost      SymbolType = "mana-cost"      // Mana cost symbols
	SymbolKeywordHaste  SymbolType = "keyword-haste"  // Haste keyword symbol
	SymbolKeywordFlying SymbolType = "keyword-flying" // Flying keyword symbol
	SymbolSetIcon       SymbolType = "set-icon"       // Set icon symbol
)

// CardObject represents a visual element that can be placed on a card
type CardObject struct {
	Type            ObjectType                 // Type of object
	InkscapeID      string                     // Original Inkscape layer ID
	SVGContent      string                     // Clean SVG markup for this object
	Style           *ObjectStyle               // Visual styling properties
	Dependencies    []ObjectType               // Required other objects
	ZIndex          int                        // Layer ordering (higher = on top)
	BlendMode       BlendMode                  // How it combines with other layers
	TransparencyMap *OpacityMap                // Position via transparency
}

// ObjectStyle defines visual styling for card objects
type ObjectStyle struct {
	Fill        string    // Fill color
	Stroke      string    // Stroke color
	StrokeWidth float64   // Stroke width
	Opacity     float64   // Opacity (0.0 to 1.0)
	Transform   string    // SVG transform attribute
	CSSClasses  []string  // CSS classes to apply
}

// TextBoundary represents a text rendering area with constraints
type TextBoundary struct {
	Type            BoundaryType      // Type of boundary
	SafeZone        Rectangle         // Text must stay within this area
	PreferredZone   Rectangle         // Optimal text placement area
	FontConstraints FontConstraints   // Font size and style limits
	ContentType     ContentType       // Text, symbol, or mixed content
	MaxCharacters   int               // Maximum character limit
	LineHeight      float64           // Line spacing
	Alignment       TextAlignment     // Text alignment
}

// FontConstraints defines font limitations for text boundaries
type FontConstraints struct {
	MinSize       float64 // Minimum font size
	MaxSize       float64 // Maximum font size
	PreferredSize float64 // Preferred font size
	FontFamily    string  // Font family name
	FontWeight    string  // Font weight (normal, bold, etc.)
	AllowBold     bool    // Whether bold text is allowed
	AllowItalic   bool    // Whether italic text is allowed
}

// ContentType defines what kind of content can be placed in a boundary
type ContentType string

const (
	ContentText   ContentType = "text"   // Regular text content
	ContentSymbol ContentType = "symbol" // Symbol content (icons, costs)
	ContentMixed  ContentType = "mixed"  // Text + symbols combined
)

// TextAlignment defines text alignment within boundaries
type TextAlignment string

const (
	AlignLeft   TextAlignment = "left"   // Left-aligned text
	AlignCenter TextAlignment = "center" // Center-aligned text
	AlignRight  TextAlignment = "right"  // Right-aligned text
)

// Rectangle represents a rectangular area (coordinates in SVG space)
type Rectangle struct {
	X      float64 // X coordinate (left)
	Y      float64 // Y coordinate (top)
	Width  float64 // Width
	Height float64 // Height
}

// Empty returns true if the rectangle has zero area
func (r Rectangle) Empty() bool {
	return r.Width <= 0 || r.Height <= 0
}

// ToImageRect converts to Go's image.Rectangle for compatibility
func (r Rectangle) ToImageRect() image.Rectangle {
	return image.Rect(int(r.X), int(r.Y), int(r.X+r.Width), int(r.Y+r.Height))
}

// TransparencyPositioning represents opacity-based positioning system
type TransparencyPositioning struct {
	Layers      map[string]*TransparencyLayer // Positioned layers
	OpacityMaps map[string]*OpacityMap        // Opacity maps for positioning
	BlendModes  map[string]BlendMode          // Blend modes for layers
}

// TransparencyLayer represents a layer positioned via transparency
type TransparencyLayer struct {
	ID              string       // Layer identifier
	ObjectType      ObjectType   // Type of object in this layer
	FullSizeImage   []byte       // Full 1500x2100 layer data (future: actual image data)
	OpacityMask     []byte       // Alpha mask for positioning (future: actual mask data)
	BlendMode       BlendMode    // How it combines with other layers
	ZIndex          int          // Layer ordering
}

// OpacityMap defines areas of different opacity for positioning
type OpacityMap struct {
	PrimaryZone      Rectangle     // Main visible area
	FadeZones        []FadeZone    // Gradient transition areas
	FullyVisible     []Rectangle   // 100% opacity areas
	PartiallyVisible []GradientZone // Variable opacity areas
	FullyHidden      []Rectangle   // 0% opacity areas
}

// FadeZone represents a gradient transition area
type FadeZone struct {
	Area         Rectangle     // Area of the fade
	StartOpacity float64       // Starting opacity (0.0 to 1.0)
	EndOpacity   float64       // Ending opacity (0.0 to 1.0)
	Direction    FadeDirection // Direction of the fade
}

// FadeDirection defines fade gradient direction
type FadeDirection string

const (
	FadeLeftToRight FadeDirection = "left-to-right"
	FadeTopToBottom FadeDirection = "top-to-bottom"
	FadeRadial      FadeDirection = "radial"
)

// GradientZone represents an area with variable opacity
type GradientZone struct {
	Area    Rectangle // Area of the gradient
	Opacity float64   // Opacity value (0.0 to 1.0)
}

// BlendMode defines how layers combine
type BlendMode string

const (
	BlendModeNormal   BlendMode = "normal"   // Normal blending
	BlendModeMultiply BlendMode = "multiply" // Multiply blending
	BlendModeScreen   BlendMode = "screen"   // Screen blending
	BlendModeOverlay  BlendMode = "overlay"  // Overlay blending
)

// TemplateMetadata contains metadata about a parsed template
type TemplateMetadata struct {
	SourceFile     string            // Original Inkscape file path
	ParsedAt       string            // When it was parsed (ISO 8601)
	ObjectCount    int               // Number of objects detected
	BoundaryCount  int               // Number of boundaries detected
	SupportedTypes []card.CardType   // Card types this template supports
	Version        string            // Template version
	Validation     ValidationResult  // Validation status
}

// ValidationResult contains validation information
type ValidationResult struct {
	IsValid bool     // Whether the template is valid
	Errors  []string // List of validation errors
	Warnings []string // List of validation warnings
}

// SVGDocument represents the parsed SVG XML structure
type SVGDocument struct {
	XMLName xml.Name `xml:"svg"`
	ViewBox string   `xml:"viewBox,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	Groups  []SVGGroup `xml:"g"`
}

// SVGGroup represents an SVG group element (usually Inkscape layers)
type SVGGroup struct {
	XMLName xml.Name `xml:"g"`
	ID      string   `xml:"id,attr"`
	Label   string   `xml:"http://www.inkscape.org/namespaces/inkscape label,attr"`
	Style   string   `xml:"style,attr"`
	Elements []interface{} `xml:",any"` // Mixed content (paths, rects, etc.)
}

// InkscapeLayer represents a layer extracted from Inkscape
type InkscapeLayer struct {
	ID          string        // Layer ID from Inkscape
	Label       string        // Human-readable label
	ObjectType  ObjectType    // Detected object type (if any)
	BoundaryType BoundaryType // Detected boundary type (if any)
	VisibleArea Rectangle     // Visible area of this layer
	FadeZones   []FadeZone    // Fade transitions in this layer
	ZIndex      int           // Layer stacking order
	BlendMode   BlendMode     // Blend mode for this layer
	Content     string        // SVG content of this layer
}

// Card type requirements mapping
var CardTypeObjectRequirements = map[card.CardType][]ObjectType{
	card.TypeCreature: {
		ObjectFrameBase,
		ObjectFrameCreature,
		ObjectNameTitle,
		ObjectEffectBody,
		ObjectArtFrame,
		ObjectStatsText,
	},
	card.TypeAnthem: {
		ObjectFrameBase,
		ObjectFrameAnthem,    // Red frame for anthems
		ObjectNameTitle,
		ObjectEffectBody,
		ObjectArtFrame,
		ObjectAnthemGlow,     // Special anthem glow
	},
	card.TypeArtifact: {
		ObjectFrameBase,
		ObjectFrameArtifact,  // Metallic frame for artifacts
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

var CardTypeBoundaryRequirements = map[card.CardType][]BoundaryType{
	card.TypeCreature: {
		BoundaryNameText,
		BoundaryEffectText,
		BoundaryCostSymbols,
		BoundaryStatsText,
	},
	card.TypeAnthem: {
		BoundaryNameText,
		BoundaryEffectText,
		BoundaryCostSymbols,
	},
	card.TypeArtifact: {
		BoundaryNameText,
		BoundaryEffectText,
		BoundaryCostSymbols,
	},
	card.TypeSpell: {
		BoundaryNameText,
		BoundaryEffectText,
		BoundaryCostSymbols,
	},
}

// GetRequiredObjectsForCardType returns required objects for a card type
func GetRequiredObjectsForCardType(cardType card.CardType) []ObjectType {
	if objects, exists := CardTypeObjectRequirements[cardType]; exists {
		return objects
	}
	return []ObjectType{}
}

// GetRequiredBoundariesForCardType returns required boundaries for a card type
func GetRequiredBoundariesForCardType(cardType card.CardType) []BoundaryType {
	if boundaries, exists := CardTypeBoundaryRequirements[cardType]; exists {
		return boundaries
	}
	return []BoundaryType{}
} 