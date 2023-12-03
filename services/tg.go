package services

import (
	"context"
	"log"
	"tg-on-ai/configs"
	"tg-on-ai/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LoopingTGNotify(ctx context.Context, bot *tgbotapi.BotAPI) {
	log.Println("LoopingTGNotify starting")
	for {
		ps, _ := models.ReadDiscretePerpetuals(ctx, "high")
		var text string
		if len(ps) > 0 {
			text = "Ratio HIGH:\n----------\n" + models.PerpetualsForHuman(ctx, ps)
		}
		ps, _ = models.ReadDiscretePerpetuals(ctx, "low")
		if len(ps) > 0 {
			if text != "" {
				text = text + "\n"
			}
			text = text + "Ratio LOW:\n----------\n"
			text = text + models.PerpetualsForHuman(ctx, ps)
		}
		buy, _ := models.ReadBestPerpetuals(ctx, "buy")
		if len(buy) > 0 {
			if text != "" {
				text = text + "\n"
			}
			text = text + "BUY:\n----\n"
			text = text + models.PerpetualsForHuman(ctx, buy)
		}
		sell, _ := models.ReadBestPerpetuals(ctx, "sell")
		if len(sell) > 0 {
			if text != "" {
				text = text + "\n"
			}
			text = text + "SELL:\n-----\n"
			text = text + models.PerpetualsForHuman(ctx, sell)
		}
		pullback, _ := models.ReadPullbackPerpetuals(ctx)
		if len(pullback) > 0 {
			if text != "" {
				text = text + "\n"
			}
			text = text + "Pullback:\n-------\n"
			text = text + models.PerpetualsForHuman(ctx, pullback)
		}
		if text != "" {
			msg := tgbotapi.NewMessage(configs.ChannelID, text)
			if _, err := bot.Send(msg); err != nil {
				time.Sleep(time.Second)
				continue
			}
		}

		log.Println("LoopingTGNotify executed at", time.Now())
		time.Sleep(time.Minute * 30)
	}
}
