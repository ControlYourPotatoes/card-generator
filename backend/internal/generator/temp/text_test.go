package text

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
)

type testProcessor struct {
	title  TitleProcessor
	effect EffectProcessor
	stats  StatsProcessor
}

func newTestProcessor() *testProcessor {
	return &testProcessor{
		title:  NewTitleProcessor(),
		effect: NewEffectProcessor(),
		stats:  NewStatsProcessor(),
	}
}

func TestMain(m *testing.M) {
	testOutputDir := "testoutput"
	os.RemoveAll(testOutputDir)
	os.MkdirAll(testOutputDir, 0755)

	result := m.Run()

	os.RemoveAll(testOutputDir)
	os.Exit(result)
}

func TestProcessors(t *testing.T) {
	tests := []struct {
		name     string
		cardData *card.CardDTO
		wantFile string
	}{
		{
			name: "Basic Creature",
			cardData: &card.CardDTO{
				Type:    card.TypeCreature,
				Name:    "Test Creature",
				Cost:    2,
				Effect:  "This is a test effect.",
				Attack:  2,
				Defense: 2,
				Trait:   "Beast",
			},
			wantFile: "basic_creature.json",
		},
		{
			name: "X Cost Creature",
			cardData: &card.CardDTO{
				Type:    card.TypeCreature,
				Name:    "Variable Beast",
				Cost:    -1,
				Effect:  "This creature enters with X +1/+1 counters.",
				Attack:  0,
				Defense: 0,
				Trait:   "Beast",
			},
			wantFile: "x_cost_creature.json",
		},
		// ... other test cases remain the same
	}

	proc := newTestProcessor()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert CardDTO to Card interface using wrapper
			cardInstance := card.Card(&cardImpl{data: tt.cardData})

			details := &TextDetails{}

			// Process Title
			titleBounds, err := proc.title.ProcessTitle(cardInstance)
			if err != nil {
				t.Fatalf("ProcessTitle() error = %v", err)
			}
			details.Title.Text = cardInstance.GetName()
			details.Title.Position = titleBounds.Rect
			details.Title.Style = titleBounds.Style

			// Process Cost
			costBounds, err := proc.title.ProcessCost(cardInstance.GetCost())
			if err != nil {
				t.Fatalf("ProcessCost() error = %v", err)
			}
			details.Title.Cost = CostInfo{
				Value:    getCostString(cardInstance.GetCost()),
				Position: costBounds.Rect.Min,
				IsXCost:  cardInstance.GetCost() < 0,
				Style:    costBounds.Style,
			}

			// Process Effect
			effectBounds, err := proc.effect.ProcessEffect(cardInstance.GetEffect())
			if err != nil {
				t.Fatalf("ProcessEffect() error = %v", err)
			}
			details.Effect.Keywords = tt.cardData.Keywords
			details.Effect.Text = cardInstance.GetEffect()
			details.Effect.Position = effectBounds.Rect
			details.Effect.Style = effectBounds.Style

			// Process Stats
			statsBounds, err := proc.stats.ProcessStats(cardInstance)
			if err != nil {
				t.Fatalf("ProcessStats() error = %v", err)
			}
			details.Stats.CardType = string(tt.cardData.Type)
			details.Stats.Subtype = tt.cardData.Trait
			// For creature stats, we'll get them from the DTO since cardImpl doesn't expose Attack/Defense directly
			if tt.cardData.Type == card.TypeCreature {
				details.Stats.Power = fmt.Sprintf("%d", tt.cardData.Attack)
				details.Stats.Toughness = fmt.Sprintf("%d", tt.cardData.Defense)
			}
			details.Stats.Position = statsBounds.Rect
			details.Stats.Style = statsBounds.Style

			// Save the processed details
			outputPath := filepath.Join("testoutput", tt.wantFile)
			saveTextDetails(t, outputPath, details)

			// Validate the details
			validateTextDetails(t, details)
		})
	}
}

// Helper functions

func getCostString(cost int) string {
	if cost < 0 {
		return "X"
	}
	return fmt.Sprintf("%d", cost)
}

func validateTextDetails(t *testing.T, details *TextDetails) {
	t.Helper()

	// Title validations
	if details.Title.Text == "" {
		t.Error("Title text should not be empty")
	}
	if details.Title.Position.Empty() {
		t.Error("Title position should be set")
	}

	// Effect validations
	if details.Effect.Text == "" {
		t.Error("Effect text should not be empty")
	}
	if details.Effect.Position.Empty() {
		t.Error("Effect position should be set")
	}

	// Stats validations
	if details.Stats.CardType == "" {
		t.Error("Card type should not be empty")
	}
	if details.Stats.Position.Empty() {
		t.Error("Stats position should be set")
	}

	// Cost validations
	if details.Title.Cost.Value == "" {
		t.Error("Cost value should not be empty")
	}
	if details.Title.Cost.Position.X == 0 && details.Title.Cost.Position.Y == 0 {
		t.Error("Cost position should be set")
	}
}

func saveTextDetails(t *testing.T, path string, details *TextDetails) {
	data, err := json.MarshalIndent(details, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal TextDetails: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("Failed to write test output: %v", err)
	}
}
