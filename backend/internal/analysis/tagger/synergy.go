package tagger

import (
	"strings"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/analysis/types"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// SynergyDetector handles detection of card synergies
type SynergyDetector struct {
	tribalTypes     []string
	mechanicGroups  map[string][]string
	synergyPatterns map[string][]string
	comboPatterns   map[string][]string
}

// NewSynergyDetector creates a new synergy detector
func NewSynergyDetector() *SynergyDetector {
	return &SynergyDetector{
		tribalTypes: []string{
			"Zombie",
			"Dragon",
			"Warrior",
			"Beast",
			"Angel",
			"Demon",
		},
		mechanicGroups: map[string][]string{
			"TOKENS": {
				"create token",
				"token enters",
				"tokens you control",
			},
			"GRAVEYARD": {
				"from graveyard",
				"when dies",
				"return from graveyard",
			},
			"COUNTERS": {
				"put counter",
				"remove counter",
				"counter on",
			},
		},
		synergyPatterns: map[string][]string{
			"SACRIFICE": {
				"sacrifice",
				"when dies",
				"creature dies",
			},
			"EQUIPMENT": {
				"equip",
				"equipped creature",
				"attach",
			},
			"SPELLS": {
				"cast spell",
				"whenever you cast",
				"spell is cast",
			},
		},
		comboPatterns: map[string][]string{
			"INFINITE_MANA": {
				"untap",
				"add mana",
			},
			"INFINITE_TOKENS": {
				"create token",
				"copy",
				"double",
			},
			"INFINITE_DAMAGE": {
				"deal damage",
				"when deals damage",
				"double damage",
			},
		},
	}
}

// AnalyzeSynergies detects potential synergies
func (sd *SynergyDetector) AnalyzeSynergies(card *card.CardDTO) []types.Tag {
	var tags []types.Tag

	// Check tribal synergies
	tribalTags := sd.detectTribalSynergies(card)
	tags = append(tags, tribalTags...)

	// Check mechanic synergies
	mechanicTags := sd.detectMechanicSynergies(card)
	tags = append(tags, mechanicTags...)

	// Check specific synergy patterns
	patternTags := sd.detectSynergyPatterns(card)
	tags = append(tags, patternTags...)

	return tags
}

// detectTribalSynergies checks for tribal-based synergies
func (sd *SynergyDetector) detectTribalSynergies(card *card.CardData) []types.Tag {
	var tags []types.Tag
	effectLower := strings.ToLower(card.Effect)

	for _, tribe := range sd.tribalTypes {
		tribeLower := strings.ToLower(tribe)

		// Check if card mentions the tribe
		if strings.Contains(effectLower, tribeLower) {
			tags = append(tags, types.Tag{
				Name:     tribe + "_SYNERGY",
				Category: types.TagSynergy,
				Weight:   2,
			})

			// Check for tribal lord effects
			if strings.Contains(effectLower, "other "+tribeLower) &&
				strings.Contains(effectLower, "get") {
				tags = append(tags, types.Tag{
					Name:     tribe + "_LORD",
					Category: types.TagSynergy,
					Weight:   3,
				})
			}
		}
	}

	return tags
}

// detectMechanicSynergies checks for mechanic-based synergies
func (sd *SynergyDetector) detectMechanicSynergies(card *card.CardData) []types.Tag {
	var tags []types.Tag
	effectLower := strings.ToLower(card.Effect)

	for mechanic, patterns := range sd.mechanicGroups {
		for _, pattern := range patterns {
			if strings.Contains(effectLower, pattern) {
				tags = append(tags, types.Tag{
					Name:     mechanic + "_SYNERGY",
					Category: types.TagSynergy,
					Weight:   2,
				})
				break
			}
		}
	}

	return tags
}

// detectSynergyPatterns checks for specific synergy patterns
func (sd *SynergyDetector) detectSynergyPatterns(card *card.CardData) []types.Tag {
	var tags []types.Tag
	effectLower := strings.ToLower(card.Effect)

	for pattern, triggers := range sd.synergyPatterns {
		for _, trigger := range triggers {
			if strings.Contains(effectLower, trigger) {
				tags = append(tags, types.Tag{
					Name:     pattern + "_ENABLER",
					Category: types.TagSynergy,
					Weight:   2,
				})
				break
			}
		}
	}

	return tags
}

// AnalyzeComboSynergies detects potential combo synergies between multiple cards
func (sd *SynergyDetector) AnalyzeComboSynergies(cards []card.CardData) []types.Tag {
	var tags []types.Tag

	// Check each combo pattern
	for comboName, patterns := range sd.comboPatterns {
		if sd.detectComboPattern(cards, patterns) {
			tags = append(tags, types.Tag{
				Name:     comboName,
				Category: types.TagCombo,
				Weight:   4,
			})
		}
	}

	// Check for powerful card advantage combinations
	if sd.detectCardAdvantageCombo(cards) {
		tags = append(tags, types.Tag{
			Name:     "CARD_ADVANTAGE_ENGINE",
			Category: types.TagCombo,
			Weight:   3,
		})
	}

	// Check for resource generation engines
	if sd.detectResourceEngine(cards) {
		tags = append(tags, types.Tag{
			Name:     "RESOURCE_ENGINE",
			Category: types.TagCombo,
			Weight:   3,
		})
	}

	return tags
}

// detectComboPattern checks if a set of cards matches a combo pattern
func (sd *SynergyDetector) detectComboPattern(cards []card.CardData, patterns []string) bool {
	patternMatches := make(map[string]bool)

	for _, card := range cards {
		effectLower := strings.ToLower(card.Effect)
		for _, pattern := range patterns {
			if strings.Contains(effectLower, pattern) {
				patternMatches[pattern] = true
			}
		}
	}

	// Check if all patterns were matched
	for _, pattern := range patterns {
		if !patternMatches[pattern] {
			return false
		}
	}

	return true
}

// detectCardAdvantageCombo checks for card draw/filtering combinations
func (sd *SynergyDetector) detectCardAdvantageCombo(cards []card.CardData) bool {
	drawCount := 0
	filterCount := 0

	for _, card := range cards {
		effectLower := strings.ToLower(card.Effect)

		if strings.Contains(effectLower, "draw") {
			drawCount++
		}
		if strings.Contains(effectLower, "scry") || strings.Contains(effectLower, "look at") {
			filterCount++
		}
	}

	return drawCount >= 2 || (drawCount >= 1 && filterCount >= 2)
}

// detectResourceEngine checks for resource generation combinations
func (sd *SynergyDetector) detectResourceEngine(cards []card.CardData) bool {
	resourceGen := 0
	resourceUse := 0

	for _, card := range cards {
		effectLower := strings.ToLower(card.Effect)

		if strings.Contains(effectLower, "add") && strings.Contains(effectLower, "mana") {
			resourceGen++
		}
		if strings.Contains(effectLower, "pay") || strings.Contains(effectLower, "cost") {
			resourceUse++
		}
	}

	return resourceGen >= 2 && resourceUse >= 1
}
