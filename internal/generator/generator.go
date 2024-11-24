package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/ControlYourPotatoes/card-generator/internal/card"
)

const (
	padding     = 20
	lineHeight  = 190
	indentWidth = 15
	maxWidth    = 2000
)

type Generator struct {
	font        *truetype.Font
	fontSize    float64
	dpi         float64
	templateImg image.Image
}

func NewGenerator() (*Generator, error) {
	// Load the built-in Go regular font
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	// Load the template image
	templatePath := filepath.Join("internal", "generator", "templates", "Sample.png")
	templateFile, err := os.Open(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open template image: %w", err)
	}
	defer templateFile.Close()

	templateImg, err := png.Decode(templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode template image: %w", err)
	}

	return &Generator{
		font:        f,
		fontSize:    84,
		dpi:         64,
		templateImg: templateImg,
	}, nil
}

func (g *Generator) GenerateImage(cardData *card.CardData, outputPath string) error {
	bounds := g.templateImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, g.templateImg, image.Point{}, draw.Over)

	c := freetype.NewContext()
	c.SetDPI(g.dpi)
	c.SetFont(g.font)
	c.SetFontSize(g.fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.White))

	lines := g.calculateLines(cardData)

	y := padding + int(c.PointToFixed(g.fontSize)>>6)
	for _, line := range lines {
		if strings.Contains(line, "\"effect\":") {
			// Handle effect text specially with wrapping
			effectText := strings.TrimPrefix(line, "  \"effect\": \"")
			effectText = strings.TrimSuffix(effectText, "\"")

			// Reduce font size for effect text
			originalSize := g.fontSize
			c.SetFontSize(g.fontSize * 0.9) // Make effect text smaller

			wrappedLines := g.wrapText(effectText, maxWidth)
			for _, wrappedLine := range wrappedLines {
				x := padding + 2*indentWidth // Extra indent for wrapped effect text
				pt := freetype.Pt(x, y)
				_, err := c.DrawString(wrappedLine, pt)
				if err != nil {
					return fmt.Errorf("failed to draw text: %w", err)
				}
				y += lineHeight / 2 // Reduced line height for wrapped text
			}
			//Add a new line after the effect text
			y += lineHeight / 2
			// Restore original font size
			c.SetFontSize(originalSize)
		} else {
			// Handle other lines normally
			indent := strings.Count(line, "  ") * indentWidth
			x := padding + indent
			pt := freetype.Pt(x, y)
			_, err := c.DrawString(line, pt)
			if err != nil {
				return fmt.Errorf("failed to draw text: %w", err)
			}
			y += lineHeight
		}
	}

	// Create output directory and save image...
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	if err := encoder.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}

// wrapText breaks text into lines that fit within maxWidth
func (g *Generator) wrapText(text string, maxWidth int) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return lines
	}

	currentLine := words[0]
	spaceWidth := int(g.fontSize * 0.3) // Approximate space width

	for _, word := range words[1:] {
		// Calculate width of current line + space + new word
		lineWidth := len(currentLine)*int(g.fontSize*0.6) + spaceWidth + len(word)*int(g.fontSize*0.6)

		if lineWidth <= maxWidth {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}
func (g *Generator) calculateLines(data *card.CardData) []string {
	var lines []string
	lines = append(lines, "{")
	lines = append(lines, fmt.Sprintf("  \"type\": \"%s\",", data.Type))
	lines = append(lines, fmt.Sprintf("  \"name\": \"%s\",", data.Name))
	lines = append(lines, fmt.Sprintf("  \"cost\": %d,", data.Cost))
	lines = append(lines, fmt.Sprintf("  \"effect\": \"%s\"", data.Effect))

	// Add type-specific fields
	switch data.Type {
	case card.TypeCreature:
		if data.Attack > 0 || data.Defense > 0 {
			lines[len(lines)-1] += ","
			lines = append(lines, fmt.Sprintf("  \"attack\": %d,", data.Attack))
			lines = append(lines, fmt.Sprintf("  \"defense\": %d", data.Defense))
		}
		if data.Trait != "" {
			lines[len(lines)-1] += ","
			lines = append(lines, fmt.Sprintf("  \"trait\": \"%s\"", data.Trait))
		}
	case card.TypeArtifact:
		if data.IsEquipment {
			lines[len(lines)-1] += ","
			lines = append(lines, "  \"is_equipment\": true")
		}
	case card.TypeSpell:
		if data.TargetType != "" {
			lines[len(lines)-1] += ","
			lines = append(lines, fmt.Sprintf("  \"target_type\": \"%s\"", data.TargetType))
		}
	case card.TypeIncantation:
		if data.Timing != "" {
			lines[len(lines)-1] += ","
			lines = append(lines, fmt.Sprintf("  \"timing\": \"%s\"", data.Timing))
		}
	case card.TypeAnthem:
		lines[len(lines)-1] += ","
		lines = append(lines, "  \"continuous\": true")
	}

	lines = append(lines, "}")
	return lines
}
