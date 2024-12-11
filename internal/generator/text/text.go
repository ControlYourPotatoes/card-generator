// internal/generator/image/text/text.go

package text

import (
    "image"
    "image/draw"

    "github.com/golang/freetype"
    "github.com/golang/freetype/truetype"
)

type TextRenderer struct {
    nameFont   *truetype.Font
    effectFont *truetype.Font
    statsFont  *truetype.Font
}

func NewTextRenderer() (*TextRenderer, error) {
    // Load fonts for different text elements
    // Could use different fonts for name, effect, and stats
    return &TextRenderer{}, nil
}

func (tr *TextRenderer) RenderName(img *image.RGBA, name string, bounds image.Rectangle) error {
    // Implement name rendering with appropriate font and size
    return nil
}

func (tr *TextRenderer) RenderEffect(img *image.RGBA, effect string, bounds image.Rectangle) error {
    // Implement effect text rendering with word wrap
    return nil
}

func (tr *TextRenderer) RenderStats(img *image.RGBA, attack, defense int, bounds StatsPosition) error {
    // Implement stats rendering
    return nil
} 