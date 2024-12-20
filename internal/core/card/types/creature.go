package types

import (
    "fmt"
    
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
    "github.com/ControlYourPotatoes/card-generator/internal/core/card/validation"
)

type Creature struct {
    card.BaseCard
    Attack  int
    Defense int
    Tribes  []card.Tribe
    Traits  []string
}

func (c *Creature) Validate() *validation.ValidationError {
    // Validate base properties
    baseValidator := validation.BaseValidator{
        Name:   c.Name,
        Cost:   c.Cost,
        Effect: c.Effect,
    }
    
    if err := baseValidator.ValidateBase(); err != nil {
        return err
    }
    
    // Use the centralized creature validation
    if err := validation.ValidateCreature(c.Attack, c.Defense, c.Tribes); err != nil {
        return err
    }
    
    return nil
}

func (c *Creature) ToData() *card.CardData {
    data := c.BaseCard.ToData()
    data.Attack = c.Attack
    data.Defense = c.Defense
    data.Tribes = c.Tribes
    data.Traits = c.Traits
    return data
}

// Helper methods
func NewCreature(name string, cost int) *Creature {
    return &Creature{
        BaseCard: card.BaseCard{
            Name: name,
            Cost: cost,
            Type: card.TypeCreature,
        },
        Tribes: make([]card.Tribe, 0),
        Traits: make([]string, 0),
    }
}

func (c *Creature) AddTribe(tribe card.Tribe) error {
    // Since validation will catch duplicates and invalid tribes,
    // we can simply append here and let validation handle the checks
    c.Tribes = append(c.Tribes, tribe)
    
    // Run validation to ensure the addition is valid
    if err := validation.ValidateCreature(c.Attack, c.Defense, c.Tribes); err != nil {
        // Remove the tribe if validation failed
        c.Tribes = c.Tribes[:len(c.Tribes)-1]
        return fmt.Errorf("invalid tribe addition: %v", err)
    }
    return nil
}

func (c *Creature) AddTrait(trait string) {
    c.Traits = append(c.Traits, trait)
}