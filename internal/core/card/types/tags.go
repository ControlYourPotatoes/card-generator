package types

// Tag represents a single card tag with its category
type Tag struct {
    Name     string      `json:"name"`
    Category TagCategory `json:"category"`
    Weight   int         `json:"weight,omitempty"`
}

// TagCategory represents different types of tags
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

// TagRule defines a rule for generating tags
type TagRule struct {
    Name        string
    Category    TagCategory
    Patterns    []Pattern
    Conditions  []Condition
    Weight      int
    Description string
}

// Pattern defines how to match card text
type Pattern struct {
    Value     string
    Type      PatternType
    Proximity int
}

// PatternType defines different ways to match patterns
type PatternType int

const (
    ExactMatch PatternType = iota
    RegexMatch
    ProximityMatch
    NegationMatch
)

// Condition defines additional requirements for a tag
type Condition struct {
    Type  ConditionType
    Value interface{}
}

// ConditionType defines different types of conditions
type ConditionType string

const (
    IsType         ConditionType = "IS_TYPE"
    HasKeyword     ConditionType = "HAS_KEYWORD"
    CostCondition  ConditionType = "COST"
    PowerCondition ConditionType = "POWER"
    ComboCondition ConditionType = "COMBO"
)

// Tagger defines the interface for tag generation
type Tagger interface {
    GenerateTags(card interface{}) ([]Tag, error)
    AddManualTag(cardID string, tag Tag) error
    ValidateTags(tags []Tag) error
}