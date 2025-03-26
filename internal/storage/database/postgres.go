package database

import (
	"database/sql"
	"fmt"
	"log"
	
	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresStore implements the Store interface using PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(connString string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &PostgresStore{db: db}, nil
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
	return s.db.Close()
}

// InitSchema sets up the necessary database schema
func (s *PostgresStore) InitSchema() error {
	// Read the schema definition from the migration file
	schema := getInitialSchema()
	
	// Execute the schema
	_, err := s.db.Exec(schema)
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
		_, err := s.db.Exec(
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
		_, err := s.db.Exec(
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
		_, err := s.db.Exec(
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

// getInitialSchema returns the SQL for initializing the database schema
func getInitialSchema() string {
	return `
-- Card Types Table
CREATE TABLE IF NOT EXISTS card_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- Keywords Table
CREATE TABLE IF NOT EXISTS keywords (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- Traits Table (for creature traits)
CREATE TABLE IF NOT EXISTS traits (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- Base Cards Table
CREATE TABLE IF NOT EXISTS cards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cost INTEGER NOT NULL CHECK (cost >= -1), -- -1 allowed for X costs
    effect TEXT NOT NULL,
    type_id INTEGER NOT NULL REFERENCES card_types(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Card Keywords Junction Table
CREATE TABLE IF NOT EXISTS card_keywords (
    card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    keyword_id INTEGER NOT NULL REFERENCES keywords(id),
    PRIMARY KEY (card_id, keyword_id)
);

-- Card Metadata Table (for flexible key-value pairs)
CREATE TABLE IF NOT EXISTS card_metadata (
    card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    PRIMARY KEY (card_id, key)
);

-- Creature Cards Table
CREATE TABLE IF NOT EXISTS creature_cards (
    card_id INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    attack INTEGER NOT NULL CHECK (attack >= 0),
    defense INTEGER NOT NULL CHECK (defense >= 0),
    trait_id INTEGER REFERENCES traits(id)
);

-- Artifact Cards Table
CREATE TABLE IF NOT EXISTS artifact_cards (
    card_id INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    is_equipment BOOLEAN NOT NULL DEFAULT FALSE
);

-- Spell Cards Table
CREATE TABLE IF NOT EXISTS spell_cards (
    card_id INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    target_type VARCHAR(50)
);

-- Incantation Cards Table
CREATE TABLE IF NOT EXISTS incantation_cards (
    card_id INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    timing VARCHAR(50)
);

-- Anthem Cards Table
CREATE TABLE IF NOT EXISTS anthem_cards (
    card_id INTEGER PRIMARY KEY REFERENCES cards(id) ON DELETE CASCADE,
    continuous BOOLEAN NOT NULL DEFAULT TRUE
);

-- Card Images Table
CREATE TABLE IF NOT EXISTS card_images (
    id SERIAL PRIMARY KEY,
    card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    image_path VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Card Sets Table (for grouping cards into expansions/sets)
CREATE TABLE IF NOT EXISTS card_sets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(10) NOT NULL UNIQUE,
    release_date DATE,
    description TEXT
);

-- Card Set Junction Table
CREATE TABLE IF NOT EXISTS card_set_cards (
    set_id INTEGER NOT NULL REFERENCES card_sets(id),
    card_id INTEGER NOT NULL REFERENCES cards(id),
    card_number VARCHAR(20) NOT NULL,
    rarity VARCHAR(20) NOT NULL,
    PRIMARY KEY (set_id, card_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_cards_type_id ON cards(type_id);
CREATE INDEX IF NOT EXISTS idx_card_keywords_card_id ON card_keywords(card_id);
CREATE INDEX IF NOT EXISTS idx_card_metadata_card_id ON card_metadata(card_id);
CREATE INDEX IF NOT EXISTS idx_card_set_cards_set_id ON card_set_cards(set_id);

-- Trigger function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$ language 'plpgsql';

-- Create trigger if it doesn't exist
DO $ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_cards_updated_at') THEN
        CREATE TRIGGER update_cards_updated_at
        BEFORE UPDATE ON cards
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $;
`
}