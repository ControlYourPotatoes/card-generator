package templates

import (
	"image"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
)

func TestTemplateCreation(t *testing.T) {
	tests := []struct {
		name     string
		cardType card.CardType
	}{
		{"Creature", card.TypeCreature},
		{"Artifact", card.TypeArtifact},
		{"Spell", card.TypeSpell},
		{"Incantation", card.TypeIncantation},
		{"Anthem", card.TypeAnthem},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := NewTemplate(tt.cardType)
			if err != nil {
				t.Fatalf("Failed to create template: %v", err)
			}
			if template == nil {
				t.Error("Template should not be nil")
			}
		})
	}
}

func TestFrameLoading(t *testing.T) {
	tests := []struct {
		name     string
		cardType card.CardType
		cardData *card.CardData
		wantFile string
	}{
		{
			name:     "Basic Creature",
			cardType: card.TypeCreature,
			cardData: &card.CardData{
				Type:    card.TypeCreature,
				Name:    "Test Creature",
				Cost:    1,
				Effect:  "Test effect",
				Attack:  1,
				Defense: 1,
			},
			wantFile: "BaseCreature.png",
		},
		{
			name:     "Basic Artifact",
			cardType: card.TypeArtifact,
			cardData: &card.CardData{
				Type:   card.TypeArtifact,
				Name:   "Test Artifact",
				Cost:   1,
				Effect: "Test effect",
			},
			wantFile: "BaseArtifact.png",
		},
		// Add more test cases for other card types
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := NewTemplate(tt.cardType)
			if err != nil {
				t.Fatalf("Failed to create template: %v", err)
			}

			frame, err := template.GetFrame(tt.cardData)
			if err != nil {
				t.Fatalf("Failed to get frame: %v", err)
			}
			if frame == nil {
				t.Error("Frame should not be nil")
			}
			// Could add image dimension checks here
		})
	}
}

func TestTextBounds(t *testing.T) {
	// Create a creature template for testing bounds
	template, err := NewTemplate(card.TypeCreature)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	cardData := &card.CardData{
		Type:    card.TypeCreature,
		Name:    "Test Creature",
		Cost:    1,
		Effect:  "Test effect",
		Attack:  1,
		Defense: 1,
	}

	bounds := template.GetTextBounds(cardData)
	if bounds == nil {
		t.Fatal("Bounds should not be nil")
	}

	// Test specific bounds
	tests := []struct {
		name       string
		got        image.Rectangle
		wantWidth  int
		wantHeight int
	}{
		{"Name bounds", bounds.Name.Bounds, 1250, 80},
		{"Effect bounds", bounds.Effect.Bounds, 1180, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if width := tt.got.Dx(); width != tt.wantWidth {
				t.Errorf("Width = %v, want %v", width, tt.wantWidth)
			}
			if height := tt.got.Dy(); height != tt.wantHeight {
				t.Errorf("Height = %v, want %v", height, tt.wantHeight)
			}
		})
	}
}

// TestCreatureStats ensures creature templates include stat positioning
func TestCreatureStats(t *testing.T) {
	template, err := NewTemplate(card.TypeCreature)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	cardData := &card.CardData{
		Type:    card.TypeCreature,
		Name:    "Test Creature",
		Attack:  1,
		Defense: 1,
	}

	bounds := template.GetTextBounds(cardData)
	if bounds.Stats == nil {
		t.Error("Creature template should include stat bounds")
	}

	// Test stat positions
	if bounds.Stats != nil {
		if bounds.Stats.Left.Bounds.Empty() {
			t.Error("Left stat bounds should not be empty")
		}
		if bounds.Stats.Right.Bounds.Empty() {
			t.Error("Right stat bounds should not be empty")
		}
	}
}
