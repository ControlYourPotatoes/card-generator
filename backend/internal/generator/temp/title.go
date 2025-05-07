package text

import (
    "fmt"
    "image"
    "strings"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

type titleProcessor struct {
    defaultStyle TextConfig
    costStyle    TextConfig
}

func NewTitleProcessor() TitleProcessor {
    return &titleProcessor{
        defaultStyle: TextConfig{
            FontSize:  72,
            Color:     "#000000",
            Bold:      true,
            Alignment: "center",
        },
        costStyle: TextConfig{
            FontSize:  72,
            Color:     "#FFFFFF",
            Bold:      true,
            Alignment: "center",
        },
    }
}

func (tp *titleProcessor) ProcessTitle(card card.Card) (TextBounds, error) {
    if err := tp.ValidateTitle(card.GetName()); err != nil {
        return TextBounds{}, err
    }

    return TextBounds{
        Rect: image.Rect(130, 100, 1370, 180),
        Style: tp.defaultStyle,
    }, nil
}

func (tp *titleProcessor) ProcessCost(cost int) (TextBounds, error) {
    // Special handling for X costs
    if cost < 0 {
        return TextBounds{
            Rect: image.Rect(1300, 90, 1360, 150),
            Style: tp.costStyle,
        }, nil
    }

    // Regular cost processing
    if cost > 99 {
        return TextBounds{}, fmt.Errorf("cost exceeds maximum value of 99")
    }

    return TextBounds{
        Rect: image.Rect(1300, 90, 1360, 150),
        Style: tp.costStyle,
    }, nil
}

func (tp *titleProcessor) ValidateTitle(title string) error {
    if title == "" {
        return fmt.Errorf("title cannot be empty")
    }

    if len(title) > 40 {
        return fmt.Errorf("title exceeds maximum length of 40 characters")
    }

    // Check for invalid characters
    invalidChars := "!@#$%^&*()_+=<>{}[]|"
    if strings.ContainsAny(title, invalidChars) {
        return fmt.Errorf("title contains invalid characters")
    }

    return nil
}