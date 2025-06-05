package database

import (
	"context"
	"fmt"
	"log"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/storage/database/migration"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore implements the Store interface using PostgreSQL
type PostgresStore struct {
	pool *pgxpool.Pool
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(connString string) (*PostgresStore, error) {
	// Create a connection pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{pool: pool}, nil
}

// Save stores a card and returns its ID
func (s *PostgresStore) Save(c card.Card) (string, error) {
	return s.saveCard(c)
}

// Load retrieves a card by its ID
func (s *PostgresStore) Load(id string) (card.Card, error) {
	return s.loadCard(id)
}

// List returns all stored cards
func (s *PostgresStore) List() ([]card.Card, error) {
	return s.listCards()
}

// Delete removes a card by its ID
func (s *PostgresStore) Delete(id string) error {
	return s.deleteCard(id)
}

// Close cleans up any resources
func (s *PostgresStore) Close() error {
	if s.pool != nil {
		s.pool.Close()
	}
	return nil
}

// InitSchema sets up the necessary database schema
func (s *PostgresStore) InitSchema() error {
	// Read the schema definition from the migration file
	schema, err := migration.GetInitialSchema()
	if err != nil {
		return fmt.Errorf("failed to get initial schema: %w", err)
	}

	// Execute the schema
	_, err = s.pool.Exec(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database schema initialized")
	return nil
}

// SeedTestData populates the database with test data
func (s *PostgresStore) SeedTestData() error {
	// First, make sure schema is initialized
	if err := s.InitSchema(); err != nil {
		return err
	}

	// Create card types if they don't exist
	cardTypes := []string{
		string(card.TypeCreature),
		string(card.TypeArtifact),
		string(card.TypeSpell),
		string(card.TypeIncantation),
		string(card.TypeAnthem),
	}

	for _, t := range cardTypes {
		_, err := s.pool.Exec(
			context.Background(),
			`INSERT INTO card_types (name) VALUES ($1) 
			 ON CONFLICT (name) DO NOTHING`,
			t,
		)
		if err != nil {
			return fmt.Errorf("failed to seed card type %s: %w", t, err)
		}
	}

	// Create some traits
	traits := []string{
		string(card.TraitBeast),
		string(card.TraitWarrior),
		string(card.TraitDragon),
		string(card.TraitDemon),
		string(card.TraitAngel),
	}

	for _, t := range traits {
		_, err := s.pool.Exec(
			context.Background(),
			`INSERT INTO traits (name) VALUES ($1) 
			 ON CONFLICT (name) DO NOTHING`,
			t,
		)
		if err != nil {
			return fmt.Errorf("failed to seed trait %s: %w", t, err)
		}
	}

	// Create some keywords
	keywords := []string{
		"HASTE",
		"CRITICAL",
		"EQUIPMENT",
		"DAMAGE",
		"BUFF",
	}

	for _, k := range keywords {
		_, err := s.pool.Exec(
			context.Background(),
			`INSERT INTO keywords (name) VALUES ($1) 
			 ON CONFLICT (name) DO NOTHING`,
			k,
		)
		if err != nil {
			return fmt.Errorf("failed to seed keyword %s: %w", k, err)
		}
	}

	// Create sample cards
	// 1. Create a creature
	creature := &card.Creature{
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
	}

	_, err := s.Save(creature)
	if err != nil {
		return fmt.Errorf("failed to seed creature: %w", err)
	}

	// 2. Create a spell
	spell := &card.Spell{
		BaseCard: card.BaseCard{
			Name:     "Lightning Bolt",
			Cost:     1,
			Effect:   "Deal 3 damage to target creature.",
			Type:     card.TypeSpell,
			Keywords: []string{"DAMAGE"},
		},
		TargetType: "Creature",
	}

	_, err = s.Save(spell)
	if err != nil {
		return fmt.Errorf("failed to seed spell: %w", err)
	}

	// 3. Create an artifact
	artifact := &card.Artifact{
		BaseCard: card.BaseCard{
			Name:     "Dragon Sword",
			Cost:     3,
			Effect:   "Equip to a creature. Equipped creature gets +2/+0.",
			Type:     card.TypeArtifact,
			Keywords: []string{"EQUIPMENT"},
		},
		IsEquipment: true,
	}

	_, err = s.Save(artifact)
	if err != nil {
		return fmt.Errorf("failed to seed artifact: %w", err)
	}

	log.Println("Test data seeded successfully")
	return nil
}

// GetPool returns the underlying connection pool
// This is used by utilities like the database cleaner
func (s *PostgresStore) GetPool() *pgxpool.Pool {
	return s.pool
}
