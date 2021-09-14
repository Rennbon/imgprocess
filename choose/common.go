package choose

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func MergeImage() {

	img1, _ := os.Open("../resource/p1.png")
	imgb, _ := png.Decode(img1)
	defer img1.Close()

	img2, _ := os.Open("../resource/p2.png")
	imga, _ := png.Decode(img2)
	defer img2.Close()

	offset := image.Pt(0, 0)
	b := imgb.Bounds()
	m := image.NewRGBA(b)
	imgW := b.Size().X
	imgH := b.Size().Y

	newImg := ImageResize(imga, imgW, imgH)

	draw.Draw(m, newImg.Bounds().Add(image.Pt(50, 100)), newImg, image.ZP, draw.Src)
	draw.Draw(m, imgb.Bounds().Add(offset), imgb, image.ZP, draw.Over)
	imgw, _ := os.Create("../resource/p3.png")
	png.Encode(imgw, m)
	defer imgw.Close()
}

// 图片大小调整
func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}

// 图片裁剪 暂时没有用
func ImageCopy(src image.Image, x, y, w, h int) (image.Image, error) {

	var subImg image.Image

	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else {

		return subImg, errors.New("图片解码失败")
	}

	return subImg, nil
}
