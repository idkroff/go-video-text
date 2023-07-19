package main

import (
	"image/png"
	"log"
	"os"

	"github.com/idkroff/go-video-text/internal/config"
	"github.com/idkroff/go-video-text/internal/generator/imagegen"
)

func main() {
	config := config.MustLoad()

	imageGen, err := imagegen.NewGenerator(config.FontPath)
	if err != nil {
		log.Fatalf("unable to create image generator: %s", err)
	}

	img, err := imageGen.NewStringImage("test123", 36)
	if err != nil {
		log.Fatalf("unable to generate string image: %s", err)
	}

	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(
		f,
		img,
	)
}
