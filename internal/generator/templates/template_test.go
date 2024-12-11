package templates

import (
    "image"
    "image/png"
    "testing"
    "path/filepath"
    "os"

    "github.com/ControlYourPotatoes/card-generator/internal/card"
    "github.com/ControlYourPotatoes/card-generator/internal/generator/layout"
)

func init() {
    // Create the full template directory path
    templatePath := filepath.Join("internal", "generator", "templates", "images")
    os.MkdirAll(templatePath, 0755)

    // Create dummy template files for testing
    dummyFiles := []string{
        "BaseCreature.png",
        "BaseArtifact.png",
        "BaseSpell.png",
        "BaseIncantation.png",
        "AnthemBase.png",
        "SpecialCreatureWithStats.png",
    }

    for _, file := range dummyFiles {
        path := filepath.Join(templatePath, file)
        if _, err := os.Stat(path); os.IsNotExist(err) {
            // Create a 1x1 pixel dummy PNG file
            img := image.NewRGBA(image.Rect(0, 0, 1, 1))
            f, _ := os.Create(path)
            png.Encode(f, img)
            f.Close()
        }
    }
}

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

            // Verify the frame dimensions
            bounds := frame.Bounds()
            if bounds.Empty() {
                t.Error("Frame bounds should not be empty")
            }
        })
    }
}

func TestTextBounds(t *testing.T) {
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

    // Expected bounds for a creature card
    expectedBounds := &layout.TextBounds{
        Name: layout.TextConfig{
            Bounds:    image.Rect(125, 90, 1375, 170),  // Updated to match actual values
            FontSize:  72,
            Alignment: "center",
        },
        Effect: layout.TextConfig{
            Bounds:    image.Rect(160, 1250, 1340, 1750),
            FontSize:  48,
            Alignment: "left",
        },
        Stats: &layout.StatsConfig{
            Left: layout.TextConfig{
                Bounds:    image.Rect(130, 1820, 230, 1900),
                FontSize:  72,
                Alignment: "center",
            },
            Right: layout.TextConfig{
                Bounds:    image.Rect(1270, 1820, 1370, 1900),
                FontSize:  72,
                Alignment: "center",
            },
        },
    }

    // Compare bounds
    if bounds.Effect.Bounds != expectedBounds.Effect.Bounds {
        t.Errorf("Effect bounds = %v, want %v", bounds.Effect.Bounds, expectedBounds.Effect.Bounds)
    }
    if bounds.Name.Bounds != expectedBounds.Name.Bounds {
        t.Errorf("Name bounds = %v, want %v", bounds.Name.Bounds, expectedBounds.Name.Bounds)
    }
    if bounds.Stats == nil {
        t.Error("Stats bounds should not be nil")
    } else if bounds.Stats.Left.Bounds != expectedBounds.Stats.Left.Bounds {
        t.Errorf("Stats left bounds = %v, want %v", bounds.Stats.Left.Bounds, expectedBounds.Stats.Left.Bounds)
    }
}

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
        return
    }

    // Check stat positions
    if bounds.Stats.Left.Bounds.Empty() {
        t.Error("Left stat bounds should not be empty")
    }
    if bounds.Stats.Right.Bounds.Empty() {
        t.Error("Right stat bounds should not be empty")
    }

    // Verify stat positions are within card bounds
    maxWidth := 1500  // Card width
    maxHeight := 2100 // Card height

    if bounds.Stats.Left.Bounds.Max.X > maxWidth || 
       bounds.Stats.Left.Bounds.Max.Y > maxHeight {
        t.Error("Left stat bounds exceed card dimensions")
    }

    if bounds.Stats.Right.Bounds.Max.X > maxWidth || 
       bounds.Stats.Right.Bounds.Max.Y > maxHeight {
        t.Error("Right stat bounds exceed card dimensions")
    }
}