package services

import (
	"context"
	"log"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"tg.ai/internel/models"
)

func LoopingExchangeInfo(ctx context.Context) {
	log.Println("LoopingExchangeInfo starting")
	for {
		time.Sleep(time.Minute * 10)
		err := fetchExchangeInfo(ctx)
		if err != nil {
			log.Printf("fetchExchangeInfo() => %#v", err)
			time.Sleep(time.Second)
			continue
		}
		log.Println("LoopingExchangeInfo executed at", time.Now())
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
		if s.QuoteAsset != "USDT" { // only fetch quote usdt
			continue
		}
		if s.Symbol != s.BaseAsset+s.QuoteAsset {
			continue
		}
		if s.Status != "TRADING" {
			continue
		}
		if filter[s.Symbol] != nil {
			delete(filter, s.Symbol)
			continue
		}
		_, err = models.CreatePerpetual(ctx, &s)
		if err != nil {
			return err
		}
	}
	for k := range filter {
		err = models.DeletePerpetual(ctx, k)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoopingPremiumIndex(ctx context.Context) {
	log.Println("LoopingPremiumIndex starting")
	for {
		time.Sleep(time.Minute)
		err := fetchPremiumIndex(ctx)
		if err != nil {
			log.Printf("fetchPremiumIndex() => %#v", err)
			time.Sleep(time.Second)
			continue
		}
		log.Println("LoopingPremiumIndex executed at", time.Now())
	}
}

func fetchPremiumIndex(ctx context.Context) error {
	client := futures.NewClient("", "")
	info, err := client.NewPremiumIndexService().Do(ctx)
	if err != nil {
		return err
	}

	filter, err := models.ReadPerpetualSet(ctx, models.PerpetualSourceBinance)
	if err != nil {
		return err
	}
	for _, in := range info {
		if filter[in.Symbol] == nil {
			continue
		}
		_, err = models.UpdatePerpetual(ctx, in.Symbol, in.MarkPrice, in.LastFundingRate, "", 0)
		if err != nil {
			return err
		}
	}
	return nil
}
