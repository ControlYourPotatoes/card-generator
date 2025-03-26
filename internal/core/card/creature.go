package card

// Trait represents a creature trait
type Trait string

// Predefined creature traits
const (
	TraitBeast     Trait = "Beast"
	TraitWarrior   Trait = "Warrior"
	TraitDragon    Trait = "Dragon"
	TraitDemon     Trait = "Demon"
	TraitAngel     Trait = "Angel"
	TraitLegendary Trait = "Legendary"
	TraitAncient   Trait = "Ancient"
	TraitDivine    Trait = "Divine"
)

// IsValid checks if a trait is valid
func (t Trait) IsValid() bool {
	validTraits := map[Trait]bool{
		TraitBeast:     true,
		TraitWarrior:   true,
		TraitDragon:    true,
		TraitDemon:     true,
		TraitAngel:     true,
		TraitLegendary: true,
		TraitAncient:   true,
		TraitDivine:    true,
	}
	return validTraits[t]
}

// Creature represents a creature card
type Creature struct {
	BaseCard
	Attack  int
	Defense int
	Trait   Trait
}

// Validate performs validation specific to Creature cards
func (c *Creature) Validate() error {
	// First validate the base card
	if err := c.BaseCard.Validate(); err != nil {
		return err
	}
	
	// Validate creature-specific fields
	if c.Attack < 0 {
		return NewValidationError("attack cannot be negative", "attack")
	}
	
	if c.Defense < 0 {
		return NewValidationError("defense cannot be negative", "defense")
	}
	
	if c.Trait != "" && !c.Trait.IsValid() {
		return NewValidationError("invalid trait: "+string(c.Trait), "trait")
	}
	
	return nil
}

// ToDTO converts Creature to CardDTO
func (c *Creature) ToDTO() *CardDTO {
	dto := c.BaseCard.ToDTO()
	dto.Attack = c.Attack
	dto.Defense = c.Defense
	dto.Trait = string(c.Trait)
	return dto
}

// ToData is for backward compatibility
func (c *Creature) ToData() *CardDTO {
	return c.ToDTO()
}

// NewCreatureFromDTO creates a Creature from a CardDTO
func NewCreatureFromDTO(dto *CardDTO) *Creature {
	return &Creature{
		BaseCard: NewBaseCardFromDTO(dto),
		Attack:   dto.Attack,
		Defense:  dto.Defense,
		Trait:    Trait(dto.Trait),
	}
}