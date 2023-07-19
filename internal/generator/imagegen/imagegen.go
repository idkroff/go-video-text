package imagegen

import (
	"image"
	"image/draw"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type ImageGenerator struct {
	Font *truetype.Font
}

func NewGenerator(fontPath string) (*ImageGenerator, error) {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return &ImageGenerator{
		Font: font,
	}, nil
}

func (g *ImageGenerator) NewStringImage(input string, size float64) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 60))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(g.Font)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	for _, s := range input {
		_, err := c.DrawString(string(s), pt)
		if err != nil {
			return nil, err
		}
		pt.X += c.PointToFixed(size * 0.6)
	}

	return img, nil
}
