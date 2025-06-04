// internal/generator/generator.go
package generator

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/art"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/templates/factory"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/text"
)

// CardGenerator defines the interface for generating cards
type CardGenerator interface {
	// GenerateCard creates a card image from the provided data and saves it to the specified path
	GenerateCard(data *card.CardDTO, outputPath string) error

	// ValidateCard checks if the card data is valid for generation
	ValidateCard(data *card.CardDTO) error

	// Close cleans up any resources used by the generator
	Close() error
}

// cardGenerator implements the CardGenerator interface
type cardGenerator struct {
	textProc text.TextProcessor
	artProc  art.ArtProcessor
}

// Configuration for the card generator
type Config struct {
	OutputPath string             // Base path for output files
	TextProc   text.TextProcessor // Custom text processor (optional)
	ArtProc    art.ArtProcessor   // Custom art processor (optional)
}

// NewCardGenerator creates a new card generator with default processors
func NewCardGenerator() (CardGenerator, error) {
	textProc, err := text.NewTextProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to create text processor: %w", err)
	}

	artProc := art.NewPlaceholderProcessor()

	return &cardGenerator{
		textProc: textProc,
		artProc:  artProc,
	}, nil
}

// NewCardGeneratorWithConfig creates a new card generator with custom configuration
func NewCardGeneratorWithConfig(cfg *Config) (CardGenerator, error) {
	var err error
	g := &cardGenerator{}

	// Initialize text processor
	if cfg.TextProc != nil {
		g.textProc = cfg.TextProc
	} else {
		g.textProc, err = text.NewTextProcessor()
		if err != nil {
			return nil, fmt.Errorf("failed to create text processor: %w", err)
		}
	}

	// Initialize art processor
	if cfg.ArtProc != nil {
		g.artProc = cfg.ArtProc
	} else {
		g.artProc = art.NewPlaceholderProcessor()
	}

	return g, nil
}

func (g *cardGenerator) ValidateCard(data *card.CardDTO) error {
	if data == nil {
		return fmt.Errorf("card data cannot be nil")
	}

	// Validate required fields
	if data.Name == "" {
		return fmt.Errorf("card name is required")
	}
	if data.Effect == "" {
		return fmt.Errorf("card effect is required")
	}
	if data.Cost < 0 && data.Cost != -1 { // -1 is allowed for X costs
		return fmt.Errorf("invalid card cost: %d", data.Cost)
	}

	// Validate type-specific requirements
	switch data.Type {
	case card.TypeCreature:
		if data.Attack < 0 {
			return fmt.Errorf("creature attack cannot be negative")
		}
		if data.Defense < 0 {
			return fmt.Errorf("creature defense cannot be negative")
		}
	case card.TypeSpell:
		if data.TargetType != "" &&
			data.TargetType != "Creature" &&
			data.TargetType != "Player" &&
			data.TargetType != "Any" {
			return fmt.Errorf("invalid target type: %s", data.TargetType)
		}
	}

	return nil
}

func (g *cardGenerator) GenerateCard(data *card.CardDTO, outputPath string) error {
	// Validate card data
	if err := g.ValidateCard(data); err != nil {
		return fmt.Errorf("invalid card data: %w", err)
	}

	// Get appropriate template for card type
	template, err := factory.NewTemplate(data.Type)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Get base frame
	frame, err := template.GetFrame(data)
	if err != nil {
		return fmt.Errorf("failed to get frame: %w", err)
	}

	// Create base image
	bounds := frame.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, frame, image.Point{}, draw.Over)

	// Process and add art first (so text can overlay if needed)
	artBounds := template.GetArtBounds()
	art, err := g.artProc.ProcessArt(data, artBounds)
	if err != nil {
		return fmt.Errorf("failed to process art: %w", err)
	}
	draw.Draw(img, artBounds, art, image.Point{}, draw.Over)

	// Process and add text
	textBounds := template.GetTextBounds(data)
	if err := g.textProc.RenderText(img, data, textBounds); err != nil {
		return fmt.Errorf("failed to render text: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save the image
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// Use png encoder with best compression
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	if err := encoder.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

func (g *cardGenerator) Close() error {
	// Clean up any resources if needed
	return nil
}
