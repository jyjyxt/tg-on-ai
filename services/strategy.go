package services

import (
	"context"
	"log"
	"tg-on-ai/models"
	"time"

	"github.com/cinar/indicator"
)

func LoopingStrategy(ctx context.Context) {
	log.Println("LoopingStrategy starting")
	for {
		ps, err := models.ReadPerpetuals(ctx, "")
		if err != nil {
			log.Printf("ReadPerpetuals() => %#v", err)
			time.Sleep(time.Second)
			continue
		}
		for _, p := range ps {
			err = fetchStrategy(ctx, p)
			if err != nil {
				log.Printf("fetchStrategy(%s) => %#v", p.Symbol, err)
			}
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 30)
	}
}

func fetchStrategy(ctx context.Context, p *models.Perpetual) error {
	asset, err := models.ReadCandlesAsAsset(ctx, p.Symbol)
	if err != nil || asset == nil {
		return err
	}

	{
		actions := indicator.MacdStrategy(asset)
		if len(actions) == 0 {
			return nil
		}
		l := len(actions)
		_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameMACD, int64(actions[l-1]), 0, asset.Date[l-1].Unix())
		if err != nil {
			return err
		}
	}
	{
		actions := indicator.DefaultKdjStrategy(asset)
		if len(actions) == 0 {
			return nil
		}
		l := len(actions)
		_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameKDJ, int64(actions[l-1]), 0, asset.Date[l-1].Unix())
		if err != nil {
			return err
		}
	}
	return nil
}
