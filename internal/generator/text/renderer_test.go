package text

import (
    "image"
    "image/draw"
    "image/png"
    "os"
    "path/filepath"
    "testing"
)

func TestRenderer_SingleLine(t *testing.T) {
    cleanupTestOutput(t)

    // Create a test image with white background
    width, height := 500, 200
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

    renderer, err := NewRenderer(img)
    if err != nil {
        t.Fatalf("Failed to create renderer: %v", err)
    }

    // Create test text details
    details := &TextDetails{
        Title: struct {
            Text      string
            Position  image.Rectangle
            Style     TextConfig
            Cost      CostInfo
        }{
            Text:     "Test Title",
            Position: image.Rect(50, 50, 450, 100),
            Style: TextConfig{
                FontSize:  24,
                Bold:     true,
                Color:    "#000000",
                Alignment: "center",
            },
            Cost: CostInfo{
                Value: "3",
                Position: image.Point{X: 400, Y: 50},
                Style: TextConfig{
                    FontSize:  24,
                    Bold:     true,
                    Color:    "#000000",
                    Alignment: "center",
                },
            },
        },
        Effect: struct {
            Keywords  []string
            Text     string
            Position image.Rectangle
            Style    TextConfig
        }{},
        Stats: struct {
            CardType    string
            Subtype     string
            Power       string
            Toughness   string
            Position    image.Rectangle
            Style       TextConfig
        }{},
    }

    renderer.SetTextDetails(details)

    bounds := TextBounds{
        Rect: image.Rect(50, 50, 450, 100),
        Style: TextConfig{
            FontSize:  24,
            Bold:     true,
            Color:    "#000000",
            Alignment: "center",
        },
    }

    if err := renderer.RenderTitle(bounds); err != nil {
        t.Errorf("Failed to render title: %v", err)
    }

    // Also render the cost
    costBounds := TextBounds{
        Rect: image.Rect(400, 50, 450, 100),
        Style: TextConfig{
            FontSize:  24,
            Bold:     true,
            Color:    "#000000",
            Alignment: "center",
        },
    }

    if err := renderer.RenderCost(costBounds); err != nil {
        t.Errorf("Failed to render cost: %v", err)
    }

    // Create output directory if it doesn't exist
    testOutputDir := "test_output"
    if err := os.MkdirAll(testOutputDir, 0755); err != nil {
        t.Fatalf("Failed to create test output directory: %v", err)
    }

    // Save the test image
    outputPath := filepath.Join(testOutputDir, "test_title.png")
    f, err := os.Create(outputPath)
    if err != nil {
        t.Fatalf("Failed to create output file: %v", err)
    }
    defer f.Close()

    if err := png.Encode(f, img); err != nil {
        t.Fatalf("Failed to encode image: %v", err)
    }

    t.Logf("Created test image at: %s", outputPath)
}

func TestRenderer_MultiLine(t *testing.T) {
    cleanupTestOutput(t)

    // Create a test image with white background
    width, height := 500, 300
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

    renderer, err := NewRenderer(img)
    if err != nil {
        t.Fatalf("Failed to create renderer: %v", err)
    }

    // Create test text details
    details := &TextDetails{
        Title: struct {
            Text      string
            Position  image.Rectangle
            Style     TextConfig
            Cost      CostInfo
        }{},
        Effect: struct {
            Keywords  []string
            Text     string
            Position image.Rectangle
            Style    TextConfig
        }{
            Keywords: []string{"HASTE", "CRITICAL"},
            Text:     "This is a longer text that should wrap across multiple lines. Let's see how the text wrapping works.",
            Position: image.Rect(50, 50, 450, 250),
            Style: TextConfig{
                FontSize:  16,
                Bold:     false,
                Color:    "#000000",
                Alignment: "left",
            },
        },
        Stats: struct {
            CardType    string
            Subtype     string
            Power       string
            Toughness   string
            Position    image.Rectangle
            Style       TextConfig
        }{},
    }

    renderer.SetTextDetails(details)

    bounds := TextBounds{
        Rect: image.Rect(50, 50, 450, 250),
        Style: TextConfig{
            FontSize:  16,
            Bold:     false,
            Color:    "#000000",
            Alignment: "left",
        },
    }

    if err := renderer.RenderEffect(bounds); err != nil {
        t.Errorf("Failed to render effect text: %v", err)
    }

    // Create output directory if it doesn't exist
    testOutputDir := "test_output"
    if err := os.MkdirAll(testOutputDir, 0755); err != nil {
        t.Fatalf("Failed to create test output directory: %v", err)
    }

    // Save the test image
    outputPath := filepath.Join(testOutputDir, "test_effect.png")
    f, err := os.Create(outputPath)
    if err != nil {
        t.Fatalf("Failed to create output file: %v", err)
    }
    defer f.Close()

    if err := png.Encode(f, img); err != nil {
        t.Fatalf("Failed to encode image: %v", err)
    }

    t.Logf("Created test image at: %s", outputPath)
}

func TestRenderer_Stats(t *testing.T) {
    cleanupTestOutput(t)

    // Create a test image with white background
    width, height := 500, 200
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

    renderer, err := NewRenderer(img)
    if err != nil {
        t.Fatalf("Failed to create renderer: %v", err)
    }

    // Create test text details
    details := &TextDetails{
        Title: struct {
            Text      string
            Position  image.Rectangle
            Style     TextConfig
            Cost      CostInfo
        }{},
        Effect: struct {
            Keywords  []string
            Text     string
            Position image.Rectangle
            Style    TextConfig
        }{},
        Stats: struct {
            CardType    string
            Subtype     string
            Power       string
            Toughness   string
            Position    image.Rectangle
            Style       TextConfig
        }{
            CardType:  "Creature",
            Subtype:   "Dragon",
            Power:     "5",
            Toughness: "5",
            Position:  image.Rect(50, 50, 450, 100),
            Style: TextConfig{
                FontSize:  24,
                Bold:     true,
                Color:    "#000000",
                Alignment: "center",
            },
        },
    }

    renderer.SetTextDetails(details)

    bounds := TextBounds{
        Rect: image.Rect(50, 50, 450, 100),
        Style: TextConfig{
            FontSize:  24,
            Bold:     true,
            Color:    "#000000",
            Alignment: "center",
        },
    }

    if err := renderer.RenderStats(bounds); err != nil {
        t.Errorf("Failed to render stats: %v", err)
    }

    // Create output directory if it doesn't exist
    testOutputDir := "test_output"
    if err := os.MkdirAll(testOutputDir, 0755); err != nil {
        t.Fatalf("Failed to create test output directory: %v", err)
    }

    // Save the test image
    outputPath := filepath.Join(testOutputDir, "test_stats.png")
    f, err := os.Create(outputPath)
    if err != nil {
        t.Fatalf("Failed to create output file: %v", err)
    }
    defer f.Close()

    if err := png.Encode(f, img); err != nil {
        t.Fatalf("Failed to encode image: %v", err)
    }

    t.Logf("Created test image at: %s", outputPath)
}

// Helper function to clean up test output
func cleanupTestOutput(t *testing.T) {
    testOutputDir := "test_output"
    if err := os.RemoveAll(testOutputDir); err != nil {
        t.Logf("Warning: Failed to clean up test output directory: %v", err)
    }
}