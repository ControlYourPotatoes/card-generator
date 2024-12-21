package validation

import (
    "strings"
    "fmt"

    "github.com/ControlYourPotatoes/card-generator/internal/core/common"

)

// CardValidator defines the interface for card validation
type CardValidator interface {
    ValidateBase() *common.ValidationError  // Updated return type
}

// BaseValidator provides validation for base card properties
type BaseValidator struct {
    Name   string
    Cost   int
    Effect string
}

func (b BaseValidator) ValidateBase() *common.ValidationError {
    if b.Name == "" {
        return common.NewValidationError(common.ErrorTypeRequired, "name cannot be empty", "name")
    }
    if len(b.Name) > 40 {
        return common.NewValidationError(common.ErrorTypeRange, "name exceeds maximum length of 40 characters", "name")
    }
    if b.Cost < -1 {
        return common.NewValidationError(common.ErrorTypeRange, "cost cannot be negative (except -1 for X costs)", "cost")
    }
    if b.Effect == "" {
        return common.NewValidationError(common.ErrorTypeRequired, "effect cannot be empty", "effect")
    }
    return nil
}


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


// Update other validation functions to use common.ValidationError
func ValidateArtifact(isEquipment bool, effect string) *common.ValidationError {
    if isEquipment && !strings.Contains(strings.ToLower(effect), "equip") {
        return common.NewValidationError(common.ErrorTypeInvalid, "equipment artifact must contain equip effect", "effect")
    }
    return nil
}

func ValidateSpell(targetType string) *common.ValidationError {
    if targetType != "" {
        validTargets := map[string]bool{
            "Creature": true,
            "Player":   true,
            "Any":      true,
        }
        if !validTargets[targetType] {
            return common.NewValidationError(common.ErrorTypeInvalid, "invalid target type", "targetType")
        }
    }
    return nil
}

func ValidateIncantation(timing string) *common.ValidationError {
    if timing != "" {
        validTimings := map[string]bool{
            "ON ANY CLASH": true,
            "ON ATTACK":    true,
        }
        if !validTimings[timing] {
            return common.NewValidationError(common.ErrorTypeInvalid, "invalid timing", "timing")
        }
    }
    return nil
}

func ValidateAnthem(continuous bool) *common.ValidationError {
    if !continuous {
        return common.NewValidationError(common.ErrorTypeInvalid, "anthem must be continuous", "continuous")
    }
    return nil
}