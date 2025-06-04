package migration

import (
	"io/ioutil"
	"path/filepath"
)

// GetInitialSchema reads the initial schema from the first migration file
func GetInitialSchema() (string, error) {
	// For simplicity, we're just returning the hardcoded schema
	// In a production system, this would read from the migration files
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
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_cards_updated_at') THEN
        CREATE TRIGGER update_cards_updated_at
        BEFORE UPDATE ON cards
        FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
`, nil
}

// ReadMigrationFile reads a migration file from disk
func ReadMigrationFile(migrationsDir string, filename string) (string, error) {
	path := filepath.Join(migrationsDir, filename)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
