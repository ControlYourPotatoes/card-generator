package common

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
