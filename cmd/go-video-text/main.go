package main

import (
	"log"

	"github.com/idkroff/go-video-text/internal/config"
	"github.com/idkroff/go-video-text/internal/generator/imagegen"
	"github.com/idkroff/go-video-text/internal/generator/videogen"
	"github.com/idkroff/go-video-text/internal/telebot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config := config.MustLoad()

	imageGen, err := imagegen.NewGenerator(config.FontPath, config.FontSize, config.MaxWidth)
	if err != nil {
		log.Fatalf("unable to create image generator: %s", err)
	}

	videoGenOptions := videogen.VideoGeneratorOptions{
		FPS:         config.VideoOptions.FPS,
		RandomDelay: config.VideoOptions.RandomDelay,
		Delay:       config.VideoOptions.Delay,
		MinDelay:    config.VideoOptions.MinDelay,
		MaxDelay:    config.VideoOptions.MaxDelay,
	}
	videoGen := videogen.NewGenerator(imageGen, videoGenOptions)

	tgbot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatalf("unable to set up bot api: %s", err)
	}

	telebot.HandleUpdates(tgbot, videoGen, config.BotStorageChatID)
}
