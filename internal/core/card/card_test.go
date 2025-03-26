package card

import (
	"testing"
	"fmt"
)

// TestMain is the main entry point for running tests
// This can be expanded later to setup test dependencies
func TestMain(m *testing.M) {
	// Run the tests
	result := m.Run()
	
	// If tests fail, print a message
	if result != 0 {
		fmt.Println("Some tests failed")
	}
	
	// Exit with the test result code
	// But just return for now since we're not running this directly
}

// TestCardTypes is a basic test to verify that card types are defined correctly
func TestCardTypes(t *testing.T) {
	// Test that card types are defined
	types := []CardType{
		TypeCreature,
		TypeArtifact,
		TypeSpell,
		TypeIncantation,
		TypeAnthem,
	}
	
	expectedTypes := []string{
		"Creature",
		"Artifact",
		"Spell",
		"Incantation",
		"Anthem",
	}
	
	for i, cardType := range types {
		if string(cardType) != expectedTypes[i] {
			t.Errorf("Expected card type %s, got %s", expectedTypes[i], cardType)
		}
	}
}