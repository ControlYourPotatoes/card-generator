package models

import (

	"github.com/ControlYourPotatoes/card-generator/internal/core/card"
)



// ToDomain converts database models to domain models
func (cm *CardModel) ToDomain(specificData TypeSpecificData, keywords []string, metadata map[string]string) (card.Card, error) {
	// Create base card
	baseCard := card.BaseCard{
		ID:        string(cm.ID),
		Name:      cm.Name,
		Cost:      cm.Cost,
		Effect:    cm.Effect,
		Type:      card.CardType(cm.TypeName),
		Keywords:  keywords,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
		Metadata:  metadata,
	}

	// Create specific card type based on TypeName
	var c card.Card

	switch card.CardType(cm.TypeName) {
	case card.TypeCreature:
		c = &card.Creature{
			BaseCard: baseCard,
			Attack:   specificData.Attack,
			Defense:  specificData.Defense,
			Trait:    card.Trait(specificData.TraitName),
		}
	case card.TypeArtifact:
		c = &card.Artifact{
			BaseCard:    baseCard,
			IsEquipment: specificData.IsEquipment,
		}
	case card.TypeSpell:
		c = &card.Spell{
			BaseCard:   baseCard,
			TargetType: specificData.TargetType,
		}
	case card.TypeIncantation:
		c = &card.Incantation{
			BaseCard: baseCard,
			Timing:   specificData.Timing,
		}
	case card.TypeAnthem:
		c = &card.Anthem{
			BaseCard:   baseCard,
			Continuous: specificData.Continuous,
		}
	default:
		// Default to a base implementation
		c = &baseCard
	}

	return c, nil
}

// FromDomain converts domain models to database models
func FromDomain(c card.Card) (*CardModel, *TypeSpecificData, []string, map[string]string, error) {
	// Extract card data
	dto := c.ToDTO()

	// Create CardModel
	cardModel := &CardModel{
		Name:      dto.Name,
		Cost:      dto.Cost,
		Effect:    dto.Effect,
		TypeName:  string(dto.Type),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}

	// Create type-specific data
	specificData := &TypeSpecificData{}

	// Set type-specific fields based on card type
	switch dto.Type {
	case card.TypeCreature:
		specificData.Attack = dto.Attack
		specificData.Defense = dto.Defense
		specificData.TraitName = dto.Trait
	case card.TypeArtifact:
		specificData.IsEquipment = dto.IsEquipment
	case card.TypeSpell:
		specificData.TargetType = dto.TargetType
	case card.TypeIncantation:
		specificData.Timing = dto.Timing
	case card.TypeAnthem:
		specificData.Continuous = dto.Continuous
	}

	return cardModel, specificData, dto.Keywords, dto.Metadata, nil
}