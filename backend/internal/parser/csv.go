package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// CSVParser parses CSV files into card domain models
type CSVParser struct {
	reader *csv.Reader
}

// NewCSVParser creates a new CSV parser
func NewCSVParser(r io.Reader) *CSVParser {
	parser := &CSVParser{
		reader: csv.NewReader(r),
	}
	// Allow variable number of fields per record
	parser.reader.FieldsPerRecord = -1
	return parser
}

// ParseCSV parses a CSV file into a slice of cards based on type
// Returns cards implementing the card.Card interface
func (p *CSVParser) ParseCSV(cardType string) ([]card.Card, error) {
	// Read header
	header, err := p.reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Create a map of column indices
	colIndex := make(map[string]int)
	for i, col := range header {
		colIndex[strings.TrimSpace(col)] = i
	}

	// Required columns for all card types
	requiredColumns := []string{"Name", "Cost", "Effect"}
	for _, col := range requiredColumns {
		if _, exists := colIndex[col]; !exists {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	var cards []card.Card
	lineNum := 1 // Start after header
	for {
		record, err := p.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read line %d: %w", lineNum, err)
		}

		// Parse the card based on the specified type
		c, err := p.parseCardByType(record, colIndex, cardType)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
		cards = append(cards, c)
		lineNum++
	}

	return cards, nil
}

// parseCardByType parses a CSV record into a specific card type
func (p *CSVParser) parseCardByType(record []string, colIndex map[string]int, cardType string) (card.Card, error) {
	// Helper function to safely get a value from the record
	getValue := func(columnName string) string {
		if idx, exists := colIndex[columnName]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Get common fields
	name := getValue("Name")
	effect := getValue("Effect")
	costStr := getValue("Cost")

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if effect == "" {
		return nil, fmt.Errorf("effect is required")
	}
	if costStr == "" {
		return nil, fmt.Errorf("cost is required")
	}

	// Parse cost as integer
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		return nil, fmt.Errorf("invalid cost '%s': %w", costStr, err)
	}

	// Parse card based on type
	switch strings.ToLower(cardType) {
	case "anthem":
		return p.parseAnthem(name, cost, effect)
	case "creature":
		return p.parseCreature(name, cost, effect, record, colIndex)
	case "spell":
		return p.parseSpell(name, cost, effect)
	case "artifact":
		return p.parseArtifact(name, cost, effect)
	case "incantation":
		return p.parseIncantation(name, cost, effect)
	default:
		return nil, fmt.Errorf("unsupported card type: %s", cardType)
	}
}

// parseAnthem creates an Anthem card
func (p *CSVParser) parseAnthem(name string, cost int, effect string) (card.Card, error) {
	// Create base card
	baseCard := card.BaseCard{
		Name:     name,
		Cost:     cost,
		Effect:   effect,
		Type:     card.TypeAnthem,
		Keywords: extractKeywords(effect),
	}

	// Create anthem
	anthem := &card.Anthem{
		BaseCard:   baseCard,
		Continuous: true, // Anthems are always continuous
	}

	return anthem, nil
}

// parseCreature creates a Creature card
func (p *CSVParser) parseCreature(name string, cost int, effect string, record []string, colIndex map[string]int) (card.Card, error) {
	// Helper function to safely get a value
	getValue := func(columnName string) string {
		if idx, exists := colIndex[columnName]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Get creature-specific fields
	attackStr := getValue("Attack")
	defenseStr := getValue("Defense")
	trait := getValue("Trait")

	if attackStr == "" {
		return nil, fmt.Errorf("attack is required for creatures")
	}
	if defenseStr == "" {
		return nil, fmt.Errorf("defense is required for creatures")
	}

	// Parse attack and defense
	attack, err := strconv.Atoi(attackStr)
	if err != nil {
		return nil, fmt.Errorf("invalid attack '%s': %w", attackStr, err)
	}
	defense, err := strconv.Atoi(defenseStr)
	if err != nil {
		return nil, fmt.Errorf("invalid defense '%s': %w", defenseStr, err)
	}

	// Create creature
	baseCard := card.BaseCard{
		Name:     name,
		Cost:     cost,
		Effect:   effect,
		Type:     card.TypeCreature,
		Keywords: extractKeywords(effect),
	}

	creature := &card.Creature{
		BaseCard: baseCard,
		Attack:   attack,
		Defense:  defense,
		Trait:    card.Trait(trait),
	}

	return creature, nil
}

// parseSpell creates a Spell card
func (p *CSVParser) parseSpell(name string, cost int, effect string) (card.Card, error) {
	// Create base card
	baseCard := card.BaseCard{
		Name:     name,
		Cost:     cost,
		Effect:   effect,
		Type:     card.TypeSpell,
		Keywords: extractKeywords(effect),
	}

	// Create spell with target type determined from effect text
	spell := &card.Spell{
		BaseCard:   baseCard,
		TargetType: card.DetermineTargetType(effect),
	}

	return spell, nil
}

// parseArtifact creates an Artifact card
func (p *CSVParser) parseArtifact(name string, cost int, effect string) (card.Card, error) {
	// Create base card
	baseCard := card.BaseCard{
		Name:     name,
		Cost:     cost,
		Effect:   effect,
		Type:     card.TypeArtifact,
		Keywords: extractKeywords(effect),
	}

	// Determine if this is equipment based on effect text
	isEquipment := card.DetermineIsEquipment(effect)

	// Create artifact
	artifact := &card.Artifact{
		BaseCard:    baseCard,
		IsEquipment: isEquipment,
	}

	return artifact, nil
}

// parseIncantation creates an Incantation card
func (p *CSVParser) parseIncantation(name string, cost int, effect string) (card.Card, error) {
	// Create base card
	baseCard := card.BaseCard{
		Name:     name,
		Cost:     cost,
		Effect:   effect,
		Type:     card.TypeIncantation,
		Keywords: extractKeywords(effect),
	}

	// Create incantation with timing determined from effect text
	incantation := &card.Incantation{
		BaseCard: baseCard,
		Timing:   card.DetermineTiming(effect),
	}

	return incantation, nil
}

// extractKeywords extracts keywords from effect text
// This is a simple implementation - could be improved with more advanced NLP
func extractKeywords(effect string) []string {
	commonKeywords := []string{
		"CRITICAL", "HASTE", "DAMAGE", "BUFF", "EQUIPMENT",
		"COUNTER", "DRAW", "DIRECT", "FLYING", "IMMUNE",
	}

	var found []string
	effectUpper := strings.ToUpper(effect)

	for _, keyword := range commonKeywords {
		if strings.Contains(effectUpper, keyword) {
			found = append(found, keyword)
		}
	}

	return found
}
