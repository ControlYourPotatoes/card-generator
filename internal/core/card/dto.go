package card

import (
	"time"
)

// CardDTO represents a data transfer object for cards
// Used for serialization, database operations, and API responses
type CardDTO struct {
	ID          string            `json:"id,omitempty"`
	Type        CardType          `json:"type"`
	Name        string            `json:"name"`
	Cost        int               `json:"cost"`
	Effect      string            `json:"effect"`
	Keywords    []string          `json:"keywords,omitempty"`
	
	// Type-specific fields
	Attack      int               `json:"attack,omitempty"`
	Defense     int               `json:"defense,omitempty"`
	Trait       string            `json:"trait,omitempty"`
	IsEquipment bool              `json:"is_equipment,omitempty"`
	TargetType  string            `json:"target_type,omitempty"`
	Timing      string            `json:"timing,omitempty"`
	Continuous  bool              `json:"continuous,omitempty"`
	
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ToDTO converts BaseCard to CardDTO
func (b BaseCard) ToDTO() *CardDTO {
	return &CardDTO{
		ID:        b.ID,
		Type:      b.Type,
		Name:      b.Name,
		Cost:      b.Cost,
		Effect:    b.Effect,
		Keywords:  b.Keywords,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		Metadata:  b.Metadata,
	}
}

// For backward compatibility - alias ToData to ToDTO
func (b BaseCard) ToData() *CardDTO {
	return b.ToDTO()
}

// NewBaseCardFromDTO creates a BaseCard from a CardDTO
func NewBaseCardFromDTO(dto *CardDTO) BaseCard {
	return BaseCard{
		ID:        dto.ID,
		Name:      dto.Name,
		Cost:      dto.Cost,
		Effect:    dto.Effect,
		Type:      dto.Type,
		Keywords:  dto.Keywords,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		Metadata:  dto.Metadata,
	}
}