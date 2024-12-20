package tagger

import (
    "regexp"
    "strings"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/dlclark/regexp2"
)

// Enhanced tag categories
type TagCategory string

const (
    TagTribal     TagCategory = "TRIBAL"
    TagMechanic   TagCategory = "MECHANIC"
    TagStrategy   TagCategory = "STRATEGY"
    TagCost       TagCategory = "COST"
    TagSynergy    TagCategory = "SYNERGY"
    TagCombo      TagCategory = "COMBO"
    TagArchetype  TagCategory = "ARCHETYPE"
    TagTiming     TagCategory = "TIMING"
)

// PatternType defines how a pattern should be matched
type PatternType int

const (
    ExactMatch PatternType = iota
    RegexMatch
    ProximityMatch
    NegationMatch
)

// Enhanced TagRule with more sophisticated matching
type TagRule struct {
    Tag         Tag
    Patterns    []Pattern
    Conditions  []Condition
    Description string
    Weight      int // Higher weight means more significant for deck building
}

type Pattern struct {
    Value     string
    Type      PatternType
    Proximity int // For ProximityMatch, how close words should be
}

type Condition struct {
    Type  ConditionType
    Value interface{}
}

type ConditionType string

const (
    IsType          ConditionType = "IS_TYPE"
    HasKeyword      ConditionType = "HAS_KEYWORD"
    CostCondition   ConditionType = "COST"
    PowerCondition  ConditionType = "POWER"
    ComboCondition  ConditionType = "COMBO"
)

// Complex pattern definitions for different card effects
var complexPatterns = map[string][]Pattern{
    "CARD_ADVANTAGE": {
        {Value: `draw (\d+) cards?`, Type: RegexMatch},
        {Value: `scry (\d+)`, Type: RegexMatch},
        {Value: "investigate", Type: ExactMatch},
    },
    "REMOVAL": {
        {Value: `destroy target`, Type: RegexMatch},
        {Value: `exile target`, Type: RegexMatch},
        {Value: `deal (\d+) damage to target`, Type: RegexMatch},
    },
    
    "TRIBAL_LORD": {
        {Value: `other .* you control get`, Type: RegexMatch},
        {Value: `creatures you control of the chosen type`, Type: ExactMatch},
    },
}

// Combo detection patterns
var comboPatterns = map[string][]Pattern{
    "INFINITE_COMBO_POTENTIAL": {
        {Value: "untap", Type: ProximityMatch, Proximity: 5},
        {Value: "return.*from.*graveyard", Type: RegexMatch},
    },
    "MANA_COMBO": {
        {Value: "add.*for each", Type: RegexMatch},
        {Value: "doubles", Type: ProximityMatch, Proximity: 3},
    },
}

// ArchetypeDetector helps identify deck archetypes
type ArchetypeDetector struct {
    patterns map[string][]Pattern
    weights  map[string]int
}

func NewArchetypeDetector() *ArchetypeDetector {
    return &ArchetypeDetector{
        patterns: map[string][]Pattern{
            "AGGRO": {
                {Value: "haste", Type: ExactMatch},
                {Value: "attack", Type: ExactMatch},
                {Value: "combat", Type: ExactMatch},
            },
            "CONTROL": {
                {Value: "counter", Type: ExactMatch},
                {Value: "destroy", Type: ExactMatch},
                {Value: "exile", Type: ExactMatch},
            },
            "COMBO": {
                {Value: "whenever", Type: ExactMatch},
                {Value: "trigger", Type: ExactMatch},
                {Value: "copy", Type: ExactMatch},
            },
            "MIDRANGE": {
                {Value: "value", Type: ExactMatch},
                {Value: "draw", Type: ExactMatch},
                {Value: "return", Type: ExactMatch},
            },
        },
        weights: map[string]int{
            "AGGRO":    1,
            "CONTROL":  1,
            "COMBO":    1,
            "MIDRANGE": 1,
        },
    }
}

// Enhanced CardTagger with more capabilities
type CardTagger struct {
    rules             []TagRule
    archetypeDetector *ArchetypeDetector
    regexCache        map[string]*regexp2.Regexp
}

func NewCardTagger() *CardTagger {
    tagger := &CardTagger{
        regexCache: make(map[string]*regexp2.Regexp),
        archetypeDetector: NewArchetypeDetector(),
    }
    tagger.initializeRules()
    return tagger
}

// GenerateTags now includes more sophisticated analysis
func (ct *CardTagger) GenerateTags(c *card.CardData) []Tag {
    var tags []Tag
    
    // Basic tags
    tags = append(tags, ct.generateBasicTags(c)...)
    
    // Complex pattern matching
    tags = append(tags, ct.matchComplexPatterns(c)...)
    
    // Combo detection
    tags = append(tags, ct.detectCombos(c)...)
    
    // Archetype analysis
    tags = append(tags, ct.analyzeArchetype(c)...)
    
    // Synergy detection
    tags = append(tags, ct.detectSynergies(c)...)
    
    return ct.deduplicateAndPrioritize(tags)
}

// Example implementation of one of the analysis methods
func (ct *CardTagger) detectCombos(c *card.CardData) []Tag {
    var tags []Tag
    effectLower := strings.ToLower(c.Effect)
    
    for comboName, patterns := range comboPatterns {
        matched := true
        for _, pattern := range patterns {
            switch pattern.Type {
            case RegexMatch:
                if !ct.matchRegex(pattern.Value, effectLower) {
                    matched = false
                    break
                }
            case ProximityMatch:
                if !ct.checkProximity(pattern.Value, effectLower, pattern.Proximity) {
                    matched = false
                    break
                }
            }
        }
        if matched {
            tags = append(tags, Tag{comboName, TagCombo})
        }
    }
    
    return tags
}

// Helper method for regex matching with caching
func (ct *CardTagger) matchRegex(pattern, text string) bool {
    regex, exists := ct.regexCache[pattern]
    if !exists {
        var err error
        regex, err = regexp2.Compile(pattern, regexp2.IgnoreCase)
        if err != nil {
            return false
        }
        ct.regexCache[pattern] = regex
    }
    
    match, _ := regex.MatchString(text)
    return match
}

// Helper method for proximity checking
func (ct *CardTagger) checkProximity(pattern, text string, maxDistance int) bool {
    words := strings.Fields(text)
    patternWords := strings.Fields(pattern)
    
    for i := 0; i <= len(words)-len(patternWords); i++ {
        matched := true
        for j, patternWord := range patternWords {
            if i+j >= len(words) || words[i+j] != patternWord {
                matched = false
                break
            }
        }
        if matched {
            return true
        }
    }
    return false
}