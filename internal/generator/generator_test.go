package generator

import (
    "errors"
    "image/png"
    "os"
    "path/filepath"
    "testing"
    "time"
    
    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/mocks"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/text"
)

const testOutputDir = "test_output"

// cleanup removes the test output directory
func cleanup() {
    os.RemoveAll(testOutputDir)
}

func TestCardGeneration(t *testing.T) {
    // Create test directory
    if err := os.MkdirAll(testOutputDir, 0755); err != nil {
        t.Fatalf("Failed to create test output directory: %v", err)
    }
    defer cleanup()

    // Create processor instances
    textProc, err := text.NewTextProcessor()
    if err != nil {
        t.Fatalf("Failed to create text processor: %v", err)
    }

    artProc := mocks.NewMockArtProcessor()
    
    // Simulate some network conditions
    artProc.SimulateNetworkError("Failed Card", errors.New("404 not found"))
    artProc.FetchDelay = 100 * time.Millisecond

    // Create generator with processors
    generator, err := NewCardGeneratorWithConfig(&Config{
        TextProc: textProc,
        ArtProc:  artProc,
    })
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }

    tests := []struct {
        name     string
        card     *card.CardData
        wantErr  bool
    }{
        {
            name: "Basic Creature",
            card: &card.CardData{
                Type:    card.TypeCreature,
                Name:    "Mountain Bear",
                Cost:    3,
                Effect:  "When this creature attacks, it gets +2/+0 until end of turn.",
                Attack:  3,
                Defense: 3,
                Trait:   "Beast",
            },
            wantErr: false,
        },
        {
            name: "Basic Spell",
            card: &card.CardData{
                Type:       card.TypeSpell,
                Name:      "Lightning Strike",
                Cost:      2,
                Effect:    "Deal 3 damage to any target.",
                TargetType: "Any",
            },
            wantErr: false,
        },
        {
            name: "Failed Art Fetch",
            card: &card.CardData{
                Type:   card.TypeCreature,
                Name:   "Failed Card",
                Cost:   1,
                Effect: "This card should fail to fetch art.",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            outputPath := filepath.Join(testOutputDir, tt.name+".png")
            
            err := generator.GenerateCard(tt.card, outputPath)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateCard() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr {
                validateOutputFile(t, outputPath)
            }
        })
    }
}

func validateOutputFile(t *testing.T, path string) {
    f, err := os.Open(path)
    if err != nil {
        t.Errorf("Failed to open output file: %v", err)
        return
    }
    defer f.Close()

    img, err := png.Decode(f)
    if err != nil {
        t.Errorf("Failed to decode output image: %v", err)
        return
    }

    bounds := img.Bounds()
    if bounds.Dx() != 1500 || bounds.Dy() != 2100 {
        t.Errorf("Incorrect image dimensions: got %dx%d, want 1500x2100",
            bounds.Dx(), bounds.Dy())
    }
}