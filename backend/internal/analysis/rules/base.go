package rules

import (
	"github.com/ControlYourPotatoes/card-generator/internal/analysis/types"
)

// BaseRules contains the fundamental rules for tag generation
var BaseRules = []types.TagRule{
	{
		Name:     "TOKEN_GENERATOR",
		Category: types.TagMechanic,
		Patterns: []types.Pattern{
			{Value: `create (\d+|a|an) .* token`, Type: types.RegexMatch},
			{Value: `creates? .* token`, Type: types.RegexMatch},
		},
		Weight:      2,
		Description: "Card creates one or more tokens",
	},
	{
		Name:     "GRAVEYARD_INTERACTION",
		Category: types.TagMechanic,
		Patterns: []types.Pattern{
			{Value: `from .*graveyard`, Type: types.RegexMatch},
			{Value: `return.*to.*battlefield`, Type: types.RegexMatch},
		},
		Weight:      2,
		Description: "Card interacts with graveyard",
	},
	{
		Name:     "CARD_ADVANTAGE",
		Category: types.TagStrategy,
		Patterns: []types.Pattern{
			{Value: `draw (\d+) cards?`, Type: types.RegexMatch},
			{Value: `scry (\d+)`, Type: types.RegexMatch},
		},
		Weight:      3,
		Description: "Card provides card advantage",
	},
	{
		Name:     "REMOVAL",
		Category: types.TagStrategy,
		Patterns: []types.Pattern{
			{Value: `destroy target`, Type: types.RegexMatch},
			{Value: `exile target`, Type: types.RegexMatch},
			{Value: `deal (\d+) damage to target`, Type: types.RegexMatch},
		},
		Weight:      3,
		Description: "Card removes threats",
	},
}

// ArchetypeRules defines rules for identifying deck archetypes
var ArchetypeRules = []types.TagRule{
	{
		Name:     "AGGRO",
		Category: types.TagArchetype,
		Patterns: []types.Pattern{
			{Value: "haste", Type: types.ExactMatch},
			{Value: "attack", Type: types.ExactMatch},
			{Value: "combat", Type: types.ExactMatch},
		},
		Weight:      1,
		Description: "Card supports aggressive strategies",
	},
	{
		Name:     "CONTROL",
		Category: types.TagArchetype,
		Patterns: []types.Pattern{
			{Value: "counter", Type: types.ExactMatch},
			{Value: "destroy", Type: types.ExactMatch},
			{Value: "exile", Type: types.ExactMatch},
		},
		Weight:      1,
		Description: "Card supports control strategies",
	},
}

// TimingRules defines rules for identifying timing-based effects
var TimingRules = []types.TagRule{
	{
		Name:     "INSTANT_SPEED",
		Category: types.TagTiming,
		Patterns: []types.Pattern{
			{Value: "flash", Type: types.ExactMatch},
			{Value: "at any time", Type: types.ExactMatch},
		},
		Weight:      1,
		Description: "Card can be played at instant speed",
	},
	{
		Name:     "TRIGGERED",
		Category: types.TagTiming,
		Patterns: []types.Pattern{
			{Value: "when", Type: types.ExactMatch},
			{Value: "whenever", Type: types.ExactMatch},
			{Value: "at the beginning of", Type: types.ExactMatch},
		},
		Weight:      1,
		Description: "Card has triggered abilities",
	},
}
