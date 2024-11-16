package parser

import (
    "encoding/csv"
    "fmt"
    "io"
    "strconv"
    "/internal/card"
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

    var cards []card.Card
    for {
        record, err := p.reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("failed to read record: %w", err)
        }

        // Parse the card based on its type
        card, err := p.parseRecord(record, colIndex)
        if err != nil {
            return nil, fmt.Errorf("failed to parse record: %w", err)
        }
        cards = append(cards, card)
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

    cost, err := strconv.Atoi(record[colIndex["Cost"]])
    if err != nil {
        return nil, fmt.Errorf("invalid cost: %w", err)
    }
    baseCard.Cost = cost

    switch cardType {
    case card.TypeCreature:
        attack, err := strconv.Atoi(record[colIndex["Attack"]])
        if err != nil {
            return nil, fmt.Errorf("invalid attack: %w", err)
        }
        defense, err := strconv.Atoi(record[colIndex["Defense"]])
        if err != nil {
            return nil, fmt.Errorf("invalid defense: %w", err)
        }
        return &card.Creature{
            BaseCard: baseCard,
            Attack:   attack,
            Defense:  defense,
        }, nil
    // Add other card types here...
    default:
        return &baseCard, nil
    }
}