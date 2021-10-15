package utils

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func TestGenerateImage(t *testing.T) {

	//font, err := LoadFont("../resource/luximb.ttf")
	font, err := LoadFont("../resource/SourceHanSansCN-Bold.ttf")
	if err != nil {
		t.Fatal(err)
	}
	width := 600
	height := 100
	m := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	draw.Draw(m, m.Bounds(), image.Black, image.Point{}, draw.Over)

	rgba := image.NewRGBA(image.Rect(100, 100, 700, 200))
	draw.Draw(m, rgba.Bounds(), image.White, image.ZP, draw.Over)

	p2 := &GenerateImageParams{
		Font:          font,
		OriginalImage: m,
		Text:          "wjh谁haha18888-+=/%4😂🙃",
		FontSize:      30,
		Align:         TextAlignCenter,
		LetterSpacing: 0,
		Left:          100,
		Top:           100,
		Color: &color.RGBA{
			R: 1, G: 1, B: 1, A: 255,
		},
		Width:       width,
		Height:      height,
		Type:        CssTypeString,
		FontSizeInc: []float64{},
	}
	err = GenerateImage(p2)
	if err != nil {
		t.Fatal(err)
	}
	offset := image.Pt(0, 200)
	draw.Draw(m, rgba.Bounds().Add(offset), image.White, image.ZP, draw.Over)
	p3 := &GenerateImageParams{
		Font:          font,
		OriginalImage: m,
		Text:          "88888",
		FontSize:      100,
		Align:         TextAlignCenter,
		LetterSpacing: 0,
		Left:          100,
		Top:           300,
		Color: &color.RGBA{
			R: 1, G: 1, B: 1, A: 255,
		},
		Width:       width,
		Height:      height,
		Type:        CssTypeNumber,
		FontSizeInc: []float64{-20, -40, -60, -60, -70, -80},
	}
	err = GenerateImage(p3)
	if err != nil {
		t.Fatal(err)
	}
	offset = image.Pt(0, 400)
	draw.Draw(m, rgba.Bounds().Add(offset), image.White, image.ZP, draw.Over)
	p4 := &GenerateImageParams{
		Font:          font,
		OriginalImage: m,
		Text:          "88",
		FontSize:      100, //realFont 80
		Align:         TextAlignCenter,
		LetterSpacing: 0,
		Left:          100,
		Top:           500,
		Color: &color.RGBA{
			R: 1, G: 1, B: 1, A: 255,
		},
		Width:       width,
		Height:      height,
		Type:        CssTypeNumber,
		FontSizeInc: []float64{-20, -40, -60, -60, -70, -80},
	}

	err = GenerateImage(p4)
	if err != nil {
		t.Fatal(err)
	}
	p4.Align = TextAlignLeft
	p4.Top = 700
	offset = image.Pt(0, 600)
	draw.Draw(m, rgba.Bounds().Add(offset), image.White, image.ZP, draw.Over)
	err = GenerateImage(p4)
	if err != nil {
		t.Fatal(err)
	}

	p4.Align = TextAlignRight
	p4.Top = 900
	offset = image.Pt(0, 800)
	draw.Draw(m, rgba.Bounds().Add(offset), image.White, image.ZP, draw.Over)
	err = GenerateImage(p4)
	if err != nil {
		t.Fatal(err)
	}
	imgw, _ := os.Create("../resource/Lark2.png")
	png.Encode(imgw, m)
	defer imgw.Close()
}
