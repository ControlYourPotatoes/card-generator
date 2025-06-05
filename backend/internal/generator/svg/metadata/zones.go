// svg/metadata/zones.go
package metadata

import (
	"image"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

// ZoneManager handles creation and management of interactive zones for different card types
type ZoneManager struct {
	// TODO: Add configuration in Phase 2
}

// NewZoneManager creates a new zone manager
func NewZoneManager() *ZoneManager {
	return &ZoneManager{}
}

// GetDefaultZones returns the default interactive zones for a given card type
func (zm *ZoneManager) GetDefaultZones(cardType card.CardType) map[string]InteractiveZone {
	// TODO: Implement type-specific zones in Phase 3
	defaultZones := make(map[string]InteractiveZone)
	
	// Common zones for all card types
	defaultZones["card-tap"] = InteractiveZone{
		ID:      "card-tap",
		Bounds:  image.Rect(0, 0, 100, 100), // Placeholder bounds
		Action:  "tap",
		Trigger: "click",
	}
	
	defaultZones["card-inspect"] = InteractiveZone{
		ID:      "card-inspect",
		Bounds:  image.Rect(0, 0, 100, 100), // Placeholder bounds
		Action:  "inspect",
		Trigger: "hover",
	}
	
	return defaultZones
}

// ValidateZone checks if an interactive zone is properly configured
func (zm *ZoneManager) ValidateZone(zone InteractiveZone) error {
	// TODO: Implement validation in Phase 2
	return nil
}

// MergeZones combines multiple zone maps with conflict resolution
func (zm *ZoneManager) MergeZones(zones ...map[string]InteractiveZone) map[string]InteractiveZone {
	// TODO: Implement in Phase 2
	result := make(map[string]InteractiveZone)
	for _, zoneMap := range zones {
		for id, zone := range zoneMap {
			result[id] = zone
		}
	}
	return result
} 