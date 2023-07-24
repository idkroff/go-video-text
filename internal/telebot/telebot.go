package telebot

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/idkroff/go-video-text/internal/generator/videogen"
)

func HandleUpdates(bot *tgbotapi.BotAPI, videoGen *videogen.VideoGenerator, botStorageChatID int64) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//usersRequests := map[string]

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		ctx, _ := context.WithCancel(context.Background())
		go handleUpdate(ctx, bot, update, videoGen, botStorageChatID)
	}
}

func handleUpdate(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update, videoGen *videogen.VideoGenerator, botStorageChatID int64) {
	if update.Message != nil {
		log.Println(fmt.Sprintf("message from: %d", update.Message.Chat.ID))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use inline menu @text_video_bot to receive video with your text")
		bot.Send(msg)
	}

	if update.InlineQuery != nil {
		fmt.Println("inlinequery")
		fmt.Println(update.InlineQuery.Query)

		videoPath, err := videoGen.NewStringVideo(update.InlineQuery.Query)
		defer os.RemoveAll(videoPath)
		if err != nil {
			log.Printf("error: %s\n", err)
			return
		}

		videoMsg := tgbotapi.NewVideo(
			botStorageChatID,
			tgbotapi.FilePath(videoPath),
		)
		videoMsgSent, err := bot.Send(videoMsg)
		if err != nil {
			log.Printf("error while sending video to itself: %s\n", err)
			return
		}

		answer := tgbotapi.NewInlineQueryResultCachedVideo(uuid.New().String(), videoMsgSent.Video.FileID, "Send video")
		bot.Send(tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			Results:       []interface{}{answer},
			CacheTime:     60,
			IsPersonal:    true,
		})
	}
}
