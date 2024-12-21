package card

import (
	"time"
     
	"github.com/ControlYourPotatoes/card-generator/internal/core/common"
)

// CardType represents the type of a card
type CardType string

const (
    TypeCreature    CardType = "Creature"
    TypeArtifact    CardType = "Artifact"
    TypeSpell       CardType = "Spell"
    TypeIncantation CardType = "Incantation"
    TypeAnthem      CardType = "Anthem"
)

// Card defines the core interface that all cards must implement
type Card interface {
    GetName() string
    GetCost() int
    GetEffect() string
    GetType() CardType
    Validate() *common.ValidationError
    ToData() *CardData
}

// BaseCard provides common functionality for all card types
type BaseCard struct {
    Name      string
    Cost      int
    Effect    string
    Type      CardType
    Keywords  []string
    CreatedAt time.Time
    UpdatedAt time.Time
    Metadata  map[string]string
}

func (b BaseCard) GetName() string   { return b.Name }
func (b BaseCard) GetCost() int     { return b.Cost }
func (b BaseCard) GetEffect() string { return b.Effect }
func (b BaseCard) GetType() CardType { return b.Type }