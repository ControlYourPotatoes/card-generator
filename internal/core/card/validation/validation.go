package validation

import (
    "strings"
    "fmt"

    "github.com/ControlYourPotatoes/card-generator/internal/core/common"

)

// CardValidator defines the interface for card validation
type CardValidator interface {
    ValidateBase() *ValidationError
}

// BaseValidator provides validation for base card properties
type BaseValidator struct {
    Name   string
    Cost   int
    Effect string
}

func (b BaseValidator) ValidateBase() *ValidationError {
    if b.Name == "" {
        return NewValidationError(ErrorTypeRequired, "name cannot be empty", "name")
    }
    if len(b.Name) > 40 {
        return NewValidationError(ErrorTypeRange, "name exceeds maximum length of 40 characters", "name")
    }
    if b.Cost < -1 { // -1 allowed for X costs
        return NewValidationError(ErrorTypeRange, "cost cannot be negative (except -1 for X costs)", "cost")
    }
    if b.Effect == "" {
        return NewValidationError(ErrorTypeRequired, "effect cannot be empty", "effect")
    }
    return nil
}

// ValidateCreature validates creature-specific properties
// ValidateCreature validates creature-specific properties
func ValidateCreature(attack, defense int, tribes []common.Tribe) *common.ValidationError {
    // Validate stats
    if attack < 0 {
        return common.NewValidationError(common.ErrorTypeRange, "attack cannot be negative", "attack")
    }
    if defense < 0 {
        return common.NewValidationError(common.ErrorTypeRange, "defense cannot be negative", "defense")
    }

    // Validate tribes
    if len(tribes) == 0 {
        return common.NewValidationError(common.ErrorTypeRequired, "creature must have at least one tribe", "tribes")
    }

    // Check each tribe is valid
    for _, tribe := range tribes {
        if !common.ValidTribes[tribe] {
            return common.NewValidationError(
                common.ErrorTypeInvalid, 
                fmt.Sprintf("invalid tribe: %s", tribe),
                "tribes",
            )
        }
    }

    // Check for duplicate tribes
    seenTribes := make(map[common.Tribe]bool)
    for _, tribe := range tribes {
        if seenTribes[tribe] {
            return common.NewValidationError(
                common.ErrorTypeInvalid,
                fmt.Sprintf("duplicate tribe: %s", tribe),
                "tribes",
            )
        }
        seenTribes[tribe] = true
    }

    return nil
}


// ValidateArtifact validates artifact-specific properties
func ValidateArtifact(isEquipment bool, effect string) *ValidationError {
    if isEquipment && !strings.Contains(strings.ToLower(effect), "equip") {
        return NewValidationError(ErrorTypeInvalid, "equipment artifact must contain equip effect", "effect")
    }
    return nil
}

// ValidateSpell validates spell-specific properties
func ValidateSpell(targetType string) *ValidationError {
    if targetType != "" {
        validTargets := map[string]bool{
            "Creature": true,
            "Player":   true,
            "Any":      true,
        }
        if !validTargets[targetType] {
            return NewValidationError(ErrorTypeInvalid, "invalid target type", "targetType")
        }
    }
    return nil
}

// ValidateIncantation validates incantation-specific properties
func ValidateIncantation(timing string) *ValidationError {
    if timing != "" {
        validTimings := map[string]bool{
            "ON ANY CLASH": true,
            "ON ATTACK":    true,
        }
        if !validTimings[timing] {
            return NewValidationError(ErrorTypeInvalid, "invalid timing", "timing")
        }
    }
    return nil
}

// ValidateAnthem validates anthem-specific properties
func ValidateAnthem(continuous bool) *ValidationError {
    if !continuous {
        return NewValidationError(ErrorTypeInvalid, "anthem must be continuous", "continuous")
    }
    return nil
}