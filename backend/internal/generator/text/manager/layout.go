// manager/layout.go
package manager

import (
	"github.com/ControlYourPotatoes/card-generator/internal/generator/text/types"
	"image"
)

type LayoutManager struct {
	cardBounds     image.Rectangle
	templateData   map[types.CardElement]types.TextElement
	customBounds   map[types.CardElement]image.Rectangle
	styleOverrides map[types.CardElement]types.TextStyle
}

func NewLayoutManager(cardBounds image.Rectangle) *LayoutManager {
	return &LayoutManager{
		cardBounds:     cardBounds,
		templateData:   types.DefaultElements,
		customBounds:   make(map[types.CardElement]image.Rectangle),
		styleOverrides: make(map[types.CardElement]types.TextStyle),
	}
}

// GetElementBounds returns the bounds for a specific element
func (lm *LayoutManager) GetElementBounds(element types.CardElement) image.Rectangle {
	// Check for custom bounds first
	if bounds, exists := lm.customBounds[element]; exists {
		return bounds
	}

	// Fall back to template data
	if elemData, exists := lm.templateData[element]; exists {
		return elemData.GetBounds()
	}

	// Return zero rectangle if element not found
	return image.Rectangle{}
}

// GetElementStyle returns the style configuration for an element
func (lm *LayoutManager) GetElementStyle(element types.CardElement) types.TextStyle {
	// Check for style override
	if style, exists := lm.styleOverrides[element]; exists {
		return style
	}

	// Get default element data
	if elemData, exists := lm.templateData[element]; exists {
		return types.TextStyle{
			FontSize:    elemData.DefaultSize,
			MinFontSize: elemData.MinSize,
			MaxFontSize: elemData.MaxSize,
			Color:       elemData.Color,
			Alignment:   elemData.Align,
			SingleLine:  elemData.OneLine,
		}
	}

	// Return default style if element not found
	return types.DefaultTextStyle()
}

// AdjustForKeywords adjusts effect text position when keywords are present
func (lm *LayoutManager) AdjustForKeywords(hasKeywords bool) {
	if !hasKeywords {
		return
	}

	if effectElem, exists := lm.templateData[types.ElementEffect]; exists {
		// Move effect text down and reduce height to accommodate keywords
		effectElem.Y += 50
		effectElem.Height -= 50
		lm.templateData[types.ElementEffect] = effectElem
	}
}

// AdjustForCardType modifies layout based on card type
func (lm *LayoutManager) AdjustForCardType(cardType string) {
	switch cardType {
	case "Creature":
		// Ensure stats areas are visible
		lm.templateData[types.ElementStats] = types.DefaultElements[types.ElementStats]

		// Adjust effect text height to accommodate stats
		if effectElem, exists := lm.templateData[types.ElementEffect]; exists {
			effectElem.Height -= 70 // Make room for stats
			lm.templateData[types.ElementEffect] = effectElem
		}

	case "Spell", "Artifact", "Incantation", "Anthem":
		// Hide stats areas
		delete(lm.templateData, types.ElementStats)

		// Allow effect text to use full height
		if effectElem, exists := lm.templateData[types.ElementEffect]; exists {
			effectElem.Height += 70 // Use space normally reserved for stats
			lm.templateData[types.ElementEffect] = effectElem
		}
	}
}

// SetCustomBounds allows overriding default bounds for an element
func (lm *LayoutManager) SetCustomBounds(element types.CardElement, bounds image.Rectangle) {
	lm.customBounds[element] = bounds
}

// SetStyleOverride allows overriding default style for an element
func (lm *LayoutManager) SetStyleOverride(element types.CardElement, style types.TextStyle) {
	lm.styleOverrides[element] = style
}

// GetTextConfiguration returns combined bounds and style info for rendering
func (lm *LayoutManager) GetTextConfiguration(element types.CardElement) types.TextConfiguration {
	bounds := lm.GetElementBounds(element)
	style := lm.GetElementStyle(element)

	return types.TextConfiguration{
		Bounds:  bounds,
		Style:   style,
		Element: lm.templateData[element],
	}
}
