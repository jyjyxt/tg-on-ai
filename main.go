package main

import (
	"context"
	"log"
	"strings"
	"tg-on-ai/configs"
	"tg-on-ai/models"
	"tg-on-ai/services"
	"tg-on-ai/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	store, err := session.OpenDataSQLite3Store(configs.Path)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = session.WithSqliteDB(ctx, store)

	go services.LoopingExchangeInfo(ctx)
	go services.LoopingPremiumIndex(ctx)
	go services.LoopingOpenInterestHist(ctx)
	// token := "6337999999:AAFimM8x_invalidetokenforexample"
	bot, err := tgbotapi.NewBotAPI(configs.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	// 打印 bot 的用户名
	log.Printf("Authorized on account %s", bot.Self.UserName)
	go services.LoopingTGNotify(ctx, bot)

	/*
		// 定义要发送到 channel 的消息
		msg := tgbotapi.NewMessage(-1001933177309, "Hey, Crypto!")

		// 调用 sendMessage 方法发送消息
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	*/

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

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
				text = models.PerpetualsForHuman(ps)
			default:
				ps, err := models.ReadPerpetualsByCategory(ctx, cmd)
				if err != nil {
					return
				}
				text = models.PerpetualsForHuman(ps)
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
