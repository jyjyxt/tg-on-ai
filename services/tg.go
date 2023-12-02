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
			text = "Ratio HIGH:\n" + models.PerpetualsForHuman(ctx, ps)
		}
		ps, _ = models.ReadDiscretePerpetuals(ctx, "low")
		if len(ps) > 0 {
			if text != "" {
				text = text + "\nRatio LOW:\n"
			}
			text = text + models.PerpetualsForHuman(ctx, ps)
		}
		buy, _ := models.ReadBestPerpetuals(ctx, "buy")
		if len(buy) > 0 {
			if text != "" {
				text = text + "\nBUY:\n"
			}
			text = text + models.PerpetualsForHuman(ctx, buy)
		}
		sell, _ := models.ReadBestPerpetuals(ctx, "sell")
		if len(sell) > 0 {
			if text != "" {
				text = text + "\nSELL:\n"
			}
			text = text + models.PerpetualsForHuman(ctx, sell)
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
