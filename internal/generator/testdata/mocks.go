// internal/generator/testdata/mocks.go
package testdata

import (
    "image"
    "image/color"
    "image/draw"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/templates"
)

// MockTemplate implements templates.Template interface for testing
type MockTemplate struct {
    templates.Template
    Frame     image.Image
    Bounds    *layout.TextBounds
    ArtBounds image.Rectangle
}

func NewMockTemplate() *MockTemplate {
    frame := image.NewRGBA(image.Rect(0, 0, 1500, 2100))
    bgColor := color.RGBA{R: 200, G: 200, B: 200, A: 255}
    draw.Draw(frame, frame.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
    
    return &MockTemplate{
        Frame:     frame,
        Bounds:    layout.GetDefaultBounds(),
        ArtBounds: image.Rect(170, 240, 1330, 1000),
    }
}

// ... implement Template interface methods

// MockArtProcessor implements art.ArtProcessor interface for testing
type MockArtProcessor struct {
    ProcessArtCalled bool
    LastBounds      image.Rectangle
}

func (m *MockArtProcessor) ProcessArt(data *card.CardData, bounds image.Rectangle) (image.Image, error) {
    m.ProcessArtCalled = true
    m.LastBounds = bounds
    img := image.NewRGBA(bounds)
    artColor := color.RGBA{R: 100, G: 150, B: 200, A: 255}
    draw.Draw(img, bounds, &image.Uniform{artColor}, image.Point{}, draw.Src)
    return img, nil
}

// MockTextProcessor implements text.TextProcessor interface for testing
type MockTextProcessor struct {
    RenderTextCalled bool
    LastBounds      *layout.TextBounds
}

func (m *MockTextProcessor) RenderText(img draw.Image, data *card.CardData, bounds *layout.TextBounds) error {
    m.RenderTextCalled = true
    m.LastBounds = bounds
    // Draw visible rectangles for text areas
    textColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}
    draw.Draw(img, bounds.Name.Bounds, &image.Uniform{textColor}, image.Point{}, draw.Over)
    draw.Draw(img, bounds.Effect.Bounds, &image.Uniform{textColor}, image.Point{}, draw.Over)
    if bounds.Stats != nil {
        draw.Draw(img, bounds.Stats.Left.Bounds, &image.Uniform{textColor}, image.Point{}, draw.Over)
        draw.Draw(img, bounds.Stats.Right.Bounds, &image.Uniform{textColor}, image.Point{}, draw.Over)
    }
    return nil
}