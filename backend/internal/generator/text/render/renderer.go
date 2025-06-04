package render

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/internal/generator/text/manager"
	"github.com/ControlYourPotatoes/card-generator/internal/generator/text/types"
)

type Renderer struct {
	context   *gg.Context
	fontMgr   *manager.FontManager
	layoutMgr *manager.LayoutManager
	styleMgr  *manager.StyleManager
}

func NewRenderer(img image.Image) (*Renderer, error) {
	bounds := img.Bounds()
	dc := gg.NewContext(bounds.Dx(), bounds.Dy())
	dc.DrawImage(img, 0, 0)

	fontMgr, err := manager.NewFontManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize font manager: %w", err)
	}

	layoutMgr := manager.NewLayoutManager(bounds)
	styleMgr := manager.NewStyleManager()

	return &Renderer{
		context:   dc,
		fontMgr:   fontMgr,
		layoutMgr: layoutMgr,
		styleMgr:  styleMgr,
	}, nil
}

// RenderElement renders a single card element
func (r *Renderer) RenderElement(element types.CardElement, text string) error {
	// Get configurations from managers
	config := r.layoutMgr.GetTextConfiguration(element)
	style := r.styleMgr.GetStyle(element)
	fontPath := r.fontMgr.GetFontPath(style.FontName)

	// Calculate best font size that fits the bounds
	bestSize := r.fitTextSize(text, config.Bounds, style, fontPath)

	// Load font face
	if err := r.context.LoadFontFace(fontPath, bestSize); err != nil {
		return fmt.Errorf("failed to load font: %w", err)
	}

	// Set text color
	r.context.SetColor(style.Color)

	// Draw text based on configuration
	if style.SingleLine {
		return r.renderSingleLine(text, config.Bounds, style)
	}
	return r.renderMultiLine(text, config.Bounds, style)
}

func (r *Renderer) fitTextSize(text string, bounds image.Rectangle, style manager.Style, fontPath string) float64 {
	for size := style.Size; size >= style.MinFontSize; size-- {
		if err := r.context.LoadFontFace(fontPath, size); err != nil {
			continue
		}

		if style.SingleLine {
			if r.fitsInSingleLine(text, bounds) {
				return size
			}
		} else {
			if r.fitsInMultiLine(text, bounds, style.LineSpacing) {
				return size
			}
		}
	}
	return style.MinFontSize
}

func (r *Renderer) fitsInSingleLine(text string, bounds image.Rectangle) bool {
	width, height := r.context.MeasureString(text)
	return int(width) <= bounds.Dx() && int(height) <= bounds.Dy()
}

func (r *Renderer) fitsInMultiLine(text string, bounds image.Rectangle, lineSpacing float64) bool {
	lines := r.measureWrappedText(text, float64(bounds.Dx()))
	fontSize := r.context.FontSize()
	totalHeight := float64(len(lines)) * fontSize * lineSpacing
	return int(totalHeight) <= bounds.Dy()
}

func (r *Renderer) renderSingleLine(text string, bounds image.Rectangle, style manager.Style) error {
	width, height := r.context.MeasureString(text)
	x := r.calculateX(bounds, width, style.Alignment)
	y := r.calculateY(bounds, height, style.AnchorY)

	r.context.DrawString(text, x, y)
	return nil
}

func (r *Renderer) renderMultiLine(text string, bounds image.Rectangle, style manager.Style) error {
	r.context.DrawStringWrapped(
		text,
		float64(bounds.Min.X),
		float64(bounds.Min.Y),
		style.AnchorX,
		style.AnchorY,
		float64(bounds.Dx()),
		style.LineSpacing,
		style.Alignment,
	)
	return nil
}

func (r *Renderer) measureWrappedText(text string, maxWidth float64) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return lines
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		width, _ := r.context.MeasureString(currentLine + " " + word)
		if width <= maxWidth {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)
	return lines
}

func (r *Renderer) calculateX(bounds image.Rectangle, width float64, align gg.Align) float64 {
	switch align {
	case gg.AlignCenter:
		return float64(bounds.Min.X) + (float64(bounds.Dx())-width)/2
	case gg.AlignRight:
		return float64(bounds.Max.X) - width
	default: // gg.AlignLeft
		return float64(bounds.Min.X)
	}
}

func (r *Renderer) calculateY(bounds image.Rectangle, height float64, anchorY float64) float64 {
	boundsHeight := float64(bounds.Dy())
	y := float64(bounds.Min.Y) + (boundsHeight-height)*anchorY + height
	return y
}

// GetImage returns the final rendered image
func (r *Renderer) GetImage() image.Image {
	return r.context.Image()
}
