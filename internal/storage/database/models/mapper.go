package models

import (
    "github.com/ControlYourPotatoes/card-generator/internal/core/card"
)

// ToDomain converts database models to domain models
func (cm *CardModel) ToDomain(specificData TypeSpecificData, keywords []string, metadata map[string]string) (card.Card, error) {
    // Implementation of conversion logic
}

// FromDomain converts domain models to database models
func FromDomain(c card.Card) (*CardModel, *TypeSpecificData, []string, map[string]string, error) {
    // Implementation of conversion logic
}