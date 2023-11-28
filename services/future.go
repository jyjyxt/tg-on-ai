package services

import (
	"context"
	"log"
	"tg-on-ai/models"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func LoopingExchangeInfo(ctx context.Context) {
	log.Println("LoopingExchangeInfo starting")
	for {
		err := fetchExchangeInfo(ctx)
		if err != nil {
			log.Printf("fetchExchangeInfo() => %#v", err)
			time.Sleep(time.Second)
			continue
		}
		log.Println("LoopingExchangeInfo executed at", time.Now())
		time.Sleep(time.Second * 30)
	}
}

func fetchExchangeInfo(ctx context.Context) error {
	client := futures.NewClient("", "")
	info, err := client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return err
	}

	filter, err := models.ReadPerpetualSet(ctx, models.PerpetualSourceBinance)
	if err != nil {
		return err
	}
	for _, s := range info.Symbols {
		if filter[s.Symbol] != nil {
			continue
		}
		if s.QuoteAsset != "USDT" { // only fetch quote usdt
			continue
		}
		_, err = models.CreatePerpetual(ctx, s.Symbol, s.BaseAsset, s.QuoteAsset, models.PerpetualSourceBinance, s.UnderlyingSubType)
		if err != nil {
			return err
		}
	}
	return nil
}
