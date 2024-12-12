package templates

import (
    "fmt"
    "image"
    "image/png"
    "os"
    "path/filepath"
    "runtime"
)

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
    
    fmt.Printf("Template directory: %s\n", b.framesPath)
    fmt.Printf("Looking for frame at: %s\n", framePath)
    
    f, err := os.Open(framePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open frame: %w", err)
    }
    defer f.Close()

    return png.Decode(f)
}

// GetArtBounds returns the art bounds
func (b *BaseTemplate) GetArtBounds() image.Rectangle {
    return b.artBounds
}

// Helper functions
func getTemplateDir() string {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return ""
    }
    return filepath.Join(filepath.Dir(filepath.Dir(filename)), "images")
}

func GetDefaultArtBounds() image.Rectangle {
    return image.Rect(170, 240, 1330, 1000)
}