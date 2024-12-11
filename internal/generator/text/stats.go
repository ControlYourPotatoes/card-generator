package text

import (
    "fmt"
    "image"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

type statsProcessor struct {
    defaultStyle TextConfig
}

func NewStatsProcessor() StatsProcessor {
    return &statsProcessor{
        defaultStyle: TextConfig{
            FontSize:  72,
            Color:     "#000000",
            Bold:      true,
            Alignment: "center",
        },
    }
}

func (sp *statsProcessor) ProcessStats(c card.Card) (TextBounds, error) {
    if err := sp.ValidateStats(c); err != nil {
        return TextBounds{}, err
    }

    bounds := TextBounds{
        Rect: image.Rect(130, 1820, 1370, 1900),
        Style: sp.defaultStyle,
    }

    // Adjust bounds based on card type
    switch cardData := c.(type) {
    case *card.Creature:
        if cardData.Trait != "" {
            // Adjust bounds to accommodate trait
            bounds.Rect.Max.Y += 30
        }
    }

    return bounds, nil
}

func (sp *statsProcessor) ValidateStats(c card.Card) error {
    switch cardData := c.(type) {
    case *card.Creature:
        if cardData.Attack < 0 {
            return fmt.Errorf("attack cannot be negative")
        }
        if cardData.Defense < 0 {
            return fmt.Errorf("defense cannot be negative")
        }
        if cardData.Attack > 99 || cardData.Defense > 99 {
            return fmt.Errorf("stats cannot exceed 99")
        }
        if cardData.Trait != "" {
            if len(cardData.Trait) > 30 {
                return fmt.Errorf("trait exceeds maximum length of 30 characters")
            }
            // Validate trait format
            validTraits := map[string]bool{
                "Beast":     true,
                "Warrior":   true,
                "Dragon":    true,
                "Demon":     true,
                "Angel":     true,
                "Legendary": true,
                "Ancient":   true,
                "Divine":    true,
            }
            if !validTraits[cardData.Trait] {
                return fmt.Errorf("invalid trait: %s", cardData.Trait)
            }
        }
    }
    return nil
}