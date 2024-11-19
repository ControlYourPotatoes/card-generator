package parser

import (
    "encoding/csv"
    "fmt"
    "io"
    "strconv"
    "strings"
    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

type Parser struct {
    reader *csv.Reader
}

func NewParser(r io.Reader) *Parser {
    parser := &Parser{
        reader: csv.NewReader(r),
    }
    // Allow variable number of fields per record
    parser.reader.FieldsPerRecord = -1
    return parser
}
func (p *Parser) Parse() ([]card.Card, error) {
    // Read header
    header, err := p.reader.Read()
    if err != nil {
        return nil, fmt.Errorf("failed to read header: %w", err)
    }

    // Create a map of column indices
    colIndex := make(map[string]int)
    for i, col := range header {
        colIndex[col] = i
    }

    // Validate required columns
    requiredColumns := []string{"Type", "Name", "Cost", "Effect"}
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

        // Parse the card based on its type
        card, err := p.parseRecord(record, colIndex)
        if err != nil {
            return nil, fmt.Errorf("line %d: %w", lineNum, err)
        }
        cards = append(cards, card)
        lineNum++
    }

    return cards, nil
}

func (p *Parser) parseRecord(record []string, colIndex map[string]int) (card.Card, error) {
    // Safely get value from record with bounds checking
    getValue := func(columnName string) string {
        if idx, exists := colIndex[columnName]; exists && idx < len(record) {
            return record[idx]
        }
        return ""
    }

    // Get required fields
    cardType := card.CardType(getValue("Type"))
    name := getValue("Name")
    effect := getValue("Effect")
    
    if name == "" || effect == "" {
        return nil, fmt.Errorf("missing required fields name or effect")
    }

    // Parse cost
    costStr := getValue("Cost")
    if costStr == "" {
        return nil, fmt.Errorf("missing required field: Cost")
    }
    cost, err := strconv.Atoi(costStr)
    if err != nil {
        return nil, fmt.Errorf("invalid cost: not a number")
    }

    baseCard := card.BaseCard{
        Name:   name,
        Effect: effect,
        Type:   cardType,
        Cost:   cost,
    }

    switch cardType {
    case card.TypeCreature:
        // Parse Attack
        attackStr := getValue("Attack")
        if attackStr == "" {
            return nil, fmt.Errorf("missing required field Attack for Creature")
        }
        attack, err := strconv.Atoi(attackStr)
        if err != nil {
            return nil, fmt.Errorf("invalid attack: not a number")
        }

        // Parse Defense
        defenseStr := getValue("Defense")
        if defenseStr == "" {
            return nil, fmt.Errorf("missing required field Defense for Creature")
        }
        defense, err := strconv.Atoi(defenseStr)
        if err != nil {
            return nil, fmt.Errorf("invalid defense: not a number")
        }

        // Get trait (optional)
        trait := getValue("Trait")

        return &card.Creature{
            BaseCard: baseCard,
            Attack:   attack,
            Defense:  defense,
            Trait:    trait,
        }, nil

    case card.TypeArtifact:
        return &card.Artifact{
            BaseCard:    baseCard,
            IsEquipment: strings.Contains(strings.ToLower(baseCard.Effect), "equip"),
        }, nil

    case card.TypeSpell:
        return &card.Spell{
            BaseCard:   baseCard,
            TargetType: determineTargetType(baseCard.Effect),
        }, nil

    case card.TypeIncantation:
        return &card.Incantation{
            BaseCard: baseCard,
            Timing:   determineTiming(baseCard.Effect),
        }, nil

    case card.TypeAnthem:
        return &card.Anthem{
            BaseCard:    baseCard,
            Continuous: true,
        }, nil

    default:
        return nil, fmt.Errorf("unknown card type: %s", cardType)
    }
}
// Helper functions remain the same
func determineTargetType(effect string) string {
    effect = strings.ToLower(effect)
    if strings.Contains(effect, "target creature") {
        return "Creature"
    }
    if strings.Contains(effect, "target player") {
        return "Player"
    }
    return "Any"
}

func determineTiming(effect string) string {
    if strings.Contains(effect, "ON ANY CLASH") {
        return "ON ANY CLASH"
    }
    if strings.Contains(effect, "ON ATTACK") {
        return "ON ATTACK"
    }
    return ""
}