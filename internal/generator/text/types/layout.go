// types/layout.go
package types

import (
    "image"
    "image/color"
)

type TextStyle struct {
    FontSize    float64
    MinFontSize float64
    MaxFontSize float64
    Color       color.Color
    Alignment   string
    SingleLine  bool
    LineSpacing float64
}

type TextConfiguration struct {
    Bounds  image.Rectangle
    Style   TextStyle
    Element TextElement  // Reference to original template element
}

func DefaultTextStyle() TextStyle {
    return TextStyle{
        FontSize:    60,
        MinFontSize: 36,
        MaxFontSize: 72,
        Color:       color.White,
        Alignment:   "left",
        SingleLine:  false,
        LineSpacing: 1.2,
    }
}