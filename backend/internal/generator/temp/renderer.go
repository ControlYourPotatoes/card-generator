package text

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"strings"
)

type Renderer struct {
	ctx     *freetype.Context
	img     draw.Image
	fonts   map[string]*truetype.Font
	details *TextDetails
}

func NewRenderer(img draw.Image) (*Renderer, error) {
	// Initialize fonts
	regularFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regular font: %w", err)
	}

	boldFont, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bold font: %w", err)
	}

	fonts := map[string]*truetype.Font{
		"regular": regularFont,
		"bold":    boldFont,
	}

	// Create and initialize freetype context
	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetDst(img)
	ctx.SetClip(img.Bounds())
	// Set a default font
	ctx.SetFont(regularFont)
	// Set a default size
	ctx.SetFontSize(12)
	// Set the source color for drawing
	ctx.SetSrc(image.NewUniform(color.Black))

	return &Renderer{
		ctx:   ctx,
		img:   img,
		fonts: fonts,
	}, nil
}

// SetTextDetails sets the text content for rendering
func (r *Renderer) SetTextDetails(details *TextDetails) {
	r.details = details
}

func (r *Renderer) RenderTitle(bounds TextBounds) error {
	if r.details == nil {
		return fmt.Errorf("text details not set")
	}

	font := r.fonts["bold"]
	r.ctx.SetFont(font)
	r.ctx.SetFontSize(bounds.Style.FontSize)

	// Calculate approximate text width
	textWidth := fixed.Int26_6(bounds.Style.FontSize * 64 * float64(len(r.details.Title.Text)) / 2)

	// Center the text
	x := bounds.Rect.Min.X + (bounds.Rect.Dx()-textWidth.Round())/2
	y := bounds.Rect.Min.Y + int(bounds.Style.FontSize)

	pt := freetype.Pt(x, y)
	_, err := r.ctx.DrawString(r.details.Title.Text, pt)
	return err
}

func (r *Renderer) RenderEffect(bounds TextBounds) error {
	if r.details == nil {
		return fmt.Errorf("text details not set")
	}

	font := r.fonts["regular"]
	r.ctx.SetFont(font)
	r.ctx.SetFontSize(bounds.Style.FontSize)
	r.ctx.SetSrc(image.NewUniform(color.Black)) // Set default color

	// Handle keywords if present
	y := bounds.Rect.Min.Y + int(bounds.Style.FontSize)
	if len(r.details.Effect.Keywords) > 0 {
		r.ctx.SetFont(r.fonts["bold"]) // Use bold for keywords
		keywordText := strings.Join(r.details.Effect.Keywords, " ")
		pt := freetype.Pt(bounds.Rect.Min.X, y)
		if _, err := r.ctx.DrawString(keywordText, pt); err != nil {
			return err
		}
		y += int(bounds.Style.FontSize * 1.5) // Add space after keywords
	}

	// Reset font for effect text
	r.ctx.SetFont(font)

	// Split text into lines based on width
	lines := r.wrapText(r.details.Effect.Text, bounds.Rect.Dx(), font, bounds.Style.FontSize)

	for _, line := range lines {
		pt := freetype.Pt(bounds.Rect.Min.X, y)
		_, err := r.ctx.DrawString(line, pt)
		if err != nil {
			return err
		}
		y += int(bounds.Style.FontSize * 1.2) // 1.2 line spacing
	}

	return nil
}

func (r *Renderer) RenderStats(bounds TextBounds) error {
	if r.details == nil {
		return fmt.Errorf("text details not set")
	}

	font := r.fonts["bold"]
	r.ctx.SetFont(font)
	r.ctx.SetFontSize(bounds.Style.FontSize)

	// Render both power and toughness if present
	if r.details.Stats.Power != "" && r.details.Stats.Toughness != "" {
		statText := r.details.Stats.Power + "/" + r.details.Stats.Toughness
		textWidth := fixed.Int26_6(bounds.Style.FontSize * 64 * float64(len(statText)) / 2)
		x := bounds.Rect.Min.X + (bounds.Rect.Dx()-textWidth.Round())/2
		y := bounds.Rect.Min.Y + int(bounds.Style.FontSize)
		pt := freetype.Pt(x, y)
		_, err := r.ctx.DrawString(statText, pt)
		return err
	}
	return nil
}

func (r *Renderer) RenderCost(bounds TextBounds) error {
	if r.details == nil {
		return fmt.Errorf("text details not set")
	}

	font := r.fonts["bold"]
	r.ctx.SetFont(font)
	r.ctx.SetFontSize(bounds.Style.FontSize)

	pt := freetype.Pt(bounds.Rect.Min.X, bounds.Rect.Min.Y+int(bounds.Style.FontSize))
	_, err := r.ctx.DrawString(r.details.Title.Cost.Value, pt)
	return err
}

// Helper function to wrap text
func (r *Renderer) wrapText(text string, maxWidth int, font *truetype.Font, fontSize float64) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return lines
	}

	currentLine := words[0]
	approxCharWidth := fontSize * 0.6 // Approximate width of a character

	for _, word := range words[1:] {
		// Calculate approximate width of current line + space + new word
		nextLine := currentLine + " " + word
		approxWidth := len(nextLine) * int(approxCharWidth)

		if approxWidth <= maxWidth {
			currentLine = nextLine
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}
