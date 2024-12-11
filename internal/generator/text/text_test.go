// text_test.go
package text

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
)

func TestMain(m *testing.M) {
    // Setup test output directory
    testOutputDir := "testoutput"
    os.RemoveAll(testOutputDir) // Clean previous test output
    os.MkdirAll(testOutputDir, 0755)
    
    // Run tests
    result := m.Run()
    
    // Cleanup
    os.RemoveAll(testOutputDir)
    os.Exit(result)
}

func TestCardProcessing(t *testing.T) {
    tests := []struct {
        name     string
        card     card.Card
        wantFile string
    }{
        {
            name: "Basic Creature",
            card: &card.Creature{
                BaseCard: card.BaseCard{
                    Name:   "Test Creature",
                    Cost:   2,
                    Effect: "This is a test effect",
                    Type:   card.TypeCreature,
                },
                Attack:  2,
                Defense: 2,
                Trait:   "Beast",
            },
            wantFile: "basic_creature.json",
        },
        {
            name: "X Cost Creature",
            card: &card.Creature{
                BaseCard: card.BaseCard{
                    Name:   "Variable Beast",
                    Cost:   -1, // Indicates X cost
                    Effect: "This creature enters with X +1/+1 counters",
                    Type:   card.TypeCreature,
                },
                Attack:  0,
                Defense: 0,
                Trait:   "Beast",
            },
            wantFile: "x_cost_creature.json",
        },
        {
            name: "Keyword Creature",
            card: &card.Creature{
                BaseCard: card.BaseCard{
                    Name:   "Elite Warrior",
                    Cost:   3,
                    Effect: "HASTE, CRITICAL • Deal 2 damage to any target",
                    Type:   card.TypeCreature,
                },
                Attack:  3,
                Defense: 2,
                Trait:   "Warrior",
            },
            wantFile: "keyword_creature.json",
        },
        // Add more test cases for other card types...
    }

    processor := NewTextProcessor() // You'll need to implement this
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            details, err := processor.ProcessCard(tt.card)
            if err != nil {
                t.Fatalf("ProcessCard() error = %v", err)
            }

            // Save the processed details to a JSON file
            outputPath := filepath.Join("testoutput", tt.wantFile)
            saveTextDetails(t, outputPath, details)

            // Compare with expected output
            // You might want to create golden files for comparison
            // or add specific assertions about the TextDetails
            validateTextDetails(t, details)
        })
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

func validateTextDetails(t *testing.T, details *TextDetails) {
    // Add validation logic here
    if details.Title.Text == "" {
        t.Error("Title text should not be empty")
    }
    if details.Title.Position.Empty() {
        t.Error("Title position should be set")
    }
    // Add more validation as needed...
}