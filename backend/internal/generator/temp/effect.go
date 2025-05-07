package text

import (
    "fmt"
    "image"
    "strings"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

type effectProcessor struct {
    defaultStyle TextConfig
    keywordStyle TextConfig
}

func NewEffectProcessor() EffectProcessor {
    return &effectProcessor{
        defaultStyle: TextConfig{
            FontSize:  48,
            Color:     "#000000",
            Bold:      false,
            Alignment: "left",
        },
        keywordStyle: TextConfig{
            FontSize:  48,
            Color:     "#000000",
            Bold:      true,
            Alignment: "center",
        },
    }
}

func (ep *effectProcessor) ProcessEffect(effect string) (TextBounds, error) {
    if err := ep.ValidateEffect(effect); err != nil {
        return TextBounds{}, err
    }

    // Default effect text position
    return TextBounds{
        Rect: image.Rect(160, 1250, 1340, 1750),
        Style: ep.defaultStyle,
    }, nil
}

// ExtractKeywords now just returns the keywords from CardData
func (ep *effectProcessor) ExtractKeywords(effect string) []string {
    // This method is kept for interface compatibility
    // but should not be used for keyword extraction anymore
    // as keywords come from CardData.Keywords
    return nil
}

func (ep *effectProcessor) ValidateEffect(effect string) error {
    
    
    // Check for maximum length
    if len(effect) > 500 {
        return fmt.Errorf("effect text exceeds maximum length of 500 characters")
    }

    // Add any additional effect text validation rules
    // For example, check for proper sentence structure
    if !strings.HasSuffix(effect, ".") && !strings.HasSuffix(effect, "!") {
        return fmt.Errorf("effect text must end with proper punctuation")
    }

    return nil
}

// Additional helper methods for effect processing

// GetKeywordBounds returns the bounds for keyword text
func (ep *effectProcessor) GetKeywordBounds(hasKeywords bool) image.Rectangle {
    if hasKeywords {
        // Position for keywords above the main effect text
        return image.Rect(160, 1200, 1340, 1240)
    }
    return image.Rectangle{}
}

// GetEffectBounds returns adjusted bounds based on whether there are keywords
func (ep *effectProcessor) GetEffectBounds(hasKeywords bool) image.Rectangle {
    if hasKeywords {
        // Move effect text down if there are keywords
        return image.Rect(160, 1280, 1340, 1750)
    }
    return image.Rect(160, 1250, 1340, 1750)
}

// GetKeywordStyle returns the style for keyword text
func (ep *effectProcessor) GetKeywordStyle() TextConfig {
    return ep.keywordStyle
}

// ProcessFullEffect processes both keywords and effect text
func (ep *effectProcessor) ProcessFullEffect(cardData *card.CardData) (*struct {
    KeywordBounds TextBounds
    EffectBounds  TextBounds
}, error) {
    hasKeywords := len(cardData.Keywords) > 0

    return &struct {
        KeywordBounds TextBounds
        EffectBounds  TextBounds
    }{
        KeywordBounds: TextBounds{
            Rect:  ep.GetKeywordBounds(hasKeywords),
            Style: ep.keywordStyle,
        },
        EffectBounds: TextBounds{
            Rect:  ep.GetEffectBounds(hasKeywords),
            Style: ep.defaultStyle,
        },
    }, nil
}