// svg/renderer/bounds.go
package renderer

import (
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/svg/metadata"
)

// BoundsCalculator handles coordinate conversion between image bounds and SVG coordinates
type BoundsCalculator struct {
	canvasWidth  float64
	canvasHeight float64
}

// NewBoundsCalculator creates a new bounds calculator
func NewBoundsCalculator(canvasWidth, canvasHeight float64) *BoundsCalculator {
	return &BoundsCalculator{
		canvasWidth:  canvasWidth,
		canvasHeight: canvasHeight,
	}
}

// ImageToSVG converts image.Rectangle to SVG coordinate system
func (c *BoundsCalculator) ImageToSVG(bounds image.Rectangle) metadata.SVGBounds {
	// TODO: Implement coordinate conversion in Phase 2
	return metadata.SVGBounds{
		X:      float64(bounds.Min.X),
		Y:      float64(bounds.Min.Y),
		Width:  float64(bounds.Dx()),
		Height: float64(bounds.Dy()),
	}
}

// SVGToImage converts SVG bounds back to image.Rectangle
func (c *BoundsCalculator) SVGToImage(bounds metadata.SVGBounds) image.Rectangle {
	// TODO: Implement coordinate conversion in Phase 2
	return image.Rect(
		int(bounds.X),
		int(bounds.Y),
		int(bounds.X+bounds.Width),
		int(bounds.Y+bounds.Height),
	)
} 