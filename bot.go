package main

import (
	"context"
	"strings"
	"tg-on-ai/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func startBot(ctx context.Context, bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	filters, err := models.ReadPerpetualSet(ctx, "")
	if err != nil {
		panic(err)
	}

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if imsg := update.ChannelPost; imsg != nil { // If we got a message
			var text string
			cmd := strings.ToLower(strings.TrimSpace(imsg.Text))
			switch cmd {
			case "help":
				cs, err := models.ReadPerpetualCategory(ctx)
				if err != nil {
					return
				}
				text = strings.Join(cs, ", ")
			case "low":
				ps, err := models.ReadLowPerpetuals(ctx)
				if err != nil {
					return
				}
				text = models.PerpetualsForHuman(ctx, ps)
			case "week":
				week, _ := models.ReadWeekStrategies(ctx, models.StrategyNameWeek, 10)
				var symbols []string
				for _, s := range week {
					symbols = append(symbols, s.Symbol)
				}
				ps, _ := models.ReadPerpetualsBySymbols(ctx, symbols)
				text = models.PerpetualsForHuman(ctx, ps)
			case "buy", "sell":
				ps, err := models.ReadBestPerpetuals(ctx, cmd)
				if err != nil {
					return
				}
				text = models.PerpetualsForHuman(ctx, ps)
			case "go":
				text = models.Notify(ctx)
			default:
				ps, err := models.ReadPerpetualsByCategory(ctx, cmd)
				if err != nil {
					return
				}
				text = models.PerpetualsForHuman(ctx, ps)
				if len(ps) == 0 {
					t := strings.ToUpper(cmd)
					if f := filters[t+"USDT"]; f != nil {
						text = models.PerpetualsForHuman(ctx, []*models.Perpetual{f})
					}
					if f := filters["1000"+t+"USDT"]; f != nil {
						text = models.PerpetualsForHuman(ctx, []*models.Perpetual{f})
					}
				}
				if text == "" {
					text = imsg.Text + " 🤟"
				}
			}
			if text == "" {
				continue
			}
			msg := tgbotapi.NewMessage(imsg.Chat.ID, text)
			msg.ReplyToMessageID = imsg.MessageID

			bot.Send(msg)
		}
	}
}
