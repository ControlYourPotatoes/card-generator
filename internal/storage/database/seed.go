package database

import (
	"fmt"
	"log"

	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
)

// Seeder is responsible for seeding the database with initial data
type Seeder struct {
	store *PostgresStore
}

// NewSeeder creates a new database seeder
func NewSeeder(store *PostgresStore) *Seeder {
	return &Seeder{
		store: store,
	}
}

// SeedCardTypes seeds the database with card types
func (s *Seeder) SeedCardTypes() error {
	cardTypes := []string{
		string(card.TypeCreature),
		string(card.TypeArtifact),
		string(card.TypeSpell),
		string(card.TypeIncantation),
		string(card.TypeAnthem),
	}
	
	for _, t := range cardTypes {
		_, err := s.store.db.Exec(
			`INSERT INTO card_types (name) VALUES ($1) 
			 ON CONFLICT (name) DO NOTHING`,
			t,
		)
		if err != nil {
			return fmt.Errorf("failed to seed card type %s: %w", t, err)
		}
	}
	
	log.Println("Seeded card types")
	return nil
}

// SeedTraits seeds the database with creature traits
func (s *Seeder) SeedTraits() error {
	traits := []string{
		string(card.TraitBeast),
		string(card.TraitWarrior),
		string(card.TraitDragon),
		string(card.TraitDemon),
		string(card.TraitAngel),
		string(card.TraitLegendary),
		string(card.TraitAncient),
		string(card.TraitDivine),
	}
	
	for _, t := range traits {
		_, err := s.store.db.Exec(
			`INSERT INTO traits (name) VALUES ($1) 
			 ON CONFLICT (name) DO NOTHING`,
			t,
		)
		if err != nil {
			return fmt.Errorf("failed to seed trait %s: %w", t, err)
		}
	}
	
	log.Println("Seeded traits")
	return nil
}

// SeedKeywords seeds the database with keywords
func (s *Seeder) SeedKeywords() error {
	keywords := []string{
		"HASTE",
		"CRITICAL",
		"EQUIPMENT",
		"DAMAGE",
		"BUFF",
		"COUNTER",
		"DRAW",
		"DIRECT",
		"FLYING",
		"IMMUNE",
	}
	
	for _, k := range keywords {
		_, err := s.store.db.Exec(
			`INSERT INTO keywords (name) VALUES ($1) 
			 ON CONFLICT (name) DO NOTHING`,
			k,
		)
		if err != nil {
			return fmt.Errorf("failed to seed keyword %s: %w", k, err)
		}
	}
	
	log.Println("Seeded keywords")
	return nil
}

// SeedSampleCards seeds the database with sample cards
func (s *Seeder) SeedSampleCards() error {
	// Create sample cards
	cards := []card.Card{
		// Creature 1
		&card.Creature{
			BaseCard: card.BaseCard{
				Name:     "Sample Dragon",
				Cost:     5,
				Effect:   "When this creature enters play, deal 2 damage to target creature.",
				Type:     card.TypeCreature,
				Keywords: []string{"DAMAGE"},
			},
			Attack:  4,
			Defense: 4,
			Trait:   card.TraitDragon,
		},
		// Creature 2
		&card.Creature{
			BaseCard: card.BaseCard{
				Name:     "Elite Warrior",
				Cost:     3,
				Effect:   "Critical: Deal double damage when attacking.",
				Type:     card.TypeCreature,
				Keywords: []string{"CRITICAL"},
			},
			Attack:  2,
			Defense: 3,
			Trait:   card.TraitWarrior,
		},
		// Spell 1
		&card.Spell{
			BaseCard: card.BaseCard{
				Name:     "Lightning Bolt",
				Cost:     1,
				Effect:   "Deal 3 damage to target creature.",
				Type:     card.TypeSpell,
				Keywords: []string{"DAMAGE"},
			},
			TargetType: "Creature",
		},
		// Spell 2
		&card.Spell{
			BaseCard: card.BaseCard{
				Name:     "Mind Control",
				Cost:     5,
				Effect:   "Gain control of target creature until end of turn.",
				Type:     card.TypeSpell,
				Keywords: []string{"CONTROL"},
			},
			TargetType: "Creature",
		},
		// Artifact
		&card.Artifact{
			BaseCard: card.BaseCard{
				Name:     "Dragon Sword",
				Cost:     3,
				Effect:   "Equip to a creature. Equipped creature gets +2/+0.",
				Type:     card.TypeArtifact,
				Keywords: []string{"EQUIPMENT"},
			},
			IsEquipment: true,
		},
		// Incantation
		&card.Incantation{
			BaseCard: card.BaseCard{
				Name:     "Combat Trick",
				Cost:     2,
				Effect:   "ON ATTACK: Target creature gets +1/+1 until end of turn.",
				Type:     card.TypeIncantation,
				Keywords: []string{"BUFF"},
			},
			Timing: "ON ATTACK",
		},
		// Anthem
		&card.Anthem{
			BaseCard: card.BaseCard{
				Name:     "Glory of Battle",
				Cost:     4,
				Effect:   "All creatures you control get +1/+1.",
				Type:     card.TypeAnthem,
				Keywords: []string{"BUFF"},
			},
			Continuous: true,
		},
	}
	
	// Save each card
	for _, c := range cards {
		_, err := s.store.Save(c)
		if err != nil {
			return fmt.Errorf("failed to seed card %s: %w", c.GetName(), err)
		}
	}
	
	log.Println("Seeded sample cards")
	return nil
}

// SeedAll seeds the database with all data
func (s *Seeder) SeedAll() error {
	// First, make sure schema is initialized
	if err := s.store.InitSchema(); err != nil {
		return err
	}
	
	// Seed card types
	if err := s.SeedCardTypes(); err != nil {
		return err
	}
	
	// Seed traits
	if err := s.SeedTraits(); err != nil {
		return err
	}
	
	// Seed keywords
	if err := s.SeedKeywords(); err != nil {
		return err
	}
	
	// Seed sample cards
	if err := s.SeedSampleCards(); err != nil {
		return err
	}
	
	log.Println("Database seeded successfully")
	return nil
}