package telebot

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/idkroff/go-video-text/internal/generator/videogen"
)

func HandleUpdates(bot *tgbotapi.BotAPI, videoGen *videogen.VideoGenerator, botStorageChatID int64) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	usersContexts := map[int64]*context.CancelFunc{}

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go handleUpdate(bot, update, videoGen, botStorageChatID, usersContexts)
	}
}

func handleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	videoGen *videogen.VideoGenerator,
	botStorageChatID int64,
	usersContexts map[int64]*context.CancelFunc,
) {
	if update.Message != nil {
		log.Println(fmt.Sprintf("message from: %d", update.Message.Chat.ID))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use inline menu @text_video_bot to receive video with your text (100 characters max)")
		bot.Send(msg)
	}

	if update.InlineQuery != nil {
		ctx, cancelContext := context.WithCancel(context.Background())

		prevCancel, exists := usersContexts[update.InlineQuery.From.ID]
		if exists && prevCancel != nil {
			(*prevCancel)()
		}
		usersContexts[update.InlineQuery.From.ID] = &cancelContext

		handleInlineQuery(
			ctx,
			videoGen,
			bot,
			update,
			botStorageChatID,
		)
		cancelContext()
	}
}

func handleInlineQuery(ctx context.Context, videoGen *videogen.VideoGenerator, bot *tgbotapi.BotAPI, update tgbotapi.Update, botStorageChatID int64) {
	fmt.Println("inlinequery")
	fmt.Println(update.InlineQuery.Query)

	if len([]rune(update.InlineQuery.Query)) > 100 {
		bot.Send(tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			Results:       []interface{}{},
			CacheTime:     600,
			IsPersonal:    false,
		})
		return
	}

	if ctx.Err() != nil {
		return
	}

	videoGenCtx, cancelVideoGenCtx := context.WithTimeout(ctx, time.Minute)
	defer func() {
		cancelVideoGenCtx()
	}()

	videoPath, err := videoGen.NewStringVideo(videoGenCtx, update.InlineQuery.Query)
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
		log.Printf("error while sending video to storage chat: %s\n", err)
		return
	}

	log.Println(fmt.Sprintf("uploaded video: %s (query: %s)", videoMsgSent.Video.FileID, update.InlineQuery.Query))

	answer := tgbotapi.NewInlineQueryResultCachedVideo(uuid.New().String(), videoMsgSent.Video.FileID, "Send video")
	bot.Send(tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       []interface{}{answer},
		CacheTime:     600,
		IsPersonal:    false,
	})
}
