package main

import (
	"fmt"
	"log"

	"github.com/idkroff/go-video-text/internal/config"
	"github.com/idkroff/go-video-text/internal/generator/imagegen"
	"github.com/idkroff/go-video-text/internal/generator/videogen"
)

func main() {
	config := config.MustLoad()

	imageGen, err := imagegen.NewGenerator(config.FontPath, config.FontSize, config.MaxWidth)
	if err != nil {
		log.Fatalf("unable to create image generator: %s", err)
	}

	input := "Test 123 bla bla bla."

	videoGenOptions := videogen.VideoGeneratorOptions{
		FPS:         config.VideoOptions.FPS,
		RandomDelay: config.VideoOptions.RandomDelay,
		Delay:       config.VideoOptions.Delay,
		MinDelay:    config.VideoOptions.MinDelay,
		MaxDelay:    config.VideoOptions.MaxDelay,
	}
	videoGen := videogen.NewGenerator(imageGen, videoGenOptions)
	path, err := videoGen.NewStringVideo(input)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	} else {
		print(path)
	}
}
