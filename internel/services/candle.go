package services

import (
	"context"
	"log"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"tg.ai/internel/models"
)

func LoopingCandle(ctx context.Context) {
	log.Println("LoopingCandle starting")
	for {
		time.Sleep(time.Second * 10)
		ps, err := models.ReadPerpetuals(ctx, "")
		if err != nil {
			log.Printf("ReadPerpetuals() => %#v", err)
			time.Sleep(time.Second)
			continue
		}
		for _, p := range ps {
			err = fetchCandle(ctx, p)
			if err != nil {
				log.Printf("fetchCandle(%s) => %#v", p.Symbol, err)
			}
		}
	}
}

func fetchCandle(ctx context.Context, p *models.Perpetual) error {
	latest, err := models.LatestCandleTime(ctx, p.Symbol)
	if err != nil {
		return err
	}
	client := futures.NewClient("", "")
	s := client.NewKlinesService()
	s.Symbol(p.Symbol)
	s.Interval("1d")
	s.StartTime(latest.OpenTime)
	klines, err := s.Do(ctx)
	if err != nil {
		return err
	}
	if len(klines) < 1 {
		return nil
	}
	for _, in := range klines {
		_, err = models.UpsertCandle(ctx, p.Symbol, in)
		if err != nil {
			return err
		}
	}
	return models.DeleteCandles(ctx, p.Symbol)
}
