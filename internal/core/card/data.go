package card

import (
    "time"

    "github.com/ControlYourPotatoes/card-generator/internal/core/common"
)



// CardData represents the serializable form of a card
type CardData struct {
    Type        CardType          `json:"type"`
    Name        string            `json:"name"`
    Cost        int               `json:"cost"`
    Effect      string            `json:"effect"`
    Attack      int               `json:"attack,omitempty"`
    Defense     int               `json:"defense,omitempty"`
    // Updated creature-specific fields
    Tribes      []common.Tribe    `json:"tribes,omitempty"`      // e.g., ["Human", "Goblin", "Zombie"]
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

