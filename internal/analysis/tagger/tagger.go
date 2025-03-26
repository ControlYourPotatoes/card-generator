package tagger

import (
    "strings"
    "sync"
    "fmt"
    
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
    "github.com/ControlYourPotatoes/card-generator/internal/analysis/types"
    "github.com/ControlYourPotatoes/card-generator/internal/analysis/rules"
    "github.com/dlclark/regexp2"
)

// CardTagger implements the main tagging logic
type CardTagger struct {
    rules        []types.TagRule
    regexCache   map[string]*regexp2.Regexp
    cacheMutex   sync.RWMutex
    manualTags   map[string][]types.Tag
    tagValidator *types.TagValidator
}

// NewCardTagger creates a new instance of CardTagger
func NewCardTagger() *CardTagger {
    tagger := &CardTagger{
        regexCache: make(map[string]*regexp2.Regexp),
        manualTags: make(map[string][]types.Tag),
        tagValidator: types.NewTagValidator(),
    }
    tagger.initializeRules()
    return tagger
}

// initializeRules loads all tagging rules
func (ct *CardTagger) initializeRules() {
    ct.rules = append(ct.rules, rules.BaseRules...)
    ct.rules = append(ct.rules, rules.TribalRules...)
    
    // Add archetype rules
    for _, rule := range rules.ArchetypeRules {
        ct.rules = append(ct.rules, rule)
    }
    
    // Add timing rules
    for _, rule := range rules.TimingRules {
        ct.rules = append(ct.rules, rule)
    }
}

// GenerateTags generates tags for a card
func (ct *CardTagger) GenerateTags(c interface{}) ([]types.Tag, error) {
    var cardData *card.CardData
    
    switch v := c.(type) {
    case *card.CardData:
        cardData = v
    case card.Card:
        cardData = v.ToData()
    default:
        return nil, fmt.Errorf("invalid card type: must be *card.CardData or card.Card")
    }

    var tags []types.Tag
    
    // Generate basic tags
    basicTags := ct.generateBasicTags(cardData)
    tags = append(tags, basicTags...)
    
    // Generate tribal tags using the card's tribes
    tribalTags := rules.GenerateTribalTags(
        string(cardData.Type),
        cardData.Effect,
        cardData.Tribes,
    )
    tags = append(tags, tribalTags...)
    
    // Generate class-based tags if the card has classes
    if len(cardData.Classes) > 0 {
        classTags := ct.generateClassTags(cardData.Classes)
        tags = append(tags, classTags...)
    }
    
    // Generate combo tags
    comboTags := rules.DetectComboPotential(cardData.Effect)
    tags = append(tags, comboTags...)
    
    // Apply pattern-based rules
    ruleTags := ct.applyRules(cardData)
    tags = append(tags, ruleTags...)
    
    // Add any manual tags
    if manualTags, exists := ct.manualTags[cardData.Name]; exists {
        tags = append(tags, manualTags...)
    }
    
    // Deduplicate and validate tags
    tags = ct.deduplicateTags(tags)
    if err := ct.tagValidator.ValidateTags(tags); err != nil {
        return nil, err
    }
    
    return tags, nil
}

// generateClassTags creates tags based on class types
func (ct *CardTagger) generateClassTags(classes []string) []types.Tag {
    var tags []types.Tag
    
    for _, class := range classes {
        tags = append(tags, types.Tag{
            Name:     class + "_CLASS",
            Category: types.TagMechanic,
            Weight:   1,
        })
    }
    
    return tags
}


// generateBasicTags creates basic tags based on card properties
func (ct *CardTagger) generateBasicTags(card *card.CardData) []types.Tag {
    var tags []types.Tag
    
    // Add type-based tag
    tags = append(tags, types.Tag{
        Name:     string(card.Type),
        Category: types.TagMechanic,
        Weight:   1,
    })
    
    // Add cost-based tags
    if card.Cost > 0 {
        costCategory := ct.categorizeCost(card.Cost)
        tags = append(tags, types.Tag{
            Name:     costCategory,
            Category: types.TagCost,
            Weight:   1,
        })
    }
    
    return tags
}

// applyRules applies all pattern-based rules to generate tags
func (ct *CardTagger) applyRules(card *card.CardData) []types.Tag {
    var tags []types.Tag
    
    for _, rule := range ct.rules {
        if ct.ruleMatches(card, rule) {
            tags = append(tags, types.Tag{
                Name:     rule.Name,
                Category: rule.Category,
                Weight:   rule.Weight,
            })
        }
    }
    
    return tags
}

// ruleMatches checks if a card matches a specific rule
func (ct *CardTagger) ruleMatches(card *card.CardData, rule types.TagRule) bool {
    for _, pattern := range rule.Patterns {
        matched := false
        
        switch pattern.Type {
        case types.ExactMatch:
            matched = strings.Contains(
                strings.ToLower(card.Effect),
                strings.ToLower(pattern.Value),
            )
            
        case types.RegexMatch:
            matched = ct.matchRegex(pattern.Value, card.Effect)
            
        case types.ProximityMatch:
            matched = ct.checkProximity(
                pattern.Value,
                card.Effect,
                pattern.Proximity,
            )
        }
        
        if !matched {
            return false
        }
    }
    
    // Check conditions if patterns match
    return ct.checkConditions(card, rule.Conditions)
}

// matchRegex performs regex matching with caching
func (ct *CardTagger) matchRegex(pattern, text string) bool {
    ct.cacheMutex.RLock()
    regex, exists := ct.regexCache[pattern]
    ct.cacheMutex.RUnlock()
    
    if !exists {
        var err error
        regex, err = regexp2.Compile(pattern, regexp2.IgnoreCase)
        if err != nil {
            return false
        }
        
        ct.cacheMutex.Lock()
        ct.regexCache[pattern] = regex
        ct.cacheMutex.Unlock()
    }
    
    match, _ := regex.MatchString(text)
    return match
}

// checkProximity checks if words appear within a certain distance
func (ct *CardTagger) checkProximity(pattern, text string, maxDistance int) bool {
    words := strings.Fields(strings.ToLower(text))
    patternWords := strings.Fields(strings.ToLower(pattern))
    
    for i := 0; i <= len(words)-len(patternWords); i++ {
        matched := true
        for j, patternWord := range patternWords {
            if i+j >= len(words) || words[i+j] != patternWord {
                matched = false
                break
            }
        }
        if matched && i <= maxDistance {
            return true
        }
    }
    return false
}

// checkConditions verifies if a card meets all conditions
func (ct *CardTagger) checkConditions(card *card.CardData, conditions []types.Condition) bool {
    for _, condition := range conditions {
        switch condition.Type {
        case types.IsType:
            if card.Type != condition.Value {
                return false
            }
            
        case types.HasKeyword:
            if !ct.hasKeyword(card, condition.Value.(string)) {
                return false
            }
            
        case types.CostCondition:
            if !ct.checkCostCondition(card.Cost, condition.Value) {
                return false
            }
        }
    }
    return true
}

// AddManualTag adds a manual tag to a card
func (ct *CardTagger) AddManualTag(cardID string, tag types.Tag) error {
    if err := ct.tagValidator.ValidateTag(tag); err != nil {
        return err
    }
    
    ct.cacheMutex.Lock()
    defer ct.cacheMutex.Unlock()
    
    ct.manualTags[cardID] = append(ct.manualTags[cardID], tag)
    return nil
}

// Helper functions
func (ct *CardTagger) deduplicateTags(tags []types.Tag) []types.Tag {
    seen := make(map[string]bool)
    result := make([]types.Tag, 0)
    
    for _, tag := range tags {
        key := tag.Name + string(tag.Category)
        if !seen[key] {
            seen[key] = true
            result = append(result, tag)
        }
    }
    return result
}

func (ct *CardTagger) categorizeCost(cost int) string {
    switch {
    case cost <= 2:
        return "LOW_COST"
    case cost <= 4:
        return "MID_COST"
    case cost <= 6:
        return "HIGH_COST"
    default:
        return "VERY_HIGH_COST"
    }
}

func (ct *CardTagger) hasKeyword(card *card.CardData, keyword string) bool {
    return strings.Contains(
        strings.ToLower(card.Effect),
        strings.ToLower(keyword),
    )
}

func (ct *CardTagger) checkCostCondition(cardCost int, condition interface{}) bool {
    switch v := condition.(type) {
    case int:
        return cardCost == v
    case map[string]int:
        if min, ok := v["min"]; ok && cardCost < min {
            return false
        }
        if max, ok := v["max"]; ok && cardCost > max {
            return false
        }
        return true
    default:
        return false
    }
}