// generator/image/generator_test.go

package image

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
)

func TestGenerator(t *testing.T) {
	// Create generator
	gen, err := NewGenerator()
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Create test output directory
	testOutputDir := filepath.Join("testdata", "output")
	if err := os.MkdirAll(testOutputDir, 0755); err != nil {
		t.Fatalf("Failed to create test output directory: %v", err)
	}
	defer os.RemoveAll(testOutputDir)

	tests := []struct {
		name     string
		cardData *card.CardData
	}{
		{
			name: "creature_card",
			cardData: &card.CardData{
				Type:    card.TypeCreature,
				Name:    "Demon Pup",
				Cost:    1,
				Effect:  "Each time you OFFER; gain +1/1.",
				Attack:  1,
				Defense: 1,
				Trait:   "Demon",
			},
		},
		{
			name: "artifact_card",
			cardData: &card.CardData{
				Type:        card.TypeArtifact,
				Name:       "Magic Sword",
				Cost:       3,
				Effect:     "Equip to a creature. It gains +2/+0.",
				IsEquipment: true,
			},
		},
		{
			name: "spell_card",
			cardData: &card.CardData{
				Type:       card.TypeSpell,
				Name:      "Fireball",
				Cost:      2,
				Effect:    "Deal 3 damage to target creature",
				TargetType: "Creature",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(testOutputDir, tt.name+".png")
			
			err := gen.GenerateImage(tt.cardData, outputPath)
			if err != nil {
				t.Errorf("GenerateImage failed: %v", err)
			}

			// Check if file was created
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("Output image was not created: %v", err)
			}
		})
	}
}