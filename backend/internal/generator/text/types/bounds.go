// types/bounds.go
package types

import (
	"image"
	"image/color"
)

// TextElement represents a card text element with positioning from the template
type TextElement struct {
	Name          string
	X             int
	Y             int
	Width         int
	Height        int
	DefaultSize   float64
	MinSize       float64
	MaxSize       float64
	Color         color.Color
	OneLine       bool
	VerticalAlign string
	Align         string
	Font          string
}

// Standard card text elements based on CardConjurer template
var DefaultElements = map[CardElement]TextElement{
	ElementTitle: {
		Name:          "Title",
		X:             125,
		Y:             90,
		Width:         1250,
		Height:        80,
		DefaultSize:   80,
		MinSize:       48,
		MaxSize:       80,
		Color:         color.White,
		OneLine:       true,
		VerticalAlign: "center",
		Font:          "regular",
	},
	ElementEffect: {
		Name:          "Rules",
		X:             160,
		Y:             1300,
		Width:         1180,
		Height:        500,
		DefaultSize:   60,
		MinSize:       36,
		MaxSize:       60,
		Color:         color.White,
		OneLine:       false,
		VerticalAlign: "center",
		Font:          "regular",
	},
	ElementType: {
		Name:          "Type",
		X:             125,
		Y:             1885,
		Width:         1250,
		Height:        70,
		DefaultSize:   70,
		MinSize:       40,
		MaxSize:       70,
		Color:         color.White,
		OneLine:       true,
		VerticalAlign: "center",
		Align:         "center",
		Font:          "regular",
	},
	ElementStats: {
		Name:          "Stats",
		X:             115,
		Y:             1885,
		Width:         165,
		Height:        70,
		DefaultSize:   70,
		MinSize:       40,
		MaxSize:       70,
		Color:         color.White,
		OneLine:       true,
		VerticalAlign: "center",
		Align:         "center",
		Font:          "regular",
	},
}

// GetBounds returns the rectangle for the element
func (te TextElement) GetBounds() image.Rectangle {
	return image.Rect(te.X, te.Y, te.X+te.Width, te.Y+te.Height)
}
