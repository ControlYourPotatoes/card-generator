package parser

import (
	"os"
	"path/filepath"
	"testing"

	"../internal/card"
)

func TestParse(t *testing.T) {
	// Open test file
	testFile, err := os.Open(filepath.Join("../../test/testdata/test_cards.csv"))
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer testFile.Close()

	// Create new parser
	p := NewParser(testFile)

	// Parse the CSV
	cards, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	// Test cases for expected cards from test_cards.csv
	tests := []struct {
		name     string
		cardType card.CardType
		count    int
	}{
		{
			name:     "Creature cards",
			cardType: card.TypeCreature,
			count:    5, // Based on test_cards.csv content
		},
		{
			name:     "Artifact cards",
			cardType: card.TypeArtifact,
			count:    3,
		},
		{
			name:     "Spell cards",
			cardType: card.TypeSpell,
			count:    3,
		},
		{
			name:     "Incantation cards",
			cardType: card.TypeIncantation,
			count:    3,
		},
		{
			name:     "Anthem cards",
			cardType: card.TypeAnthem,
			count:    3,
		},
	}

	// Count cards by type
	cardCounts := make(map[card.CardType]int)
	for _, c := range cards {
		cardCounts[c.GetType()]++
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := cardCounts[tt.cardType]
			if count != tt.count {
				t.Errorf("got %d %s cards, want %d", count, tt.cardType, tt.count)
			}
		})
	}
}

func TestParseCreatureDetails(t *testing.T) {
	// Test specific creature card parsing
	testData := `Type,Name,Cost,Effect,Attack,Defense
Creature,Demon Pup,1,Each time you OFFER; gain +1/1.,1,1`

	// Create a temporary file with test data
	tmpFile, err := os.CreateTemp("", "test_creature.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Reopen the file for reading
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	// Parse the CSV
	p := NewParser(file)
	cards, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	if len(cards) != 1 {
		t.Fatalf("Expected 1 card, got %d", len(cards))
	}

	// Type assert and check creature details
	creature, ok := cards[0].(*card.Creature)
	if !ok {
		t.Fatal("First card is not a Creature")
	}

	// Verify creature fields
	tests := []struct {
		name     string
		got      interface{}
		want     interface{}
	}{
		{"Name", creature.Name, "Demon Pup"},
		{"Cost", creature.Cost, 1},
		{"Effect", creature.Effect, "Each time you OFFER; gain +1/1."},
		{"Attack", creature.Attack, 1},
		{"Defense", creature.Defense, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name    string
		csv     string
		wantErr string
	}{
		{
			name: "invalid cost",
			csv: `Type,Name,Cost,Effect,Attack,Defense
Creature,Test Card,invalid,Test effect,1,1`,
			wantErr: "invalid cost",
		},
		{
			name: "invalid attack",
			csv: `Type,Name,Cost,Effect,Attack,Defense
Creature,Test Card,1,Test effect,invalid,1`,
			wantErr: "invalid attack",
		},
		{
			name: "invalid defense",
			csv: `Type,Name,Cost,Effect,Attack,Defense
Creature,Test Card,1,Test effect,1,invalid`,
			wantErr: "invalid defense",
		},
		{
			name: "missing required column",
			csv: `Type,Name,Effect,Attack,Defense
Creature,Test Card,Test effect,1,1`,
			wantErr: "missing required column",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with test data
			tmpFile, err := os.CreateTemp("", "test_errors.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			
			if _, err := tmpFile.WriteString(tt.csv); err != nil {
				t.Fatalf("Failed to write test data: %v", err)
			}
			
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			// Reopen the file for reading
			file, err := os.Open(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer file.Close()

			p := NewParser(file)
			_, err = p.Parse()
			if err == nil {
				t.Errorf("Expected error containing %q, got nil", tt.wantErr)
			} else if !contains(err.Error(), tt.wantErr) {
				t.Errorf("Expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[:len(substr)] == substr
}