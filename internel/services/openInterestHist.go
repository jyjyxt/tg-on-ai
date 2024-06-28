package services

import (
	"context"
	"log"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"tg.ai/internel/models"
)

func LoopingOpenInterestHist(ctx context.Context) {
	log.Println("LoopingOpenInterestHist starting")
	for {
		ps, err := models.ReadPerpetuals(ctx, "")
		if err != nil {
			log.Printf("ReadPerpetuals() => %#v", err)
			time.Sleep(time.Second)
			continue
		}
		for _, p := range ps {
			err = fetchOpenInterestHist(ctx, p)
			if err != nil {
				log.Printf("fetchOpenInterestHist(%s) => %#v", p.Symbol, err)
			}
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 30)
	}
}

func fetchOpenInterestHist(ctx context.Context, p *models.Perpetual) error {
	t := time.UnixMilli(p.UpdatedAt)
	if t.Add(time.Minute * 30).After(time.Now()) {
		return nil
	}
	p.UpdatedAt = time.Now().Add(time.Minute * -30).UnixMilli()
	client := futures.NewClient("", "")
	s := client.NewOpenInterestStatisticsService()
	s.Symbol(p.Symbol)
	s.Period("30m")
	s.StartTime(p.UpdatedAt)
	info, err := s.Do(ctx)
	if err != nil {
		return err
	}
	if len(info) < 1 {
		return nil
	}
	in := info[len(info)-1]
	_, err = models.UpdatePerpetual(ctx, in.Symbol, "", "", in.SumOpenInterestValue, in.Timestamp)
	return err
}
