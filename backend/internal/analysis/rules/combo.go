package rules

import (
    
    "github.com/ControlYourPotatoes/card-generator/internal/analysis/types"
)
// ComboPatterns defines patterns that indicate combo potential
var ComboPatterns = map[string][]types.Pattern{
    "INFINITE_COMBO": {
        {Value: "untap", Type: types.ProximityMatch, Proximity: 5},
        {Value: "return.*from.*graveyard", Type: types.RegexMatch},
    },
    "MANA_COMBO": {
        {Value: "add.*for each", Type: types.RegexMatch},
        {Value: "doubles", Type: types.ProximityMatch, Proximity: 3},
    },
    "TOKEN_COMBO": {
        {Value: "create.*token", Type: types.RegexMatch},
        {Value: "double.*tokens", Type: types.RegexMatch},
    },
}

// ComboTriggers defines common combo trigger words
var ComboTriggers = []string{
    "whenever",
    "at the beginning of",
    "at end of",
    "for each",
}

// ComboActions defines common combo action words
var ComboActions = []string{
    "create",
    "draw",
    "return",
    "copy",
    "search",
}

// ComboRule represents a specific combo detection rule
type ComboRule struct {
    Name        string
    Description string
    Patterns    []types.Pattern
    MinPieces   int
    Weight      int
}

// PredefinedCombos contains known powerful card combinations
var PredefinedCombos = []ComboRule{
    {
        Name:        "RESURRECTION_LOOP",
        Description: "Cards that can repeatedly resurrect creatures",
        Patterns: []types.Pattern{
            {Value: "return.*from.*graveyard", Type: types.RegexMatch},
            {Value: "when.*dies", Type: types.RegexMatch},
        },
        MinPieces: 2,
        Weight:    3,
    },
    {
        Name:        "TOKEN_DOUBLING",
        Description: "Cards that can exponentially increase tokens",
        Patterns: []types.Pattern{
            {Value: "create.*token", Type: types.RegexMatch},
            {Value: "double.*tokens", Type: types.RegexMatch},
        },
        MinPieces: 2,
        Weight:    3,
    },
}

// DetectComboPotential checks if a card has combo potential
func DetectComboPotential(effect string) []types.Tag {
    var tags []types.Tag
    
    // Check for combo triggers
    hasTrigger := false
    for _, trigger := range ComboTriggers {
        if contains(effect, trigger) {
            hasTrigger = true
            break
        }
    }
    
    // Check for combo actions
    hasAction := false
    for _, action := range ComboActions {
        if contains(effect, action) {
            hasAction = true
            break
        }
    }
    
    // If card has both trigger and action, it has combo potential
    if hasTrigger && hasAction {
        tags = append(tags, types.Tag{
            Name:     "COMBO_POTENTIAL",
            Category: types.TagCombo,
            Weight:   2,
        })
    }
    
    // Check predefined combo patterns
    for _, combo := range PredefinedCombos {
        matchesAll := true
        for _, pattern := range combo.Patterns {
            if !matchPattern(effect, pattern) {
                matchesAll = false
                break
            }
        }
        if matchesAll {
            tags = append(tags, types.Tag{
                Name:     combo.Name,
                Category: types.TagCombo,
                Weight:   combo.Weight,
            })
        }
    }
    
    return tags
}

// Helper functions
func contains(text, substring string) bool {
    return strings.Contains(strings.ToLower(text), strings.ToLower(substring))
}

func matchPattern(text string, pattern types.Pattern) bool {
    // Pattern matching implementation would go here
    // This is a simplified version
    return contains(text, pattern.Value)
}