# SVG Template System Implementation Plan

## **Template System Architecture for Card Generator**

### **Purpose**

This document outlines the implementation of the SVG template system with proper separation between Inkscape ingestion, object management, text boundary systems, and transparency-based positioning.

---

## **Architecture Overview**

```
Inkscape Files → Ingestion Pipeline → Template System → SVG Output
     ↓                    ↓                ↓           ↓
Raw .svg files    Clean Objects    Composed Cards   Game-ready SVG
```

### **Separation of Concerns**

1. **Ingestion Pipeline**: Cleans and structures raw Inkscape files
2. **Object Library**: Manages visual elements (frames, backgrounds, decorative)
3. **Boundary System**: Defines text rendering guardrails
4. **Template Composer**: Combines everything into final cards
5. **Positioning Engine**: Uses transparency-based positioning

---

## **Phase 3a: Ingestion Pipeline** ⏱️ **Week 1**

### **Objectives**

- Parse raw Inkscape SVG files into structured components
- Extract objects, boundaries, and positioning metadata
- Create clean, template-ready data structures

### **Implementation**

#### **3a.1 Inkscape Parser**

```go
// Location: backend/internal/generator/svg/ingestion/inkscape_parser.go
type InkscapeParser struct {
    layerExtractor   *LayerExtractor
    objectDetector   *ObjectDetector
    boundaryFinder   *BoundaryFinder
    metadataBuilder  *MetadataBuilder
}

func (p *InkscapeParser) ParseSVGFile(filepath string) (*ParsedTemplate, error) {
    // Parse SVG DOM
    // Extract layers by naming convention
    // Identify objects vs boundaries
    // Build structured metadata
}

type ParsedTemplate struct {
    Objects     map[ObjectType]*CardObject
    Boundaries  map[BoundaryType]*TextBoundary
    Positioning *TransparencyPositioning
    Metadata    *TemplateMetadata
}
```

#### **3a.2 Object Detection**

```go
// Detects visual elements by Inkscape naming convention
type ObjectDetector struct {
    namingRules map[string]ObjectType
}

// Inkscape layer naming → Object types
var InkscapeObjectMapping = map[string]ObjectType{
    "frame-base":           ObjectFrameBase,
    "frame-border":         ObjectFrameBorder,
    "frame-creature":       ObjectFrameCreature,
    "frame-anthem":         ObjectFrameAnthem,
    "frame-artifact":       ObjectFrameArtifact,
    "frame-spell":          ObjectFrameSpell,
    "text-style-name":      ObjectNameTitle,
    "text-style-effect":    ObjectEffectBody,
    "art-frame":           ObjectArtFrame,
    // ... more mappings
}
```

#### **3a.2.1 Card Type-Specific Object Selection**

```go
// Card types determine which objects and styles are used
type CardTypeObjectMap struct {
    objectSelector map[card.CardType][]ObjectType
    styleVariants  map[card.CardType]map[ObjectType]*ObjectStyle
}

// Example: Anthem cards use red frames, creatures use standard frames
var CardTypeObjects = map[card.CardType][]ObjectType{
    card.TypeCreature: {
        ObjectFrameBase,
        ObjectFrameCreature,  // Standard creature frame
        ObjectNameTitle,
        ObjectEffectBody,
        ObjectArtFrame,
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
    // ... more card types
}

// Style variations per card type
var CardTypeStyles = map[card.CardType]map[ObjectType]*ObjectStyle{
    card.TypeAnthem: {
        ObjectFrameBase: &ObjectStyle{
            Fill:        "#8B0000",  // Dark red for anthem
            Stroke:      "#FF4500",  // Orange-red border
            StrokeWidth: 3.0,
            // ... more styling
        },
        ObjectAnthemGlow: &ObjectStyle{
            Fill:        "#FF6347",  // Tomato red glow
            Opacity:     0.7,
            BlendMode:   BlendModeScreen,
        },
    },
    card.TypeCreature: {
        ObjectFrameBase: &ObjectStyle{
            Fill:        "#228B22",  // Forest green for creatures
            Stroke:      "#006400",  // Dark green border
            StrokeWidth: 2.0,
        },
    },
    // ... more type-specific styles
}
```

#### **3a.3 Boundary Detection**

```go
// Finds text rendering areas by naming convention
type BoundaryFinder struct {
    textAreaRules map[string]BoundaryType
}

// Only areas that contain text/symbols get boundaries
var InkscapeBoundaryMapping = map[string]BoundaryType{
    "boundary-name-text":      BoundaryNameText,
    "boundary-effect-text":    BoundaryEffectText,
    "boundary-cost-symbols":   BoundaryCostSymbols,
    "boundary-keyword-symbols": BoundaryKeywordSymbols,
    "boundary-stats-text":     BoundaryStatsText,
    "boundary-set-icon":       BoundarySetIcon,
    // Only text areas get boundaries
}
```

#### **3a.4 Transparency Mapping**

```go
// Creates opacity-based positioning from layers
type TransparencyMapper struct {
    layerComposer *LayerComposer
    opacityBuilder *OpacityBuilder
}

type TransparencyPositioning struct {
    Layers      map[string]*TransparencyLayer
    OpacityMaps map[string]*OpacityMap
    BlendModes  map[string]BlendMode
}
```

---

## **Phase 3b: Template System Core** ⏱️ **Week 2**

### **Objectives**

- Implement structured object and boundary management
- Create template composition system
- Build transparency-based positioning engine

### **Implementation**

#### **3b.1 Object Library**

```go
// Location: backend/internal/generator/svg/templates/object_library.go
type ObjectLibrary struct {
    objects map[ObjectType]*CardObject
    cache   map[string]*CardObject
}

type CardObject struct {
    Type            ObjectType
    InkscapeID      string          // Original layer ID
    SVGContent      string          // Clean SVG markup
    Style           *ObjectStyle    // Visual styling
    Dependencies    []ObjectType    // Required other objects
    ZIndex          int             // Layer ordering
    BlendMode       BlendMode       // How it combines
    TransparencyMap *OpacityMap     // Position via transparency
}

type ObjectStyle struct {
    Fill        string
    Stroke      string
    StrokeWidth float64
    Opacity     float64
    Transform   string
    CSSClasses  []string
}
```

#### **3b.2 Boundary Manager**

```go
// Location: backend/internal/generator/svg/templates/boundary_manager.go
type BoundaryManager struct {
    boundaries map[BoundaryType]*TextBoundary
    validator  *BoundaryValidator
}

type TextBoundary struct {
    Type            BoundaryType
    SafeZone        image.Rectangle  // Text must stay within
    PreferredZone   image.Rectangle  // Optimal text area
    FontConstraints FontConstraints  // Size/style limits
    ContentType     ContentType      // Text vs Symbol
    MaxCharacters   int              // Length limits
    LineHeight      float64          // Spacing rules
    Alignment       TextAlignment    // Left/Center/Right
}

type ContentType string
const (
    ContentText    ContentType = "text"     // Regular text
    ContentSymbol  ContentType = "symbol"   // Cost/keyword symbols
    ContentMixed   ContentType = "mixed"    // Text + symbols
)

type FontConstraints struct {
    MinSize     float64
    MaxSize     float64
    PreferredSize float64
    FontFamily  string
    FontWeight  string
    AllowBold   bool
    AllowItalic bool
}
```

#### **3b.3 Symbol Registry**

```go
// Location: backend/internal/generator/svg/templates/symbol_registry.go
type SymbolRegistry struct {
    symbols map[SymbolType]*Symbol
    cache   map[string]*Symbol
}

type Symbol struct {
    Type        SymbolType
    SVGContent  string          // Symbol markup
    Boundary    BoundaryType    // Where it can be placed
    Size        SymbolSize      // Sizing constraints
    Style       *SymbolStyle    // Visual properties
}

type SymbolSize struct {
    MinWidth   float64
    MaxWidth   float64
    MinHeight  float64
    MaxHeight  float64
    AspectRatio float64  // Maintain proportions
}
```

#### **3b.4 Template Composer**

```go
// Location: backend/internal/generator/svg/templates/template_composer.go
type TemplateComposer struct {
    objectLibrary   *ObjectLibrary
    boundaryManager *BoundaryManager
    symbolRegistry  *SymbolRegistry
    positionEngine  *TransparencyPositionEngine
}

func (tc *TemplateComposer) ComposeCard(
    cardType card.CardType,
    objects []ObjectType,
    boundaries []BoundaryType,
) (*ComposedTemplate, error) {
    // Step 1: Get card type-specific objects (e.g., red frame for anthems)
    typeObjects := CardTypeObjects[cardType]
    typeStyles := CardTypeStyles[cardType]

    // Step 2: Merge requested objects with type requirements
    allObjects := tc.mergeObjectTypes(typeObjects, objects)

    // Step 3: Apply card type-specific styling
    styledObjects := tc.applyTypeStyles(allObjects, typeStyles)

    // Step 4: Validate boundary compatibility
    validBoundaries := tc.validateBoundaries(boundaries, cardType)

    // Step 5: Build transparency positioning
    positioning := tc.buildPositioning(styledObjects, validBoundaries)

    // Step 6: Generate final SVG template
    return tc.generateTemplate(cardType, styledObjects, validBoundaries, positioning)
}

type ComposedTemplate struct {
    SVGTemplate     string                      // Final template
    Objects         map[ObjectType]*CardObject  // Included objects
    Boundaries      map[BoundaryType]*TextBoundary // Text areas
    InteractiveZones []InteractiveZone          // Game features
    AnimationTargets []AnimationTarget          // Effect targets
}
```

---

## **Card Type Configuration Flexibility**

### **Type-Specific Styling Examples**

The template system **maintains full card type configurability** while focusing on objects/boundaries separation:

#### **Anthem Cards: Red Frame Example**

```go
// Anthem cards automatically get red frames
anthemTemplate := templateComposer.ComposeCard(
    card.TypeAnthem,           // Card type drives object selection
    []ObjectType{},            // Additional objects (optional)
    []BoundaryType{           // Text boundaries (same for all types)
        BoundaryNameText,
        BoundaryEffectText,
        BoundaryCostSymbols,
    },
)

// Results in:
// - ObjectFrameAnthem (red frame) instead of ObjectFrameCreature
// - ObjectAnthemGlow (special red glow effect)
// - Red color styling applied automatically
// - Same text boundaries as other cards (consistency)
```

#### **Artifact Cards: Metallic Frame Example**

```go
// Artifact cards get metallic/silver styling
artifactTemplate := templateComposer.ComposeCard(
    card.TypeArtifact,         // Drives silver/metallic frame selection
    []ObjectType{},
    []BoundaryType{
        BoundaryNameText,
        BoundaryEffectText,
        BoundaryCostSymbols,
    },
)

// Results in:
// - ObjectFrameArtifact (metallic frame)
// - Silver/gray color scheme
// - Artifact-specific decorative elements
// - Standard text boundaries (reusable across types)
```

### **Key Benefits of This Approach**

1. **Type Differentiation Preserved**: Each card type gets its unique visual identity
2. **Boundaries Stay Consistent**: Text areas work the same across all card types
3. **Easy Customization**: Change anthem frame color by updating `CardTypeStyles`
4. **Maintainable**: Objects handle visuals, boundaries handle text placement
5. **Extensible**: Add new card types by defining their object/style mappings

---

## **Phase 3c: Transparency Positioning Engine** ⏱️ **Week 2-3**

### **Objectives**

- Implement transparency-based positioning instead of coordinate math
- Create layering system for fixed-size cards (1500x2100)
- Enable easy repositioning via opacity manipulation

### **Implementation**

#### **3c.1 Transparency Position Engine**

```go
// Location: backend/internal/generator/svg/positioning/transparency_engine.go
type TransparencyPositionEngine struct {
    canvasSize   Size             // Always 1500x2100
    layerStack   []*PositionLayer // Ordered layers
    compositor   *LayerCompositor // Combines layers
}

type PositionLayer struct {
    ID              string
    ObjectType      ObjectType
    FullSizeImage   image.Image      // Full 1500x2100 layer
    OpacityMask     image.Image      // Alpha mask for positioning
    BlendMode       BlendMode        // How it combines
    ZIndex          int              // Layer order
}

// Instead of coordinates, position via transparency
func (tpe *TransparencyPositionEngine) PositionObject(
    objectType ObjectType,
    visibleArea image.Rectangle,
) *PositionLayer {
    // Create full-size layer with object
    // Generate opacity mask for visible area
    // Return positioned layer
}
```

#### **3c.2 Opacity Map System**

```go
type OpacityMap struct {
    PrimaryZone    image.Rectangle  // Main visible area
    FadeZones      []FadeZone       // Gradient transitions
    FullyVisible   []image.Rectangle // 100% opacity areas
    PartiallyVisible []GradientZone // Variable opacity
    FullyHidden    []image.Rectangle // 0% opacity areas
}

type FadeZone struct {
    Area        image.Rectangle
    StartOpacity float64
    EndOpacity   float64
    Direction    FadeDirection  // Left-to-right, top-to-bottom, etc.
}

// Benefits: Moving objects = changing opacity maps, not recalculating coordinates
```

---

## **Text Boundary Examples**

### **Name Title Area**

```go
nameBoundary := &TextBoundary{
    Type: BoundaryNameText,
    SafeZone: image.Rect(135, 100, 1365, 160),  // Text must stay here
    PreferredZone: image.Rect(150, 110, 1350, 150), // Ideal placement
    FontConstraints: FontConstraints{
        MinSize: 24.0,
        MaxSize: 48.0,
        PreferredSize: 36.0,
        FontFamily: "serif",
        FontWeight: "bold",
    },
    ContentType: ContentText,
    MaxCharacters: 25,
    Alignment: AlignCenter,
}
```

### **Cost Symbol Area**

```go
costBoundary := &TextBoundary{
    Type: BoundaryCostSymbols,
    SafeZone: image.Rect(1320, 90, 1420, 190), // Top right
    PreferredZone: image.Rect(1330, 100, 1410, 180),
    FontConstraints: FontConstraints{
        MinSize: 20.0,
        MaxSize: 40.0,
        PreferredSize: 32.0,
    },
    ContentType: ContentSymbol,
    MaxCharacters: 5,  // Max symbols
    Alignment: AlignCenter,
}
```

---

## **Integration with Current SVG System**

### **Template Interface Updates**

```go
// Extend existing SVGTemplate interface
type EnhancedSVGTemplate interface {
    SVGTemplate  // Existing interface

    // New methods for enhanced system
    GetObjects() map[ObjectType]*CardObject
    GetBoundaries() map[BoundaryType]*TextBoundary
    GetSymbols() map[SymbolType]*Symbol
    GetTransparencyLayers() map[string]*PositionLayer
}
```

### **Generator Integration**

```go
// Update existing svgGenerator to use new system
func (g *svgGenerator) GenerateSVG(data *card.CardDTO, outputPath string) error {
    // Load enhanced template
    enhancedTemplate := g.templateComposer.ComposeCard(data.Type, ...)

    // Process with transparency positioning
    processedSVG := g.transparencyEngine.RenderCard(enhancedTemplate, data)

    // Write to file
    return os.WriteFile(outputPath, []byte(processedSVG), 0644)
}
```

---

## **Acceptance Criteria**

### **Phase 3a: Ingestion Pipeline**

- [ ] Parse Inkscape SVG files by naming convention
- [ ] Extract objects and boundaries separately
- [ ] Create transparency positioning metadata
- [ ] Generate clean, structured template data

### **Phase 3b: Template System**

- [ ] Object library with visual elements
- [ ] Boundary manager for text areas only
- [ ] Symbol registry for text symbols
- [ ] Template composer combining everything

### **Phase 3c: Positioning Engine**

- [ ] Transparency-based positioning working
- [ ] Fixed 1500x2100 canvas optimization
- [ ] Easy repositioning via opacity changes
- [ ] Integration with existing SVG generator

---

## **Testing Strategy**

### **Unit Tests**

- Inkscape parser with sample files
- Object detection accuracy
- Boundary validation rules
- Transparency positioning math

### **Integration Tests**

- Full pipeline: Inkscape → Template → SVG
- Text boundary guardrails working
- Symbol placement accuracy
- Performance vs coordinate-based system

### **Visual Tests**

- Generated cards match Inkscape designs
- Text stays within boundaries during testing
- Symbol positioning consistent
- Transparency positioning accurate

---

## **Future Enhancements**

### **Advanced Features**

- Dynamic boundary adjustment during beta testing
- Smart symbol auto-sizing within boundaries
- Multi-language text boundary adaptation
- Performance optimization for transparency operations

### **Tooling**

- Inkscape plugin for proper layer naming
- Boundary visualization tool for designers
- Template validation and preview system
- Performance profiling for transparency vs coordinates

```

```
