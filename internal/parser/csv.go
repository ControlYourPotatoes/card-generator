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
    return &Parser{
        reader: csv.NewReader(r),
    }
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
    cardType := card.CardType(record[colIndex["Type"]])
    
    baseCard := card.BaseCard{
        Name:   record[colIndex["Name"]],
        Effect: record[colIndex["Effect"]],
        Type:   cardType,
    }

    // Check if Cost column exists
    if _, exists := colIndex["Cost"]; !exists {
        return nil, fmt.Errorf("missing required column: Cost")
    }

    cost, err := strconv.Atoi(record[colIndex["Cost"]])
    if err != nil {
        return nil, fmt.Errorf("invalid cost: not a number")
    }
    baseCard.Cost = cost

    switch cardType {
    case card.TypeCreature:
        // Check required columns for creature
        if _, exists := colIndex["Attack"]; !exists {
            return nil, fmt.Errorf("missing required column: Attack")
        }
        if _, exists := colIndex["Defense"]; !exists {
            return nil, fmt.Errorf("missing required column: Defense")
        }

        attack, err := strconv.Atoi(record[colIndex["Attack"]])
        if err != nil {
            return nil, fmt.Errorf("invalid attack: not a number")
        }
        defense, err := strconv.Atoi(record[colIndex["Defense"]])
        if err != nil {
            return nil, fmt.Errorf("invalid defense: not a number")
        }
        return &card.Creature{
            BaseCard: baseCard,
            Attack:   attack,
            Defense:  defense,
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