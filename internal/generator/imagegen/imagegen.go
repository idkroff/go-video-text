package imagegen

import (
	"image"
	"image/draw"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const (
	XSpaceFactor = 0.6
	YSpaceFactor = 1.1
)

type ImageGenerator struct {
	Font     *truetype.Font
	FontSize float64
	MaxWidth int
	HPadding int
	WPadding int
}

func NewGenerator(fontPath string, fontSize float64, maxWidth int) (*ImageGenerator, error) {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return &ImageGenerator{
		Font:     font,
		FontSize: fontSize,
		MaxWidth: maxWidth,
	}, nil
}

func (g *ImageGenerator) NewStringImage(input string, w, h int) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetFont(g.Font)
	c.SetFontSize(g.FontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	pt := freetype.Pt(int(g.FontSize*XSpaceFactor), int(g.FontSize*YSpaceFactor))
	for _, s := range input {
		_, err := c.DrawString(string(s), pt)
		if err != nil {
			return nil, err
		}
		pt.X += c.PointToFixed(g.FontSize * XSpaceFactor)
	}

	return img, nil
}

func (g *ImageGenerator) UpdateStringImage(img *image.RGBA, input string, xStartIndex, yStartIndex int) (*image.RGBA, error) {
	c := freetype.NewContext()
	c.SetFont(g.Font)
	c.SetFontSize(g.FontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	pt := freetype.Pt(
		int(g.FontSize*XSpaceFactor)+xStartIndex*int(g.FontSize*XSpaceFactor),
		int(g.FontSize*YSpaceFactor)+yStartIndex*int(g.FontSize*YSpaceFactor),
	)
	for _, s := range input {
		_, err := c.DrawString(string(s), pt)
		if err != nil {
			return nil, err
		}
		pt.X += c.PointToFixed(g.FontSize * XSpaceFactor)
	}

	return img, nil
}

func (g *ImageGenerator) GetRows(input string) []string {
	maxLettersPerRow := g.MaxWidth / int(g.FontSize*XSpaceFactor)

	rows := []string{""}
	word := ""

	for _, v := range input {
		if string(v) == " " || len(word) >= maxLettersPerRow {
			lastRow := rows[len(rows)-1]
			spaceShift := 0

			if len(lastRow) > 0 {
				spaceShift = 1
			}

			if len(lastRow)+spaceShift+len(word) <= maxLettersPerRow {
				// put on current row
				if spaceShift == 1 {
					rows[len(rows)-1] += " "
				}
				rows[len(rows)-1] += word
			} else {
				// put on next row
				rows = append(rows, word)
			}
			word = ""
		}

		if string(v) != " " {
			word += string(v)
		}
	}
	lastRow := rows[len(rows)-1]
	spaceShift := 0

	if len(lastRow) > 0 {
		spaceShift = 1
	}

	if len(lastRow)+spaceShift+len(word) <= maxLettersPerRow {
		// put on current row
		if spaceShift == 1 {
			rows[len(rows)-1] += " "
		}
		rows[len(rows)-1] += word
	} else {
		// put on next row
		rows = append(rows, word)
	}

	return rows
}

func (g *ImageGenerator) CalculateWH(rows []string) (int, int) {
	if len(rows) == 1 {
		return int(float64(len(rows[0]))*g.FontSize*XSpaceFactor) + int(g.FontSize*XSpaceFactor*2),
			int(g.FontSize * YSpaceFactor * 2)
	}

	return g.MaxWidth + int(g.FontSize*XSpaceFactor)*2,
		(len(rows) + 1) * int(g.FontSize*YSpaceFactor)
}
