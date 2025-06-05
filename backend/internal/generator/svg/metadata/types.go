// svg/metadata/types.go
package metadata

import (
	"image"
)

// SVGMetadata contains game-ready metadata for SVG cards
type SVGMetadata struct {
	CardID           string                     `json:"card_id"`
	InteractiveZones []InteractiveZone          `json:"interactive_zones"`
	AnimationTargets []string                   `json:"animation_targets"`
	GameState        map[string]string          `json:"game_state"`
	Version          string                     `json:"version"`
	GeneratedAt      string                     `json:"generated_at"`
}

// InteractiveZone defines clickable/hoverable areas on the card for game interaction
type InteractiveZone struct {
	ID      string          `json:"id"`
	Bounds  image.Rectangle `json:"bounds"`
	Action  string          `json:"action"`  // "tap", "target", "inspect", etc.
	Trigger string          `json:"trigger"` // "click", "hover", "contextmenu", etc.
	Data    map[string]interface{} `json:"data,omitempty"` // Additional zone-specific data
}

// AnimationTarget identifies elements that can be animated for game effects
type AnimationTarget struct {
	ElementID     string                 `json:"element_id"`
	AnimationType string                 `json:"animation_type"` // "glow", "shake", "pulse", "rotate", etc.
	Properties    map[string]interface{} `json:"properties"`
	Duration      string                 `json:"duration,omitempty"`      // CSS duration format
	Trigger       string                 `json:"trigger,omitempty"`       // What triggers the animation
}

// SVGBounds represents coordinate boundaries within the SVG coordinate system
type SVGBounds struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// TextPlacement defines how text should be positioned and styled in SVG
type TextPlacement struct {
	ElementID   string            `json:"element_id"`
	X           float64           `json:"x"`
	Y           float64           `json:"y"`
	Width       float64           `json:"width,omitempty"`
	Height      float64           `json:"height,omitempty"`
	TextAnchor  string            `json:"text_anchor,omitempty"`  // "start", "middle", "end"
	FontFamily  string            `json:"font_family,omitempty"`
	FontSize    float64           `json:"font_size,omitempty"`
	FontWeight  string            `json:"font_weight,omitempty"`
	Fill        string            `json:"fill,omitempty"`
	Stroke      string            `json:"stroke,omitempty"`
	Classes     []string          `json:"classes,omitempty"`      // CSS classes to apply
	Attributes  map[string]string `json:"attributes,omitempty"`   // Additional SVG attributes
}

// GameStateClasses maps game states to CSS classes for dynamic styling
type GameStateClasses struct {
	Default  []string `json:"default"`
	Tapped   []string `json:"tapped,omitempty"`
	Selected []string `json:"selected,omitempty"`
	Targeted []string `json:"targeted,omitempty"`
	Disabled []string `json:"disabled,omitempty"`
	Summoned []string `json:"summoned,omitempty"`
	Custom   map[string][]string `json:"custom,omitempty"`
} 