package card

import (
	"encoding/json"
	"fmt"
)

// CardData represents the serializable form of a card
// This is used for storing/loading cards
type CardData struct {
	Type        CardType            `json:"type"`
	Name        string             `json:"name"`
	Cost        int                `json:"cost"`
	Effect      string             `json:"effect"`
	Attack      int                `json:"attack,omitempty"`
	Defense     int                `json:"defense,omitempty"`
	SubType     string             `json:"subtype,omitempty"`
	IsEquipment bool               `json:"is_equipment,omitempty"`
	TargetType  string             `json:"target_type,omitempty"`
	Timing      string             `json:"timing,omitempty"`
	Continuous  bool               `json:"continuous,omitempty"`
	Keywords    []string           `json:"keywords,omitempty"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
}

// CardStore interface defines methods for storing and retrieving cards
type CardStore interface {
	Save(card Card) error
	Load(id string) (Card, error)
	List() ([]Card, error)
	Delete(id string) error
}

// CardFactory handles the creation of different card types
type CardFactory struct {
	store CardStore
}

// NewCardFactory creates a new card factory
func NewCardFactory(store CardStore) *CardFactory {
	return &CardFactory{
		store: store,
	}
}

// CreateFromData creates a card from CardData
func (f *CardFactory) CreateFromData(data *CardData) (Card, error) {
	baseCard := BaseCard{
		Name:   data.Name,
		Cost:   data.Cost,
		Effect: data.Effect,
		Type:   data.Type,
	}

	var card Card
	switch data.Type {
	case TypeCreature:
		card = &Creature{
			BaseCard: baseCard,
			Attack:   data.Attack,
			Defense:  data.Defense,
			SubType:  data.SubType,
		}
	case TypeArtifact:
		card = &Artifact{
			BaseCard:    baseCard,
			IsEquipment: data.IsEquipment,
		}
	case TypeSpell:
		card = &Spell{
			BaseCard:   baseCard,
			TargetType: data.TargetType,
		}
	case TypeIncantation:
		card = &Incantation{
			BaseCard: baseCard,
			Timing:   data.Timing,
		}
	case TypeAnthem:
		card = &Anthem{
			BaseCard:    baseCard,
			Continuous: true,
		}
	default:
		return nil, fmt.Errorf("unknown card type: %s", data.Type)
	}

	// Validate the card before returning
	if err := card.Validate(); err != nil {
		return nil, fmt.Errorf("invalid card data: %w", err)
	}

	return card, nil
}

// ToData converts a Card interface to CardData
func ToData(card Card) *CardData {
	data := &CardData{
		Type:   card.GetType(),
		Name:   card.GetName(),
		Cost:   card.GetCost(),
		Effect: card.GetEffect(),
	}

	// Type-specific conversions
	switch c := card.(type) {
	case *Creature:
		data.Attack = c.Attack
		data.Defense = c.Defense
		data.SubType = c.SubType
	case *Artifact:
		data.IsEquipment = c.IsEquipment
	case *Spell:
		data.TargetType = c.TargetType
	case *Incantation:
		data.Timing = c.Timing
	case *Anthem:
		data.Continuous = c.Continuous
	}

	return data
}

// Serialization methods for CardData
func (d *CardData) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func (d *CardData) FromJSON(data []byte) error {
	return json.Unmarshal(data, d)
}