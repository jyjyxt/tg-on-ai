package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg-on-ai/configs"
	"tg-on-ai/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LoopingTGNotify(ctx context.Context, bot *tgbotapi.BotAPI) {
	log.Println("LoopingTGNotify starting")
	for {
		ps, err := models.ReadDiscretePerpetuals(ctx)
		if err != nil {
			log.Printf("ReadDiscretePerpetuals() => %#v", err)
			time.Sleep(time.Second)
			continue
		}

		var texts []string
		for _, p := range ps {
			texts = append(texts, fmt.Sprintf("%s, %s, Price %f, Rate %f, Value %f", p.Symbol, p.Categories, p.MarkPrice, p.LastFundingRate, p.SumOpenInterestValue))
		}
		if len(texts) > 0 {
			msg := tgbotapi.NewMessage(configs.ChannelID, strings.Join(texts, "\n"))
			if _, err := bot.Send(msg); err != nil {
				time.Sleep(time.Second)
				continue
			}
		}

		log.Println("LoopingTGNotify executed at", time.Now())
		time.Sleep(time.Minute * 30)
	}
}
