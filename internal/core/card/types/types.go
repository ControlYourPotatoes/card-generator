package types

// CardType represents the type of a card
type CardType string

const (
    TypeCreature    CardType = "Creature"
    TypeArtifact    CardType = "Artifact"
    TypeSpell       CardType = "Spell"
    TypeIncantation CardType = "Incantation"
    TypeAnthem      CardType = "Anthem"
)

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
