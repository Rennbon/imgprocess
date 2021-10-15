package utils

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"io/ioutil"
	"strings"
)

const dpi = 72

func LoadFont(path string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type CssType uint8
type TextAlign string

const (
	CssTypeString CssType = 0
	CssTypeNumber CssType = 1

	TextAlignCenter TextAlign = "center"
	TextAlignLeft   TextAlign = "left"
	TextAlignRight  TextAlign = "right"
)

type GenerateImageParams struct {
	Font          *truetype.Font
	OriginalImage *image.RGBA

	Text          string
	FontSize      float64
	Align         TextAlign
	LetterSpacing int
	Left          int
	Top           int
	Color         *color.RGBA
	Width         int
	Height        int
	Type          CssType
	FontSizeInc   []float64
}

func (p *GenerateImageParams) realFontSize() float64 {
	switch p.Type {
	case CssTypeString:
		return p.FontSize
	case CssTypeNumber:
		textLen := len(p.Text)
		incLen := len(p.FontSizeInc)
		if incLen == 0 {
			return p.FontSize
		}
		if textLen >= incLen {
			return p.FontSize + p.FontSizeInc[incLen-1]
		}
		return p.FontSize + p.FontSizeInc[textLen-2]
	default:
		return 12
	}
}
func (p *GenerateImageParams) realText() []string {
	p.Text = strings.TrimSpace(p.Text)
	if p.Type == CssTypeString {
		return []string{p.Text}
	}
	numbers := make([]string, len(p.Text))
	for k, v := range []rune(p.Text) {
		numbers[k] = string(v)
	}
	return numbers
}
func (p *GenerateImageParams) textWidth() int {
	width, _ := p.fontBounds()
	textLen := len(p.Text)
	spaceLen := 0
	if textLen > 0 {
		spaceLen = textLen - 1
	}
	return width + spaceLen*p.LetterSpacing
}
func (p *GenerateImageParams) textXOffset() int {
	offset := 0
	switch p.Align {
	case TextAlignLeft:
		offset = p.Left
		break
	case TextAlignRight:
		offset = p.Left + p.Width - p.textWidth()
		break
	default:
		offset = p.Left + (p.Width-p.textWidth())/2
		break
	}
	return offset
}

func (p *GenerateImageParams) textYOffset() int {
	_, height := p.fontBounds()
	offset := (p.Height - height) / 2
	return p.Top + offset + height // int(PointToFixed(p.realFontSize())>>6)
}

var (
	ErrGenerateImageParamsIllegal = errors.New("GenerateImageParams illegal")
)

func (p *GenerateImageParams) fontBounds() (width, height int) {
	size := p.realFontSize()
	fontFace := truetype.NewFace(p.Font, &truetype.Options{
		Size: size,
		DPI:  72,
	})
	height = int(fontFace.Metrics().Height-fontFace.Metrics().Descent) >> 6
	for _, v := range []rune(p.Text) {
		adv, ok := fontFace.GlyphAdvance(v)
		if ok {
			width += int(adv) >> 6
		} else {
			width += int(size)
		}
	}
	return width, height
}

func PointToFixed(x float64) fixed.Int26_6 {
	return fixed.Int26_6(x * float64(dpi) * (64.0 / 72.0))
}
func (p *GenerateImageParams) Check() error {
	if p == nil || p.Font == nil || p.OriginalImage == nil || len(p.Text) == 0 {
		return ErrGenerateImageParamsIllegal
	}
	return nil
}
func GenerateImage(p *GenerateImageParams) error {
	err := p.Check()
	if err != nil {
		return err
	}
	color := image.NewUniform(p.Color)
	fontSize := p.realFontSize()
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(p.Font)
	c.SetFontSize(fontSize)
	c.SetClip(p.OriginalImage.Bounds())
	c.SetDst(p.OriginalImage)
	c.SetSrc(color)
	c.SetHinting(font.HintingNone)

	textXOffset := p.textXOffset()
	// top + (height-fontSize)/2
	// Note shift/truncate 6 bits first
	textYOffset := p.textYOffset() //
	pt := freetype.Pt(textXOffset, textYOffset)
	for _, s := range p.realText() {
		pt, err = c.DrawString(s, pt)
		if err != nil {
			return err
		}
		pt.X += fixed.Int26_6(p.LetterSpacing << 6)
	}
	return nil
}
