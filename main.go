package main

import (
	"context"
	"log"
	"net/http"
	"tg-on-ai/configs"
	"tg-on-ai/middlewares"
	"tg-on-ai/routes"
	"tg-on-ai/services"
	"tg-on-ai/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unrolled/render"
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
	go services.LoopingCandle(ctx)
	go services.LoopingStrategy(ctx)
	// token := "6337999999:AAFimM8x_invalidetokenforexample"
	bot, err := tgbotapi.NewBotAPI(configs.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	// 打印 bot 的用户名
	log.Printf("Authorized on account %s", bot.Self.UserName)
	go services.LoopingTGNotify(ctx, bot)
	go startBot(ctx, bot)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	handler := middlewares.Constraint(mux)
	handler = middlewares.Context(handler, store, render.New())
	handler = middlewares.Stats(handler)

	http.ListenAndServe("localhost:8090", handler)
}
