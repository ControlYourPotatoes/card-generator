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