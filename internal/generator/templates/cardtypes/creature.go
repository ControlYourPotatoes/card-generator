package cardtypes

import (
    "image"
    "path/filepath"
    "strings"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

type CreatureTemplate struct {
    framesPath string
    artBounds  image.Rectangle
}

func NewCreatureTemplate() (*CreatureTemplate, error) {
    return &CreatureTemplate{
        framesPath: filepath.Join("internal", "generator", "templates", "images"),
        artBounds: image.Rect(170, 240, 1330, 1000), // Default art bounds
    }, nil
}

func (t *CreatureTemplate) GetFrame(data *card.CardData) (image.Image, error) {
    framePath := filepath.Join(t.framesPath, "BaseCreature.png")
    if t.isSpecialFrame(data) {
        framePath = filepath.Join(t.framesPath, "SpecialCreatureWithStats.png")
    }
    
    // Load the frame (this function needs to be accessible)
    return LoadFrame(framePath)
}

func (t *CreatureTemplate) GetTextBounds(data *card.CardData) *layout.TextBounds {
    bounds := layout.GetDefaultBounds()
    
    // Add stats positioning for creatures
    bounds.Stats = &layout.StatsConfig{
        Left: layout.TextConfig{
            Bounds:    image.Rect(130, 1820, 230, 1900),
            FontSize:  72,
            Alignment: "center",
        },
        Right: layout.TextConfig{
            Bounds:    image.Rect(1270, 1820, 1370, 1900),
            FontSize:  72,
            Alignment: "center",
        },
    }
    
    // Adjust effect box to account for stats
    bounds.Effect.Bounds = image.Rect(160, 1250, 1340, 1750)
    return bounds
}

func (t *CreatureTemplate) GetArtBounds() image.Rectangle {
    return t.artBounds
}

func (t *CreatureTemplate) isSpecialFrame(data *card.CardData) bool {
    // Check for special traits or keywords that would trigger special frame
    specialTraits := []string{"Legendary", "Ancient", "Divine"}
    for _, trait := range specialTraits {
        if strings.Contains(data.Trait, trait) {
            return true
        }
    }
    
    // Check for special keywords in effect text
    specialKeywords := []string{"CRITICAL", "DOUBLE STRIKE", "INDESTRUCTIBLE"}
    for _, keyword := range specialKeywords {
        if strings.Contains(data.Effect, keyword) {
            return true
        }
    }
    
    return false
}

