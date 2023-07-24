package videogen

import (
	"fmt"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/idkroff/go-video-text/internal/generator/imagegen"
)

type VideoGeneratorOptions struct {
	FPS         int
	RandomDelay bool
	Delay       float64
	MinDelay    float64
	MaxDelay    float64
}

type VideoGenerator struct {
	ImageGen *imagegen.ImageGenerator
	Options  VideoGeneratorOptions
}

func NewGenerator(imageGen *imagegen.ImageGenerator, options VideoGeneratorOptions) *VideoGenerator {
	return &VideoGenerator{
		ImageGen: imageGen,
		Options:  options,
	}
}

func (g *VideoGenerator) GenerateFrames(input string) (string, int, string, error) {
	const op = "internal.generator.videogen.GenerateFrames"

	FPS := g.Options.FPS

	videoID := uuid.New().String()
	framesPath := filepath.Join("tmp", videoID)
	if _, err := os.Stat(framesPath); os.IsNotExist(err) {
		err = os.MkdirAll(framesPath, 0700)
		if err != nil {
			log.Print(op+": unable to mkdir: "+framesPath, err)
		}
	}

	rows := g.ImageGen.GetRows(input)
	w, h := g.ImageGen.CalculateWH(rows)

	img, err := g.ImageGen.NewStringImage("", w, h)
	if err != nil {
		os.RemoveAll(framesPath)
		return "", 0, "", err
	}

	currentFrame := 0
	for i := 0; i < int(g.GetDelay()*float64(FPS)); i++ {
		currentFrame++
		f, err := os.Create(filepath.Join(framesPath, fmt.Sprintf("%d.png", currentFrame)))
		if err != nil {
			os.RemoveAll(framesPath)
			return "", 0, "", err
		}

		png.Encode(f, img)
		f.Close()
	}

	for rowI, row := range rows {
		for rowShift := 0; rowShift < len(row); rowShift++ {
			img, err := g.ImageGen.UpdateStringImage(img, row[:rowShift+1], 0, rowI)
			if err != nil {
				os.RemoveAll(framesPath)
				return "", 0, "", err
			}

			for i := 0; i < int(g.GetDelay()*float64(FPS)); i++ {
				currentFrame++
				f, err := os.Create(filepath.Join(framesPath, fmt.Sprintf("%d.png", currentFrame)))
				if err != nil {
					os.RemoveAll(framesPath)
					return "", 0, "", err
				}

				png.Encode(f, img)
				f.Close()
			}
		}
	}

	return framesPath, currentFrame, videoID, nil
}

func (g *VideoGenerator) NewStringVideo(input string) (string, error) {
	//TODO: add context with timeout to prevent too long generation

	framesPath, _, id, err := g.GenerateFrames(input)
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(framesPath)

	cmd := exec.Command(
		"ffmpeg",
		"-framerate", fmt.Sprintf("%d", g.Options.FPS),
		"-i", filepath.Join(framesPath, "%d.png"),
		filepath.Join("tmp", fmt.Sprintf("%s.mp4", id)),
	)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %s: %w", string(b), err)
	}
	return filepath.Join("tmp", fmt.Sprintf("%s.mp4", id)), nil
}

func (g *VideoGenerator) GetDelay() float64 {
	if g.Options.RandomDelay {
		rand.Seed(time.Now().UnixNano())
		return g.Options.MinDelay + rand.Float64()*(g.Options.MaxDelay-g.Options.MinDelay)
	} else {
		return g.Options.Delay
	}
}
