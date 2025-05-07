package types

import (
	"fmt"
)

// TagValidator handles validation of tags
type TagValidator struct {
    validCategories map[TagCategory]bool
}

// NewTagValidator creates a new tag validator
func NewTagValidator() *TagValidator {
    return &TagValidator{
        validCategories: map[TagCategory]bool{
            TagTribal:    true,
            TagMechanic:  true,
            TagStrategy:  true,
            TagCost:      true,
            TagSynergy:   true,
            TagCombo:     true,
            TagArchetype: true,
            TagTiming:    true,
        },
    }
}

// ValidateTag checks if a single tag is valid
func (tv *TagValidator) ValidateTag(tag Tag) error {
    if tag.Name == "" {
        return fmt.Errorf("tag name cannot be empty")
    }
    
    if !tv.validCategories[tag.Category] {
        return fmt.Errorf("invalid tag category: %s", tag.Category)
    }
    
    return nil
}

// ValidateTags checks if a slice of tags is valid
func (tv *TagValidator) ValidateTags(tags []Tag) error {
    for _, tag := range tags {
        if err := tv.ValidateTag(tag); err != nil {
            return err
        }
    }
    return nil
}