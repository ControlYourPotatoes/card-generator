package card

import (
    "time"
    
)

// Tribe represents a creature tribe
type Tribe string

const (
    TribeZombie  Tribe = "Zombie"
    TribeHuman   Tribe = "Human"
    TribeDemon   Tribe = "Demon"
    TribeGoblin  Tribe = "Goblin"
    TribeVampire Tribe = "Vampire"
    TribeGod     Tribe = "God"
)

// ValidTribes contains all currently valid tribes
var ValidTribes = map[Tribe]bool{
    TribeZombie:  true,
    TribeHuman:   true,
    TribeDemon:   true,
    TribeGoblin:  true,
    TribeVampire: true,
    TribeGod:     true,
}

// CardData represents the serializable form of a card
type CardData struct {
    Type        CardType           `json:"type"`
    Name        string            `json:"name"`
    Cost        int               `json:"cost"`
    Effect      string            `json:"effect"`
    Attack      int               `json:"attack,omitempty"`
    Defense     int               `json:"defense,omitempty"`
    // Updated creature-specific fields
    Tribes      []Tribe         `json:"tribes,omitempty"`      // e.g., ["Human", "Goblin", "Zombie"]
    Classes     []string          `json:"classes,omitempty"`     // e.g., ["Warrior", "Wizard", "Cleric"]
    Traits      []string          `json:"traits,omitempty"`

    IsEquipment bool              `json:"is_equipment,omitempty"`
    TargetType  string            `json:"target_type,omitempty"`
    Timing      string            `json:"timing,omitempty"`
    Continuous  bool              `json:"continuous,omitempty"`
    Keywords    []string          `json:"keywords,omitempty"`
    Tags       []string          `json:"tags,omitempty"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}

// Keywords represents valid card keywords
var Keywords = map[string]bool{
    "CRITICAL":       true,
    "DOUBLE STRIKE":  true,
    "INDESTRUCTIBLE": true,
    "BREAKTHROUGH":   true,
    "EVASIVE":        true,
    "HASTE":          true,
    "GUIDE":          true,
    "OFFER":          true,
    "COMMAND":        true,
    "DMG GAIN":       true,
    "FACEOFF":        true,
}

// ToData converts BaseCard to CardData
func (b BaseCard) ToData() *CardData {
    return &CardData{
        Type:      b.Type,
        Name:      b.Name,
        Cost:      b.Cost,
        Effect:    b.Effect,
        Keywords:  b.Keywords,
        CreatedAt: b.CreatedAt,
        UpdatedAt: b.UpdatedAt,
        Metadata:  b.Metadata,
    }
}

