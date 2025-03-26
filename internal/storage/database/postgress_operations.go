package database

import (
	"context"
	"fmt"
	"time"

	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
	"github.com/jackc/pgx/v5"
)

// saveCard stores a card in the database
func (s *PostgresStore) saveCard(c card.Card) (string, error) {
	// Validate the card first
	if err := c.Validate(); err != nil {
		return "", fmt.Errorf("invalid card: %w", err)
	}
	
	// Start a transaction
	tx, err := s.pool.Begin(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()
	
	// Get card data
	data := c.ToDTO()
	
	// Get card type ID
	var typeID int
	err = tx.QueryRow(
		context.Background(),
		`SELECT id FROM card_types WHERE name = $1`,
		data.Type,
	).Scan(&typeID)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			// Insert card type if not exists
			err = tx.QueryRow(
				context.Background(),
				`INSERT INTO card_types (name) VALUES ($1) RETURNING id`,
				data.Type,
			).Scan(&typeID)
			
			if err != nil {
				return "", fmt.Errorf("failed to create card type: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to get card type: %w", err)
		}
	}
	
	// Insert card
	var cardID int
	now := time.Now()
	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO cards (name, cost, effect, type_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`,
		data.Name, data.Cost, data.Effect, typeID, now, now,
	).Scan(&cardID)
	
	if err != nil {
		return "", fmt.Errorf("failed to insert card: %w", err)
	}
	
	// Insert type-specific data based on card type
	if err = s.saveTypeSpecificData(tx, cardID, data); err != nil {
		return "", err
	}
	
	// Insert keywords
	if err = s.saveCardKeywords(tx, cardID, data.Keywords); err != nil {
		return "", err
	}
	
	// Insert metadata
	if err = s.saveCardMetadata(tx, cardID, data.Metadata); err != nil {
		return "", err
	}
	
	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return fmt.Sprintf("%d", cardID), nil
}

// saveTypeSpecificData saves the type-specific data for a card
func (s *PostgresStore) saveTypeSpecificData(tx pgx.Tx, cardID int, data *card.CardDTO) error {
	switch data.Type {
	case card.TypeCreature:
		return s.saveCreatureData(tx, cardID, data)
	case card.TypeArtifact:
		return s.saveArtifactData(tx, cardID, data)
	case card.TypeSpell:
		return s.saveSpellData(tx, cardID, data)
	case card.TypeIncantation:
		return s.saveIncantationData(tx, cardID, data)
	case card.TypeAnthem:
		return s.saveAnthemData(tx, cardID, data)
	default:
		return fmt.Errorf("unsupported card type: %s", data.Type)
	}
}

// saveCreatureData saves creature-specific data
func (s *PostgresStore) saveCreatureData(tx pgx.Tx, cardID int, data *card.CardDTO) error {
	if data.Trait == "" {
		// If no trait, just insert creature data without trait_id
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO creature_cards (card_id, attack, defense)
			VALUES ($1, $2, $3)`,
			cardID, data.Attack, data.Defense,
		)
		if err != nil {
			return fmt.Errorf("failed to insert creature data: %w", err)
		}
		return nil
	}
	
	// Try to get trait ID
	var traitID int
	err := tx.QueryRow(
		context.Background(),
		`SELECT id FROM traits WHERE name = $1`,
		data.Trait,
	).Scan(&traitID)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			// Trait doesn't exist, create it
			err = tx.QueryRow(
				context.Background(),
				`INSERT INTO traits (name) VALUES ($1) RETURNING id`,
				data.Trait,
			).Scan(&traitID)
			
			if err != nil {
				return fmt.Errorf("failed to create trait: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get trait: %w", err)
		}
	}
	
	// Insert creature data with trait
	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO creature_cards (card_id, attack, defense, trait_id)
		VALUES ($1, $2, $3, $4)`,
		cardID, data.Attack, data.Defense, traitID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert creature data: %w", err)
	}
	
	return nil
}

// saveArtifactData saves artifact-specific data
func (s *PostgresStore) saveArtifactData(tx pgx.Tx, cardID int, data *card.CardDTO) error {
	_, err := tx.Exec(
		context.Background(),
		`INSERT INTO artifact_cards (card_id, is_equipment)
		VALUES ($1, $2)`,
		cardID, data.IsEquipment,
	)
	if err != nil {
		return fmt.Errorf("failed to insert artifact data: %w", err)
	}
	return nil
}

// saveSpellData saves spell-specific data
func (s *PostgresStore) saveSpellData(tx pgx.Tx, cardID int, data *card.CardDTO) error {
	_, err := tx.Exec(
		context.Background(),
		`INSERT INTO spell_cards (card_id, target_type)
		VALUES ($1, $2)`,
		cardID, data.TargetType,
	)
	if err != nil {
		return fmt.Errorf("failed to insert spell data: %w", err)
	}
	return nil
}

// saveIncantationData saves incantation-specific data
func (s *PostgresStore) saveIncantationData(tx pgx.Tx, cardID int, data *card.CardDTO) error {
	_, err := tx.Exec(
		context.Background(),
		`INSERT INTO incantation_cards (card_id, timing)
		VALUES ($1, $2)`,
		cardID, data.Timing,
	)
	if err != nil {
		return fmt.Errorf("failed to insert incantation data: %w", err)
	}
	return nil
}

// saveAnthemData saves anthem-specific data
func (s *PostgresStore) saveAnthemData(tx pgx.Tx, cardID int, data *card.CardDTO) error {
	_, err := tx.Exec(
		context.Background(),
		`INSERT INTO anthem_cards (card_id, continuous)
		VALUES ($1, $2)`,
		cardID, data.Continuous,
	)
	if err != nil {
		return fmt.Errorf("failed to insert anthem data: %w", err)
	}
	return nil
}

// saveCardKeywords saves the keywords for a card
func (s *PostgresStore) saveCardKeywords(tx pgx.Tx, cardID int, keywords []string) error {
	for _, keyword := range keywords {
		// Get keyword ID or create it
		var keywordID int
		err := tx.QueryRow(
			context.Background(),
			`SELECT id FROM keywords WHERE name = $1`,
			keyword,
		).Scan(&keywordID)
		
		if err != nil {
			if err == pgx.ErrNoRows {
				err = tx.QueryRow(
					context.Background(),
					`INSERT INTO keywords (name) VALUES ($1) RETURNING id`,
					keyword,
				).Scan(&keywordID)
				
				if err != nil {
					return fmt.Errorf("failed to create keyword: %w", err)
				}
			} else {
				return fmt.Errorf("failed to get keyword: %w", err)
			}
		}
		
		// Insert keyword relation
		_, err = tx.Exec(
			context.Background(),
			`INSERT INTO card_keywords (card_id, keyword_id)
			VALUES ($1, $2)`,
			cardID, keywordID,
		)
		
		if err != nil {
			return fmt.Errorf("failed to insert keyword relation: %w", err)
		}
	}
	return nil
}

// saveCardMetadata saves the metadata for a card
func (s *PostgresStore) saveCardMetadata(tx pgx.Tx, cardID int, metadata map[string]string) error {
	for key, value := range metadata {
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO card_metadata (card_id, key, value)
			VALUES ($1, $2, $3)`,
			cardID, key, value,
		)
		
		if err != nil {
			return fmt.Errorf("failed to insert metadata: %w", err)
		}
	}
	return nil
}

// loadCard loads a card from the database by ID
func (s *PostgresStore) loadCard(id string) (card.Card, error) {
	// Parse ID
	cardID, err := parseCardID(id)
	if err != nil {
		return nil, err
	}

	// Query for the base card data and card type
	var (
		dbID      int
		name      string
		cost      int
		effect    string
		typeID    int
		typeName  string
		createdAt time.Time
		updatedAt time.Time
	)
	
	err = s.pool.QueryRow(
		context.Background(),
		`SELECT c.id, c.name, c.cost, c.effect, c.type_id, ct.name, c.created_at, c.updated_at
		FROM cards c
		JOIN card_types ct ON c.type_id = ct.id
		WHERE c.id = $1`,
		cardID,
	).Scan(&dbID, &name, &cost, &effect, &typeID, &typeName, &createdAt, &updatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("card not found: %s", id)
		}
		return nil, fmt.Errorf("failed to load card: %w", err)
	}
	
	// Query for keywords
	keywords, err := s.loadCardKeywords(cardID)
	if err != nil {
		return nil, err
	}
	
	// Query for metadata
	metadata, err := s.loadCardMetadata(cardID)
	if err != nil {
		return nil, err
	}
	
	// Create base card
	baseCard := card.BaseCard{
		ID:        id,
		Name:      name,
		Cost:      cost,
		Effect:    effect,
		Type:      card.CardType(typeName),
		Keywords:  keywords,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Metadata:  metadata,
	}
	
	// Create type-specific card
	var c card.Card
	
	switch card.CardType(typeName) {
	case card.TypeCreature:
		c, err = s.loadCreatureCard(cardID, baseCard)
	case card.TypeArtifact:
		c, err = s.loadArtifactCard(cardID, baseCard)
	case card.TypeSpell:
		c, err = s.loadSpellCard(cardID, baseCard)
	case card.TypeIncantation:
		c, err = s.loadIncantationCard(cardID, baseCard)
	case card.TypeAnthem:
		c, err = s.loadAnthemCard(cardID, baseCard)
	default:
		// Default to a base card for unknown types
		c = &baseCard
	}
	
	if err != nil {
		return nil, err
	}
	
	return c, nil
}

// loadCardKeywords loads the keywords for a card
func (s *PostgresStore) loadCardKeywords(cardID int) ([]string, error) {
	rows, err := s.pool.Query(
		context.Background(),
		`SELECT k.name
		FROM card_keywords ck
		JOIN keywords k ON ck.keyword_id = k.id
		WHERE ck.card_id = $1`,
		cardID,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to load keywords: %w", err)
	}
	defer rows.Close()
	
	var keywords []string
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, fmt.Errorf("failed to scan keyword: %w", err)
		}
		keywords = append(keywords, keyword)
	}
	
	return keywords, nil
}

// loadCardMetadata loads the metadata for a card
func (s *PostgresStore) loadCardMetadata(cardID int) (map[string]string, error) {
	rows, err := s.pool.Query(
		context.Background(),
		`SELECT key, value
		FROM card_metadata
		WHERE card_id = $1`,
		cardID,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}
	defer rows.Close()
	
	metadata := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("failed to scan metadata: %w", err)
		}
		metadata[key] = value
	}
	
	return metadata, nil
}

// loadCreatureCard loads a creature card from the database
func (s *PostgresStore) loadCreatureCard(cardID int, baseCard card.BaseCard) (card.Card, error) {
	var (
		attack  int
		defense int
		trait   pgx.NullString
	)
	
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT cc.attack, cc.defense, t.name
		FROM creature_cards cc
		LEFT JOIN traits t ON cc.trait_id = t.id
		WHERE cc.card_id = $1`,
		cardID,
	).Scan(&attack, &defense, &trait)
	
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to load creature data: %w", err)
	}
	
	traitStr := ""
	if trait.Valid {
		traitStr = trait.String
	}
	
	return &card.Creature{
		BaseCard: baseCard,
		Attack:   attack,
		Defense:  defense,
		Trait:    card.Trait(traitStr),
	}, nil
}

// loadArtifactCard loads an artifact card from the database
func (s *PostgresStore) loadArtifactCard(cardID int, baseCard card.BaseCard) (card.Card, error) {
	var isEquipment bool
	
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT is_equipment
		FROM artifact_cards
		WHERE card_id = $1`,
		cardID,
	).Scan(&isEquipment)
	
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to load artifact data: %w", err)
	}
	
	return &card.Artifact{
		BaseCard:    baseCard,
		IsEquipment: isEquipment,
	}, nil
}

// loadSpellCard loads a spell card from the database
func (s *PostgresStore) loadSpellCard(cardID int, baseCard card.BaseCard) (card.Card, error) {
	var targetType pgx.NullString
	
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT target_type
		FROM spell_cards
		WHERE card_id = $1`,
		cardID,
	).Scan(&targetType)
	
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to load spell data: %w", err)
	}
	
	targetTypeStr := ""
	if targetType.Valid {
		targetTypeStr = targetType.String
	}
	
	return &card.Spell{
		BaseCard:   baseCard,
		TargetType: targetTypeStr,
	}, nil
}

// loadIncantationCard loads an incantation card from the database
func (s *PostgresStore) loadIncantationCard(cardID int, baseCard card.BaseCard) (card.Card, error) {
	var timing pgx.NullString
	
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT timing
		FROM incantation_cards
		WHERE card_id = $1`,
		cardID,
	).Scan(&timing)
	
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to load incantation data: %w", err)
	}
	
	timingStr := ""
	if timing.Valid {
		timingStr = timing.String
	}
	
	return &card.Incantation{
		BaseCard: baseCard,
		Timing:   timingStr,
	}, nil
}

// loadAnthemCard loads an anthem card from the database
func (s *PostgresStore) loadAnthemCard(cardID int, baseCard card.BaseCard) (card.Card, error) {
	var continuous bool
	
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT continuous
		FROM anthem_cards
		WHERE card_id = $1`,
		cardID,
	).Scan(&continuous)
	
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to load anthem data: %w", err)
	}
	
	return &card.Anthem{
		BaseCard:   baseCard,
		Continuous: continuous,
	}, nil
}

// parseCardID converts a string ID to an integer
func parseCardID(id string) (int, error) {
	var cardID int
	var err error
	
	// Try to parse as integer
	cardID, err = stringToInt(id)
	if err != nil {
		return 0, fmt.Errorf("invalid card ID: %s", id)
	}
	
	return cardID, nil
}

// stringToInt tries to convert a string to an integer
func stringToInt(s string) (int, error) {
	// Try to parse as integer
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// deleteCard deletes a card and its related data
func (s *PostgresStore) deleteCard(id string) error {
	// Parse ID
	cardID, err := parseCardID(id)
	if err != nil {
		return err
	}
	
	// Start a transaction
	tx, err := s.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()
	
	// Delete related records first (foreign key constraints)
	tables := []string{
		"card_keywords",
		"card_metadata",
		"creature_cards",
		"artifact_cards",
		"spell_cards",
		"incantation_cards",
		"anthem_cards",
		"card_images",
		"card_set_cards",
	}
	
	for _, table := range tables {
		_, err = tx.Exec(
			context.Background(),
			fmt.Sprintf("DELETE FROM %s WHERE card_id = $1", table),
			cardID,
		)
		if err != nil {
			return fmt.Errorf("failed to delete from %s: %w", table, err)
		}
	}
	
	// Delete the card itself
	commandTag, err := tx.Exec(
		context.Background(),
		"DELETE FROM cards WHERE id = $1",
		cardID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}
	
	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("card not found: %s", id)
	}
	
	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// listCards returns all cards in the database
func (s *PostgresStore) listCards() ([]card.Card, error) {
	// Query for all card IDs
	rows, err := s.pool.Query(
		context.Background(),
		`SELECT id FROM cards ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list cards: %w", err)
	}
	defer rows.Close()
	
	var ids []string
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan id: %w", err)
		}
		ids = append(ids, fmt.Sprintf("%d", id))
	}
	
	var cards []card.Card
	for _, id := range ids {
		card, err := s.loadCard(id)
		if err != nil {
			return nil, fmt.Errorf("failed to load card %s: %w", id, err)
		}
		cards = append(cards, card)
	}
	
	return cards, nil
}