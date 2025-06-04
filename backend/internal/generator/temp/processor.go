// internal/generator/text/processor.go
package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/ControlYourPotatoes/card-generator/backend/internal/core/card"
	"github.com/ControlYourPotatoes/card-generator/backend/internal/generator/layout"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// TextProcessor coordinates all text rendering operations
type TextProcessor interface {
	RenderText(img draw.Image, data *card.CardDTO, bounds *layout.TextBounds) error
}

type textProcessor struct {
	titleProc  TitleProcessor
	effectProc EffectProcessor
	statsProc  StatsProcessor
	context    *freetype.Context
	font       *truetype.Font
}

func NewTextProcessor() (TextProcessor, error) {
	// Parse the built-in regular font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	return &textProcessor{
		titleProc:  NewTitleProcessor(),
		effectProc: NewEffectProcessor(),
		statsProc:  NewStatsProcessor(),
		font:       font,
	}, nil
}

func (tp *textProcessor) RenderText(img draw.Image, data *card.CardDTO, bounds *layout.TextBounds) error {
	// Initialize freetype context for text rendering
	if err := tp.initContext(img); err != nil {
		return err
	}

	if err := tp.renderTitle(data, bounds); err != nil {
		return fmt.Errorf("failed to render title: %w", err)
	}

	if err := tp.renderEffect(data, bounds); err != nil {
		return fmt.Errorf("failed to render effect: %w", err)
	}

	if data.Type == card.TypeCreature {
		if err := tp.renderStats(data, bounds); err != nil {
			return fmt.Errorf("failed to render stats: %w", err)
		}
	}

	return nil
}

func (tp *textProcessor) initContext(img draw.Image) error {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(tp.font)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.Black))

	tp.context = c
	return nil
}

func (tp *textProcessor) renderTitle(data *card.CardDTO, bounds *layout.TextBounds) error {
	// Create a Card interface implementation from CardDTO
	cardImpl := &cardImpl{data}

	// Process title using the Card interface
	_, err := tp.titleProc.ProcessTitle(cardImpl)
	if err != nil {
		return err
	}

	tp.context.SetFontSize(bounds.Name.FontSize)
	pt := freetype.Pt(
		bounds.Name.Bounds.Min.X,
		bounds.Name.Bounds.Min.Y+int(tp.context.PointToFixed(bounds.Name.FontSize)>>6),
	)

	_, err = tp.context.DrawString(data.Name, pt)
	return err
}

func (tp *textProcessor) renderEffect(data *card.CardDTO, bounds *layout.TextBounds) error {
	_, err := tp.effectProc.ProcessEffect(data.Effect)
	if err != nil {
		return err
	}

	tp.context.SetFontSize(bounds.Effect.FontSize)
	pt := freetype.Pt(
		bounds.Effect.Bounds.Min.X,
		bounds.Effect.Bounds.Min.Y+int(tp.context.PointToFixed(bounds.Effect.FontSize)>>6),
	)

	_, err = tp.context.DrawString(data.Effect, pt)
	return err
}

func (tp *textProcessor) renderStats(data *card.CardDTO, bounds *layout.TextBounds) error {
	if bounds.Stats == nil {
		return fmt.Errorf("stats bounds not provided for creature card")
	}

	// Create a Card interface implementation from CardDTO
	cardImpl := &cardImpl{data}

	_, err := tp.statsProc.ProcessStats(cardImpl)
	if err != nil {
		return err
	}

	// Render power
	tp.context.SetFontSize(bounds.Stats.Left.FontSize)
	ptLeft := freetype.Pt(
		bounds.Stats.Left.Bounds.Min.X,
		bounds.Stats.Left.Bounds.Min.Y+int(tp.context.PointToFixed(bounds.Stats.Left.FontSize)>>6),
	)

	_, err = tp.context.DrawString(fmt.Sprintf("%d", data.Attack), ptLeft)
	if err != nil {
		return err
	}

	// Render toughness
	tp.context.SetFontSize(bounds.Stats.Right.FontSize)
	ptRight := freetype.Pt(
		bounds.Stats.Right.Bounds.Min.X,
		bounds.Stats.Right.Bounds.Min.Y+int(tp.context.PointToFixed(bounds.Stats.Right.FontSize)>>6),
	)

	_, err = tp.context.DrawString(fmt.Sprintf("%d", data.Defense), ptRight)
	return err
}

// cardImpl implements the card.Card interface for CardDTO
type cardImpl struct {
	data *card.CardDTO
}

func (c *cardImpl) GetName() string        { return c.data.Name }
func (c *cardImpl) GetCost() int           { return c.data.Cost }
func (c *cardImpl) GetEffect() string      { return c.data.Effect }
func (c *cardImpl) GetType() card.CardType { return c.data.Type }
func (c *cardImpl) Validate() error        { return nil } // Add proper validation if needed
