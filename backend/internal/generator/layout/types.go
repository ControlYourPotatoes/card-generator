// File: internal/generator/layout/types.go
package layout

import (
	"image"
)

// TextBounds defines positions for text elements
type TextBounds struct {
	Name          TextConfig
	ManaCost      TextConfig
	Type          TextConfig
	Effect        TextConfig
	Stats         *StatsConfig // Only for creatures
	CollectorInfo TextConfig
}

// TextConfig holds positioning and styling data
type TextConfig struct {
	Bounds    image.Rectangle
	FontSize  float64
	Alignment string // "left", "center", "right"
}

// StatsConfig holds creature stats positioning
type StatsConfig struct {
	Left  TextConfig // Power/Attack
	Right TextConfig // Defense/Toughness
}

// Get default text bounds used by most cards
func GetDefaultBounds() *TextBounds {
	return &TextBounds{
		Name: TextConfig{
			Bounds:    image.Rect(125, 90, 1375, 170),
			FontSize:  80,
			Alignment: "left",
		},
		ManaCost: TextConfig{
			Bounds:    image.Rect(125, 90, 1375, 170),
			FontSize:  80,
			Alignment: "right",
		},
		Type: TextConfig{
			Bounds:    image.Rect(125, 1885, 1375, 1955),
			FontSize:  70,
			Alignment: "center",
		},
		Effect: TextConfig{
			Bounds:    image.Rect(160, 1300, 1340, 1800),
			FontSize:  60,
			Alignment: "left",
		},
		CollectorInfo: TextConfig{
			Bounds:    image.Rect(110, 2010, 750, 2090),
			FontSize:  32,
			Alignment: "left",
		},
	}
}
