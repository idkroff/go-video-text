package videogen

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/idkroff/go-video-text/internal/generator/imagegen"
	"github.com/idkroff/go-video-text/lib/clone"
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

	wg := sync.WaitGroup{}

	for rowI, row := range rows {
		for rowShift := 0; rowShift < len(row); rowShift++ {
			wg.Add(1)
			framesCount := int(g.GetDelay() * float64(FPS))

			img, err := g.ImageGen.UpdateStringImage(img, row[:rowShift+1], 0, rowI)
			if err != nil {
				os.RemoveAll(framesPath)

				//TODO: add error handling?
				return "", 0, "", err
			}

			go func(rowI int, row string, rowShift int, img *image.RGBA, framesStart int, framesCount int) {
				defer wg.Done()

				for i := 0; i < framesCount; i++ {
					currentFrameLocal := framesStart + i

					log.Println(fmt.Sprintf("frame %d %s", currentFrameLocal, row[:rowShift+1]))
					f, err := os.Create(filepath.Join(framesPath, fmt.Sprintf("%d.png", currentFrameLocal)))
					if err != nil {
						os.RemoveAll(framesPath)

						//TODO: add error handling?
						return
					}

					png.Encode(f, img)
					f.Close()
				}
			}(rowI, row, rowShift, clone.CloneImageAsRGBA(img), currentFrame, framesCount)

			currentFrame += framesCount
		}
	}

	wg.Wait()

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
