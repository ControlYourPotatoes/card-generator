package card

import (
    "fmt"
    "strings"
)

// CardType represents the type of card
type CardType string

const (
    TypeCreature    CardType = "Creature"
    TypeArtifact    CardType = "Artifact"
    TypeSpell       CardType = "Spell"
    TypeIncantation CardType = "Incantation"
    TypeAnthem      CardType = "Anthem"
)

// Keywords represents common card keywords
var Keywords = map[string]bool{
    "OFFER":          true,
    "COMMAND":        true,
    "HASTE":          true,
    "CRITICAL":       true,
    "DMG GAIN":       true,
    "FACEOFF":        true,
    "INDESTRUCTIBLE": true,
    "DOUBLE STRIKE":  true,
    "BREAKTHROUGH":   true,
    "EVASIVE":        true,
    "GUIDE":          true,
}

// Card represents the basic interface that all cards must implement
type Card interface {
    GetName() string
    GetCost() int
    GetEffect() string
    GetType() CardType
    Validate() error
}

// BaseCard contains common fields for all cards
type BaseCard struct {
    Name   string
    Cost   int
    Effect string
    Type   CardType
}

func (b BaseCard) GetName() string    { return b.Name }
func (b BaseCard) GetCost() int      { return b.Cost }
func (b BaseCard) GetEffect() string  { return b.Effect }
func (b BaseCard) GetType() CardType  { return b.Type }

// Base validation for all cards
func (b BaseCard) Validate() error {
    if b.Name == "" {
        return fmt.Errorf("name cannot be empty")
    }
    if b.Cost < -1 {
        return fmt.Errorf("cost cannot be negative")
    }
    if b.Effect == "" {
        return fmt.Errorf("effect cannot be empty")
    }
    return nil
}

// Creature represents a creature card
type Creature struct {
    BaseCard
    Attack  int
    Defense int
    Trait string // e.g., "Demon", "Goblin"
}

func (c Creature) Validate() error {
    if err := c.BaseCard.Validate(); err != nil {
        return err
    }
    if c.Attack < 0 {
        return fmt.Errorf("attack cannot be negative")
    }
    if c.Defense < 0 {
        return fmt.Errorf("defense cannot be negative")
    }
    return nil
}

// Artifact represents an artifact card
type Artifact struct {
    BaseCard
    IsEquipment bool // Determined by checking if effect contains "Equip"
}

func (a Artifact) Validate() error {
    if err := a.BaseCard.Validate(); err != nil {
        return err
    }
    // Additional artifact-specific validation
    if a.IsEquipment && !strings.Contains(strings.ToLower(a.Effect), "equip") {
        return fmt.Errorf("equipment artifact must contain equip effect")
    }
    return nil
}

// Spell represents a spell card
type Spell struct {
    BaseCard
    TargetType string // "Creature", "Player", "Any"
}

func (s Spell) Validate() error {
    if err := s.BaseCard.Validate(); err != nil {
        return err
    }
    // Spells often have specific targeting requirements
    if s.TargetType != "" && s.TargetType != "Creature" && s.TargetType != "Player" && s.TargetType != "Any" {
        return fmt.Errorf("invalid target type: %s", s.TargetType)
    }
    return nil
}

// Incantation represents an incantation card (instant-speed spells)
type Incantation struct {
    BaseCard
    Timing string // "ON ANY CLASH", "ON ATTACK", etc.
}

func (i Incantation) Validate() error {
    if err := i.BaseCard.Validate(); err != nil {
        return err
    }
    // Validate timing phrases
    validTimings := map[string]bool{
        "ON ANY CLASH": true,
        "ON ATTACK":    true,
    }
    if i.Timing != "" && !validTimings[i.Timing] {
        return fmt.Errorf("invalid timing: %s", i.Timing)
    }
    return nil
}

// Anthem represents a continuous effect card
type Anthem struct {
    BaseCard
    Continuous bool // All anthems are continuous by nature
}

func (a Anthem) Validate() error {
    if err := a.BaseCard.Validate(); err != nil {
        return err
    }
    // Anthems should always be continuous
    if !a.Continuous {
        return fmt.Errorf("anthem must be continuous")
    }
    return nil
}

// Helper function to create a card based on type
func NewCard(cardType CardType) Card {
    baseCard := BaseCard{Type: cardType}
    switch cardType {
    case TypeCreature:
        return &Creature{BaseCard: baseCard}
    case TypeArtifact:
        return &Artifact{BaseCard: baseCard}
    case TypeSpell:
        return &Spell{BaseCard: baseCard}
    case TypeIncantation:
        return &Incantation{BaseCard: baseCard}
    case TypeAnthem:
        return &Anthem{BaseCard: baseCard, Continuous: true}
    default:
        return &baseCard
    }
}


