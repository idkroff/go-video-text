package russian_letters_image_test

import (
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/idkroff/go-video-text/internal/generator/imagegen"
)

func Test(t *testing.T) {
	gen, err := imagegen.NewGenerator(
		"../fonts/UbuntuMono-Regular.ttf",
		64,
		500,
	)
	if err != nil {
		log.Fatal(err)
	}

	input := "тестовый русский текст."

	img, err := gen.NewStringImage(input, 1000, 100)
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Create("test.png")
	png.Encode(f, img)
	f.Close()
}
