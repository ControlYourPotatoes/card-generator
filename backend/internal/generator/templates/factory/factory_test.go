package factory

import (
	"image"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

func TestTemplateFactory(t *testing.T) {
	tests := []struct {
		name      string
		cardType  card.CardType
		wantError bool
	}{
		{"Creature Template", card.TypeCreature, false},
		{"Artifact Template", card.TypeArtifact, false},
		{"Spell Template", card.TypeSpell, false},
		{"Incantation Template", card.TypeIncantation, false},
		{"Anthem Template", card.TypeAnthem, false},
		{"Invalid Type", "InvalidType", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := NewTemplate(tt.cardType)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if template == nil {
				t.Fatal("template is nil")
			}

			// Test GetFrame
			cardData := &card.CardData{
				Type:   tt.cardType,
				Name:   "Test Card",
				Cost:   1,
				Effect: "Test effect",
			}

			frame, err := template.GetFrame(cardData)
			if err != nil {
				t.Errorf("GetFrame() error = %v", err)
			}
			if frame == nil {
				t.Error("frame is nil")
			}

			// Test GetTextBounds
			bounds := template.GetTextBounds(cardData)
			if bounds == nil {
				t.Error("text bounds is nil")
			}

			// Test GetArtBounds
			artBounds := template.GetArtBounds()
			if artBounds == image.ZR {
				t.Error("art bounds is empty")
			}
		})
	}
}
