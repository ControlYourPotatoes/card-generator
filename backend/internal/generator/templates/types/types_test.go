package types

import (
	"path/filepath"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
)

func TestCreatureTemplate(t *testing.T) {
	tests := []struct {
		name          string
		cardData      *card.CardData
		expectedFrame string
	}{
		{
			name: "Basic Creature",
			cardData: &card.CardData{
				Type:    card.TypeCreature,
				Name:    "Test Creature",
				Cost:    2,
				Effect:  "Basic effect",
				Attack:  2,
				Defense: 2,
			},
			expectedFrame: "BaseCreature.png",
		},
		{
			name: "Legendary Creature",
			cardData: &card.CardData{
				Type:    card.TypeCreature,
				Name:    "Legendary Creature",
				Cost:    2,
				Effect:  "Legendary effect",
				Attack:  2,
				Defense: 2,
				Trait:   "Legendary",
			},
			expectedFrame: "SpecialCreatureWithStats.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := NewCreatureTemplate()
			if err != nil {
				t.Fatalf("failed to create template: %v", err)
			}

			// Test frame selection
			frame, err := template.GetFrame(tt.cardData)
			if err != nil {
				t.Errorf("GetFrame() error = %v", err)
			}
			if frame == nil {
				t.Error("frame should not be nil")
			}

			// Test bounds calculation
			bounds := template.GetTextBounds(tt.cardData)
			if bounds == nil {
				t.Error("text bounds should not be nil")
			}
			if bounds.Stats == nil {
				t.Error("creature template should have stats bounds")
			}
		})
	}
}

func TestSpellTemplate(t *testing.T) {
	tests := []struct {
		name          string
		cardData      *card.CardData
		expectedFrame string
	}{
		{
			name: "Basic Spell",
			cardData: &card.CardData{
				Type:       card.TypeSpell,
				Name:       "Test Spell",
				Cost:       2,
				Effect:     "Deal 3 damage",
				TargetType: "Creature",
			},
			expectedFrame: "BaseSpell.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := NewSpellTemplate()
			if err != nil {
				t.Fatalf("failed to create template: %v", err)
			}

			// Test frame selection
			frame, err := template.GetFrame(tt.cardData)
			if err != nil {
				t.Errorf("GetFrame() error = %v", err)
			}
			if frame == nil {
				t.Error("frame should not be nil")
			}

			// Test bounds calculation
			bounds := template.GetTextBounds(tt.cardData)
			if bounds == nil {
				t.Error("text bounds should not be nil")
			}
			// Spells don't have stats, so no need to check for stats bounds
		})
	}
}

// Helper function to check the template directory structure
func TestTemplateDirectory(t *testing.T) {
	requiredFiles := []string{
		"BaseCreature.png",
		"SpecialCreature.png",
		"BaseSpell.png",
		"SpecialCreatureWithStats.png",
		"BaseArtifact.png",
		"BaseAnthem.png",
		"BaseIncantation.png",
	}

	templateDir := filepath.Join("images") // Adjust path as needed
	for _, file := range requiredFiles {
		t.Logf("Checking for template file: %s", file)
		fullPath := filepath.Join(templateDir, file)
		t.Logf("Full path: %s", fullPath)
	}
}
