package database

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
    
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
    "github.com/lib/pq" // PostgreSQL driver
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
    // For initial testing, implement a very basic version
    data := c.ToData()
    
    var id int
    err := s.db.QueryRow(
        `INSERT INTO cards (name, cost, effect, card_type) 
         VALUES ($1, $2, $3, $4) 
         RETURNING id`,
        data.Name, data.Cost, data.Effect, data.Type,
    ).Scan(&id)
    
    if err != nil {
        return "", fmt.Errorf("failed to save card: %w", err)
    }
    
    return fmt.Sprintf("%d", id), nil
}

// Load retrieves a card by its ID
func (s *PostgresStore) Load(id string) (card.Card, error) {
    // For initial testing, implement a very basic version
    var (
        name     string
        cost     int
        effect   string
        cardType string
    )
    
    err := s.db.QueryRow(
        `SELECT name, cost, effect, card_type FROM cards WHERE id = $1`,
        id,
    ).Scan(&name, &cost, &effect, &cardType)
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("card not found: %s", id)
        }
        return nil, fmt.Errorf("failed to load card: %w", err)
    }
    
    // For initial testing, create a basic card
    baseCard := card.BaseCard{
        Name:   name,
        Cost:   cost,
        Effect: effect,
        Type:   card.CardType(cardType),
    }
    
    // Handle different card types with minimal implementation
    switch card.CardType(cardType) {
    case card.TypeCreature:
        return &card.Creature{BaseCard: baseCard}, nil
    case card.TypeSpell:
        return &card.Spell{BaseCard: baseCard}, nil
    case card.TypeArtifact:
        return &card.Artifact{BaseCard: baseCard}, nil
    default:
        // Default to a generic card for testing
        return &card.Creature{BaseCard: baseCard}, nil
    }
}

// List returns all stored cards
func (s *PostgresStore) List() ([]card.Card, error) {
    // For initial testing, implement a very basic version
    rows, err := s.db.Query(`SELECT id FROM cards`)
    if err != nil {
        return nil, fmt.Errorf("failed to list cards: %w", err)
    }
    defer rows.Close()
    
    var ids []string
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, fmt.Errorf("failed to scan id: %w", err)
        }
        ids = append(ids, id)
    }
    
    var cards []card.Card
    for _, id := range ids {
        card, err := s.Load(id)
        if err != nil {
            return nil, fmt.Errorf("failed to load card %s: %w", id, err)
        }
        cards = append(cards, card)
    }
    
    return cards, nil
}

// Delete removes a card by its ID
func (s *PostgresStore) Delete(id string) error {
    // For initial testing, implement a very basic version
    _, err := s.db.Exec(`DELETE FROM cards WHERE id = $1`, id)
    if err != nil {
        return fmt.Errorf("failed to delete card: %w", err)
    }
    
    return nil
}

// Close cleans up any resources
func (s *PostgresStore) Close() error {
    return s.db.Close()
}

// InitSchema sets up the necessary database schema
func (s *PostgresStore) InitSchema() error {
    // Basic schema for testing
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS cards (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            cost INTEGER NOT NULL,
            effect TEXT NOT NULL,
            card_type VARCHAR(50) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
    
    if err != nil {
        return fmt.Errorf("failed to create schema: %w", err)
    }
    
    log.Println("Database schema initialized")
    return nil
}