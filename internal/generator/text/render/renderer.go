// render/renderer.go
package render

import (
    "fmt"
    "image"
    "image/color"
    "github.com/fogleman/gg"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/text/types"
)

type Renderer struct {
    context *gg.Context
    fontMgr *FontManager
}

func NewRenderer(img image.Image) (*Renderer, error) {
    bounds := img.Bounds()
    dc := gg.NewContext(bounds.Dx(), bounds.Dy())
    dc.DrawImage(img, 0, 0)
    
    fontMgr, err := NewFontManager()
    if err != nil {
        return nil, err
    }
    
    return &Renderer{
        context: dc,
        fontMgr: fontMgr,
    }, nil
}

// RenderElement renders text for a specific card element
func (r *Renderer) RenderElement(element types.CardElement, text string) error {
    // Get element configuration
    config, exists := types.DefaultElements[element]
    if !exists {
        return fmt.Errorf("unknown element type: %s", element)
    }

    // Get appropriate font
    fontPath := r.fontMgr.GetFontPath(config.Font)

    // Calculate best font size
    bestSize := r.fitTextSize(text, config, fontPath)

    // Load font with calculated size
    if err := r.context.LoadFontFace(fontPath, bestSize); err != nil {
        return fmt.Errorf("failed to load font: %w", err)
    }

    // Set color
    r.context.SetColor(config.Color)

    // Calculate position based on alignment
    x, y := r.calculatePosition(text, config, bestSize)

    // Draw text
    if config.OneLine {
        r.context.DrawString(text, x, y)
    } else {
        r.context.DrawStringWrapped(
            text,
            x, y,
            0.0, 0.0,
            float64(config.Width),
            1.2, // line height
            convertAlignment(config.Align),
        )
    }

    return nil
}

// fitTextSize finds the largest font size that fits within bounds
func (r *Renderer) fitTextSize(text string, config types.TextElement, fontPath string) float64 {
    for size := config.MaxSize; size >= config.MinSize; size-- {
        if err := r.context.LoadFontFace(fontPath, size); err != nil {
            continue
        }

        width, height := r.context.MeasureString(text)
        
        if config.OneLine {
            if int(width) <= config.Width && int(height) <= config.Height {
                return size
            }
        } else {
            // For multi-line text, measure wrapped
            lines := r.measureWrappedText(text, float64(config.Width), size)
            totalHeight := float64(len(lines)) * size * 1.2
            if totalHeight <= float64(config.Height) {
                return size
            }
        }
    }
    
    return config.MinSize
}

// measureWrappedText returns lines after wrapping
func (r *Renderer) measureWrappedText(text string, maxWidth float64, fontSize float64) []string {
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

// calculatePosition determines x,y coordinates based on alignment
func (r *Renderer) calculatePosition(text string, config types.TextElement, fontSize float64) (float64, float64) {
    width, height := r.context.MeasureString(text)
    
    // Calculate X position
    var x float64
    switch config.Align {
    case "center":
        x = float64(config.X) + (float64(config.Width)-width)/2
    case "right":
        x = float64(config.X + config.Width) - width
    default: // left
        x = float64(config.X)
    }
    
    // Calculate Y position
    var y float64
    switch config.VerticalAlign {
    case "center":
        y = float64(config.Y) + (float64(config.Height)+height)/2
    case "bottom":
        y = float64(config.Y + config.Height)
    default: // top
        y = float64(config.Y) + height
    }
    
    return x, y
}

func convertAlignment(align string) gg.Align {
    switch align {
    case "center":
        return gg.AlignCenter
    case "right":
        return gg.AlignRight
    default:
        return gg.AlignLeft
    }
}