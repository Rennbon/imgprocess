package choose

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func MergeImage2() {

	img1, _ := os.Open("../resource/Lark.png")
	imga, _ := png.Decode(img1)
	defer img1.Close()

	m := image.NewRGBA(imga.Bounds())
	draw.Draw(m, imga.Bounds(), imga, image.Point{}, draw.Over)
	//addLabel(m, 200, 400, "20")
	for i := 0; i < 1; i++ {
		generateImage(m, "0123456789", 50, 200+i*5, 350*1)
	}
	imgw, _ := os.Create("../resource/Lark2.png")
	png.Encode(imgw, m)
	defer imgw.Close()
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{R: 200, G: 100, B: 0, A: 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	face := &basicfont.Face{
		Advance: 10, //间距
		Width:   100,
		Height:  200,
		Ascent:  11,
		Descent: 2,
		Mask:    img,
		Ranges: []basicfont.Range{
			{'\u0020', '\u007f', 0},
			{'\ufffd', '\ufffe', 95},
		},
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}
func loadFont() (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile("../resource/luximb.ttf")
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func generateImage(rgba *image.RGBA, textContent string, fontSize float64, x, y int) {
	f, _ := loadFont()
	fgColor := image.Black

	code := strings.Replace(textContent, "\t", "    ", -1) // convert tabs into spaces
	text := strings.Split(code, "\n")                      // split newlines into arrays

	fg := image.NewUniform(fgColor)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	textXOffset := x
	textYOffset := y + int(c.PointToFixed(fontSize)>>6) // Note shift/truncate 6 bits first

	pt := freetype.Pt(textXOffset, textYOffset)
	for _, s := range text {
		_, err := c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(12)
	}

}
