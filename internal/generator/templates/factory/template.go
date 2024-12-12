package templates

import (
    "fmt"
    "image"
    "image/png"
    "os"
    "path/filepath"
    "runtime"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

// Template defines what each card template must provide
type Template interface {
    // GetFrame should load and return the appropriate frame image
    GetFrame(data *card.CardData) (image.Image, error)
    GetTextBounds(data *card.CardData) *layout.TextBounds
    GetArtBounds() image.Rectangle
}

// BaseTemplate provides common template functionality
type BaseTemplate struct {
    framesPath string
    artBounds  image.Rectangle
}

func NewBaseTemplate() *BaseTemplate {
    return &BaseTemplate{
        framesPath: getTemplateDir(),
        artBounds:  GetDefaultArtBounds(),
    }
}

// LoadFrame is a helper function for loading frame images
func (b *BaseTemplate) LoadFrame(imageName string) (image.Image, error) {
    framePath := filepath.Join(b.framesPath, imageName)
    
    // Debug logging
    fmt.Printf("Template directory: %s\n", b.framesPath)
    fmt.Printf("Looking for frame at: %s\n", framePath)
    
    f, err := os.Open(framePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open frame: %w", err)
    }
    defer f.Close()

    return png.Decode(f)
}

// getTemplateDir returns the absolute path to templates directory
func getTemplateDir() string {
    // Get the current file's location
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return ""
    }
    
    // Navigate to the images directory
    return filepath.Join(filepath.Dir(filename), "images")
}

// GetDefaultArtBounds returns default art placement bounds
func GetDefaultArtBounds() image.Rectangle {
    return image.Rect(170, 240, 1330, 1000)
}