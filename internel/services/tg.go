package services

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg.ai/internel/configs"
	"tg.ai/internel/models"
)

func LoopingTGNotify(ctx context.Context, bot *tgbotapi.BotAPI) {
	log.Println("LoopingTGNotify starting")
	for {
		text := models.Notify(ctx)
		if text != "" {
			msg := tgbotapi.NewMessage(configs.ChannelID, text)
			if _, err := bot.Send(msg); err != nil {
				time.Sleep(time.Second)
				continue
			}
		}

		time.Sleep(time.Minute * 30)
		log.Println("LoopingTGNotify executed at", time.Now())
	}
}
