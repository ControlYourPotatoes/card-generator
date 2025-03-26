package rules

import (
    
    "github.com/ControlYourPotatoes/card-generator/internal/analysis/types"
    "github.com/ControlYourPotatoes/card-generator/internal/core/common"
)

// TribalTypes defines all supported tribal types
var TribalRules = []types.TagRule{
    {
        Name:     "TRIBAL_LORD",
        Category: types.TagTribal,
        Patterns: []types.Pattern{
            {Value: `other .* you control get`, Type: types.RegexMatch},
            {Value: `creatures you control of the chosen type`, Type: types.ExactMatch},
        },
        Weight:      3,
        Description: "Card boosts creatures of a specific type",
    },
    {
        Name:     "TRIBAL_SYNERGY",
        Category: types.TagTribal,
        Patterns: []types.Pattern{
            {Value: `whenever .* (you control|enters)`, Type: types.RegexMatch},
            {Value: `for each .* you control`, Type: types.RegexMatch},
        },
        Weight:      2,
        Description: "Card synergizes with specific creature types",
    },
}

// TribalEffectPatterns maps tribal types to their common effect patterns
var TribalEffectPatterns = map[common.Tribe][]types.Pattern{
    common.TribeZombie: {
        {Value: "create.*Zombie", Type: types.RegexMatch},
        {Value: "return.*from.*graveyard", Type: types.RegexMatch},
    },
    common.TribeVampire: {
        {Value: "gain.*life", Type: types.RegexMatch},
        {Value: "drain", Type: types.ExactMatch},
    },
    common.TribeGoblin: {
        {Value: "haste", Type: types.ExactMatch},
        {Value: "attack", Type: types.ExactMatch},
    },
    common.TribeDemon: {
        {Value: "sacrifice", Type: types.ExactMatch},
        {Value: "deal.*damage", Type: types.RegexMatch},
    },
    // Add patterns for other tribes as needed
}


// GenerateTribalTags generates tribal-specific tags
func GenerateTribalTags(cardType string, effect string, tribes []string) []types.Tag {
    var tags []types.Tag
    
    // Check if card is of a tribal type
    for _, tribalType := range TribalTypes {
        if cardType == tribalType {
            tags = append(tags, types.Tag{
                Name:     tribalType + "_TRIBAL",
                Category: types.TagTribal,
                Weight:   2,
            })
        }
    }
    
    // Check for tribal effect patterns
    for tribalType, patterns := range TribalEffectPatterns {
        for _, pattern := range patterns {
            // Pattern matching implementation would go here
            // For brevity, using simple contains check
            if pattern.Type == types.ExactMatch && effect == pattern.Value {
                tags = append(tags, types.Tag{
                    Name:     tribalType + "_SYNERGY",
                    Category: types.TagSynergy,
                    Weight:   1,
                })
            }
        }
    }
    
    return tags
}