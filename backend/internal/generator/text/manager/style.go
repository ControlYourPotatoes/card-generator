package manager

import (
	"github.com/ControlYourPotatoes/card-generator/internal/generator/text/types"
	"github.com/fogleman/gg"
	"image/color"
)

type StyleManager struct {
	defaultStyles map[types.CardElement]Style
}

type Style struct {
	FontName    string
	Size        float64
	Color       color.Color
	Alignment   gg.Align
	LineSpacing float64
	AnchorX     float64
	AnchorY     float64
	Wrap        bool
}

func NewStyleManager() *StyleManager {
	return &StyleManager{
		defaultStyles: map[types.CardElement]Style{
			types.ElementTitle: {
				FontName:    "bold",
				Size:        72,
				Color:       color.Black,
				Alignment:   gg.AlignCenter,
				LineSpacing: 1.0,
				AnchorX:     0.5,
				AnchorY:     0.5,
				Wrap:        false,
			},
			types.ElementEffect: {
				FontName:    "regular",
				Size:        48,
				Color:       color.Black,
				Alignment:   gg.AlignLeft,
				LineSpacing: 1.5,
				AnchorX:     0.0,
				AnchorY:     0.0,
				Wrap:        true,
			},
			// Add other element styles...
		},
	}
}

func (sm *StyleManager) GetStyle(element types.CardElement) Style {
	if style, exists := sm.defaultStyles[element]; exists {
		return style
	}
	return sm.defaultStyles[types.ElementEffect] // Default fallback
}
