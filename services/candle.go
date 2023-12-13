package services

import (
	"context"
	"log"
	"tg-on-ai/models"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

func LoopingCandle(ctx context.Context) {
	log.Println("LoopingCandle starting")
	for {
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
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 30)
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
	s.Interval("4h")
	s.StartTime(latest.OpenTime)
	info, err := s.Do(ctx)
	if err != nil {
		return err
	}
	if len(info) < 1 {
		return nil
	}
	for _, in := range info {
		_, err = models.CreateCandle(ctx, p.Symbol, in.Open, in.High, in.Low, in.Close, in.Volume, in.OpenTime, in.CloseTime)
		if err != nil {
			return err
		}
	}
	return models.DeleteCandles(ctx, p.Symbol)
}
