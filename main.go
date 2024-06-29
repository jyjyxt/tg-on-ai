package main

import (
	"context"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/unrolled/render"
	"tg.ai/internel/configs"
	"tg.ai/internel/middlewares"
	"tg.ai/internel/routes"
	"tg.ai/internel/services"
	"tg.ai/internel/session"
)

func main() {
	store, err := session.OpenDataSQLite3Store(configs.Path)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = session.WithSqliteDB(ctx, store)

	go services.Root(ctx)
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

	router := httptreemux.New()
	routes.RegisterRoutes(router)
	handler := middlewares.Constraint(router)
	handler = middlewares.Context(handler, store, render.New())
	handler = middlewares.Stats(handler)

	log.Println("localhost:8090")
	http.ListenAndServe("localhost:8090", handler)
}
