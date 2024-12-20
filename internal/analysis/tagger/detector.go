package tagger

import (
    "regexp"
    "strings"
    "./types"
)

// EffectDetector handles complex effect analysis
type EffectDetector struct {
    keywordPatterns map[string]*regexp.Regexp
    effectTypes     map[string][]string
}

// NewEffectDetector creates a new effect detector
func NewEffectDetector() *EffectDetector {
    detector := &EffectDetector{
        keywordPatterns: make(map[string]*regexp.Regexp),
        effectTypes: map[string][]string{
            "DESTRUCTION": {
                "destroy",
                "sacrifice",
                "exile",
            },
            "BUFF": {
                "gets? \\+\\d+/\\+\\d+",
                "gains? \\+\\d+/\\+\\d+",
            },
            "CARD_DRAW": {
                "draw \\d+ cards?",
                "draw a card",
            },
            "MANA_GENERATION": {
                "add \\w+ mana",
                "adds? one mana",
            },
        },
    }
    detector.compilePatterns()
    return detector
}

// compilePatterns precompiles regex patterns
func (ed *EffectDetector) compilePatterns() {
    for category, patterns := range ed.effectTypes {
        for _, pattern := range patterns {
            ed.keywordPatterns[category+"_"+pattern] = regexp.MustCompile(
                "(?i)" + pattern,
            )
        }
    }
}

// AnalyzeEffect performs detailed analysis of card effects
func (ed *EffectDetector) AnalyzeEffect(effect string) []types.Tag {
    var tags []types.Tag
    
    // Check each effect category
    for category, patterns := range ed.effectTypes {
        for _, pattern := range patterns {
            regexKey := category + "_" + pattern
            if ed.keywordPatterns[regexKey].MatchString(effect) {
                tags = append(tags, types.Tag{
                    Name:     category,
                    Category: types.TagMechanic,
                    Weight:   2,
                })
                break // Found match for this category
            }
        }
    }
    
    // Detect targeting
    if ed.hasTargeting(effect) {
        tags = append(tags, types.Tag{
            Name:     "TARGETED",
            Category: types.TagMechanic,
            Weight:   1,
        })
    }
    
    // Detect conditional effects
    if ed.isConditional(effect) {
        tags = append(tags, types.Tag{
            Name:     "CONDITIONAL",
            Category: types.TagMechanic,
            Weight:   1,
        })
    }
    
    // Add complexity tag
    complexity := ed.AnalyzeComplexity(effect)
    if complexity > 3 {
        tags = append(tags, types.Tag{
            Name:     "HIGH_COMPLEXITY",
            Category: types.TagMechanic,
            Weight:   1,
        })
    }
    
    return tags
}

// hasTargeting checks if the effect targets something
func (ed *EffectDetector) hasTargeting(effect string) bool {
    targetPatterns := []string{
        "target \\w+",
        "choose \\w+",
        "selected",
    }
    
    for _, pattern := range targetPatterns {
        if matched, _ := regexp.MatchString("(?i)"+pattern, effect); matched {
            return true
        }
    }
    return false
}

// isConditional checks if the effect is conditional
func (ed *EffectDetector) isConditional(effect string) bool {
    conditionalPatterns := []string{
        "if \\w+",
        "when \\w+",
        "whenever \\w+",
        "unless",
    }
    
    for _, pattern := range conditionalPatterns {
        if matched, _ := regexp.MatchString("(?i)"+pattern, effect); matched {
            return true
        }
    }
    return false
}

// DetectKeywords identifies keywords in effect text
func (ed *EffectDetector) DetectKeywords(effect string) []string {
    var keywords []string
    keywordPatterns := map[string]string{
        "FLYING":        "flying",
        "FIRST_STRIKE":  "first strike",
        "TRAMPLE":       "trample",
        "HASTE":         "haste",
        "VIGILANCE":     "vigilance",
        "DEATHTOUCH":    "deathtouch",
        "DOUBLE_STRIKE": "double strike",
    }
    
    for keyword, pattern := range keywordPatterns {
        if strings.Contains(strings.ToLower(effect), pattern) {
            keywords = append(keywords, keyword)
        }
    }
    
    return keywords
}

// AnalyzeComplexity determines the complexity of an effect
func (ed *EffectDetector) AnalyzeComplexity(effect string) int {
    complexity := 0
    
    // Count conditional statements
    complexity += strings.Count(strings.ToLower(effect), "if")
    complexity += strings.Count(strings.ToLower(effect), "when")
    complexity += strings.Count(strings.ToLower(effect), "whenever")
    
    // Count targeting
    complexity += strings.Count(strings.ToLower(effect), "target")
    
    // Count choices
    complexity += strings.Count(strings.ToLower(effect), "choose")
    complexity += strings.Count(strings.ToLower(effect), "may")
    
    // Count numbers (usually indicating quantity)
    numberPattern := regexp.MustCompile("\\d+")
    complexity += len(numberPattern.FindAllString(effect, -1))
    
    // Count keywords
    keywords := ed.DetectKeywords(effect)
    complexity += len(keywords)
    
    // Count sentence separators
    complexity += strings.Count(effect, ".")
    complexity += strings.Count(effect, ";")
    
    return complexity
}

// AnalyzeInteractions detects card interactions
func (ed *EffectDetector) AnalyzeInteractions(effect string) []types.Tag {
    var tags []types.Tag
    
    // Check for graveyard interaction
    if strings.Contains(strings.ToLower(effect), "graveyard") {
        tags = append(tags, types.Tag{
            Name:     "GRAVEYARD_INTERACTION",
            Category: types.TagMechanic,
            Weight:   2,
        })
    }
    
    // Check for library interaction
    if strings.Contains(strings.ToLower(effect), "library") {
        tags = append(tags, types.Tag{
            Name:     "LIBRARY_INTERACTION",
            Category: types.TagMechanic,
            Weight:   2,
        })
    }
    
    // Check for token creation
    if matched, _ := regexp.MatchString("create.*token", strings.ToLower(effect)); matched {
        tags = append(tags, types.Tag{
            Name:     "TOKEN_CREATOR",
            Category: types.TagMechanic,
            Weight:   2,
        })
    }
    
    // Check for counter manipulation
    if strings.Contains(strings.ToLower(effect), "counter") {
        tags = append(tags, types.Tag{
            Name:     "COUNTER_MANIPULATION",
            Category: types.TagMechanic,
            Weight:   2,
        })
    }
    
    return tags
}

// GetEffectType determines the primary type of an effect
func (ed *EffectDetector) GetEffectType(effect string) string {
    effect = strings.ToLower(effect)
    
    // Check for different effect types in order of priority
    if strings.Contains(effect, "win the game") || strings.Contains(effect, "lose the game") {
        return "GAME_ENDING"
    }
    
    if strings.Contains(effect, "damage") || strings.Contains(effect, "destroy") || 
       strings.Contains(effect, "exile") {
        return "REMOVAL"
    }
    
    if strings.Contains(effect, "draw") {
        return "CARD_ADVANTAGE"
    }
    
    if strings.Contains(effect, "search") {
        return "TUTOR"
    }
    
    if matched, _ := regexp.MatchString("\\+\\d+/\\+\\d+", effect); matched {
        return "BUFF"
    }
    
    if matched, _ := regexp.MatchString("-\\d+/-\\d+", effect); matched {
        return "DEBUFF"
    }
    
    return "UTILITY"
}