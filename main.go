package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// token := "6337999999:AAFimM8x_invalidetokenforexample"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	// 打印 bot 的用户名
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// 定义要发送到 channel 的消息
	msg := tgbotapi.NewMessage(-1001933177309, "Hey, Crypto!")

	// 调用 sendMessage 方法发送消息
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
