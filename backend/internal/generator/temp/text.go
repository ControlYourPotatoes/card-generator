// text.go
package text

import (
    "image"
    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

// TextConfig defines the styling configuration for text elements
type TextConfig struct {
    FontSize  float64
    Color     string
    Bold      bool
    Alignment string
}

// TextDetails contains all text-related information for a card
type TextDetails struct {
    Title struct {
        Text      string
        Position  image.Rectangle
        Style     TextConfig
        Cost      CostInfo
    }
    Effect struct {
        Keywords     []string
        Text        string
        Position    image.Rectangle
        Style       TextConfig
    }
    Stats struct {
        CardType    string
        Subtype     string
        Power       string
        Toughness   string
        Position    image.Rectangle
        Style       TextConfig
    }
    IsSpecialFrame bool
    FramePath     string
}

type CostInfo struct {
    Value     string
    Position  image.Point
    IsXCost   bool
    Style     TextConfig
}


// TextBounds holds positioning information
type TextBounds struct {
    Rect     image.Rectangle
    Style    TextConfig
}

// TitleProcessor handles title and cost text processing
type TitleProcessor interface {
    ProcessTitle(card card.Card) (TextBounds, error)
    ProcessCost(cost int) (TextBounds, error)
    ValidateTitle(title string) error
}

// EffectProcessor handles effect text processing
type EffectProcessor interface {
    ProcessEffect(effect string) (TextBounds, error)
    ExtractKeywords(effect string) []string
    ValidateEffect(effect string) error
}

// StatsProcessor handles stats text processing
type StatsProcessor interface {
    ProcessStats(card card.Card) (TextBounds, error)
    ValidateStats(card card.Card) error
}

// TextRenderer handles the final rendering of all text elements
type TextRenderer interface {
    RenderTitle(bounds TextBounds) error
    RenderEffect(bounds TextBounds) error
    RenderStats(bounds TextBounds) error
    RenderCost(bounds TextBounds) error
}

// TextOutput represents the final processed text information
type TextOutput struct {
    Title    TextBounds
    Cost     TextBounds
    Effect   TextBounds
    Stats    TextBounds
    Keywords []string
}