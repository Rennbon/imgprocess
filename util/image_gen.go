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
	CssTypeString CssType = 0 // 文本类型 NOTE:字符串类型超过width长度会自动切断
	CssTypeNumber CssType = 1 // 数值类型

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
func (p *GenerateImageParams) lettersSpacing() int {
	textLen := len(p.Text)
	spaceLen := 0
	if textLen > 0 {
		spaceLen = textLen - 1
	}
	return spaceLen * p.LetterSpacing
}



func (p *GenerateImageParams) textXYOffset() (x, y int) {
	width, height := p.fontBounds()
	margin := (p.Height - height) / 2
	y = p.Top + margin + height

	switch p.Align {
	case TextAlignLeft:
		x = p.Left
		break
	case TextAlignRight:
		x = p.Left + p.Width - width + p.lettersSpacing()
		break
	default:
		x = p.Left + (p.Width-width-p.lettersSpacing())/2
		break
	}
	return x, y
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
	tmp := width
	height = int(fontFace.Metrics().Height-fontFace.Metrics().Descent) >> 6
	for i, v := range []rune(p.Text) {
		adv, ok := fontFace.GlyphAdvance(v)
		if ok {
			tmp += int(adv) >> 6
		} else {
			tmp += int(size)
		}
		if p.Type == CssTypeString && p.Width < tmp+i*p.LetterSpacing {
			p.Text = string([]rune(p.Text)[:i])
			break
		}
		width = tmp
	}
	return width, height
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

	textXOffset, textYOffset := p.textXYOffset()
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
