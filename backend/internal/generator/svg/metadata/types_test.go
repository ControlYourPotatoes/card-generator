// svg/metadata/types_test.go
package metadata

import (
	"encoding/json"
	"image"
	"testing"
)

func TestSVGMetadataJSONSerialization(t *testing.T) {
	metadata := SVGMetadata{
		CardID:           "test-card-001",
		InteractiveZones: []InteractiveZone{},
		AnimationTargets: []string{"element1", "element2"},
		GameState:        map[string]string{"tapped": "false"},
		Version:          "1.0.0",
		GeneratedAt:      "2024-01-01T00:00:00Z",
	}

	// Test JSON marshaling
	data, err := json.Marshal(metadata)
	if err != nil {
		t.Fatalf("Failed to marshal SVGMetadata: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled SVGMetadata
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SVGMetadata: %v", err)
	}

	// Verify round-trip
	if unmarshaled.CardID != metadata.CardID {
		t.Errorf("CardID mismatch: got %s, want %s", unmarshaled.CardID, metadata.CardID)
	}
	if unmarshaled.Version != metadata.Version {
		t.Errorf("Version mismatch: got %s, want %s", unmarshaled.Version, metadata.Version)
	}
}

func TestInteractiveZoneValidation(t *testing.T) {
	zone := InteractiveZone{
		ID:      "test-zone",
		Bounds:  image.Rect(0, 0, 100, 100),
		Action:  "tap",
		Trigger: "click",
		Data:    map[string]interface{}{"cost": 3},
	}

	// Test JSON serialization
	data, err := json.Marshal(zone)
	if err != nil {
		t.Fatalf("Failed to marshal InteractiveZone: %v", err)
	}

	var unmarshaled InteractiveZone
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal InteractiveZone: %v", err)
	}

	if unmarshaled.ID != zone.ID {
		t.Errorf("ID mismatch: got %s, want %s", unmarshaled.ID, zone.ID)
	}
	if unmarshaled.Action != zone.Action {
		t.Errorf("Action mismatch: got %s, want %s", unmarshaled.Action, zone.Action)
	}
}

func TestSVGBoundsConversion(t *testing.T) {
	bounds := SVGBounds{
		X:      10.5,
		Y:      20.5,
		Width:  100.0,
		Height: 200.0,
	}

	// Test JSON serialization
	data, err := json.Marshal(bounds)
	if err != nil {
		t.Fatalf("Failed to marshal SVGBounds: %v", err)
	}

	var unmarshaled SVGBounds
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal SVGBounds: %v", err)
	}

	if unmarshaled.X != bounds.X {
		t.Errorf("X mismatch: got %f, want %f", unmarshaled.X, bounds.X)
	}
	if unmarshaled.Width != bounds.Width {
		t.Errorf("Width mismatch: got %f, want %f", unmarshaled.Width, bounds.Width)
	}
} 