package services

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/cinar/indicator"
	"tg.ai/internel/models"
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
		_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameMACD, int64(actions[l-1]), 0, 0, asset.Date[l-1].Unix())
		if err != nil {
			return err
		}
	}
	{
		result := indicator.WilliamsR(asset.Low, asset.High, asset.Closing)
		if len(result) == 0 {
			return nil
		}
		l := len(result)
		_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameWilliamsR, 0, result[l-1], 0, asset.Date[l-1].Unix())
		if err != nil {
			return err
		}
	}
	{
		aroonUp, aroonDown := indicator.Aroon(asset.High, asset.Low)
		if len(aroonUp) == 0 {
			return nil
		}
		l := len(aroonUp)
		_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameAroon, 0, aroonUp[l-1], aroonDown[l-1], asset.Date[l-1].Unix())
		if err != nil {
			return err
		}
	}
	{
		_, atr := indicator.Atr(14, asset.High, asset.Low, asset.Closing)
		if len(atr) > 24 {
			l := len(atr)
			now := atr[len(atr)-1]
			r := atr[14:]
			sort.Slice(r, func(i, j int) bool { return r[i] > r[j] })
			max := r[0]
			min := r[len(r)-1]
			_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameATR, 0, (now-min)/(max-min), now, asset.Date[l-1].Unix())
			if err != nil {
				return err
			}
		}
	}
	// week
	{
		offset := 42
		if l := len(asset.Closing); l > offset {
			closings := asset.Closing[l-offset : l-(offset/3)]
			sort.Float64s(closings)
			min := closings[0]
			max := closings[len(closings)-1]
			_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameWeek, 0, p.MarkPrice/max, min/max, asset.Date[l-1].Unix())
		}
	}
	// two week
	{
		offset := 84
		if l := len(asset.Closing); l > offset {
			closings := asset.Closing[l-offset : l-(offset/3)]
			sort.Float64s(closings)
			min := closings[0]
			max := closings[len(closings)-1]
			_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameWeekTwo, 0, p.MarkPrice/max, min/max, asset.Date[l-1].Unix())
		}
	}
	// four week
	{
		offset := 168
		if l := len(asset.Closing); l > offset {
			closings := asset.Closing[l-offset : l-(offset/3)]
			sort.Float64s(closings)
			min := closings[0]
			max := closings[len(closings)-1]
			_, err = models.CreateStrategy(ctx, p.Symbol, models.StrategyNameWeekFour, 0, p.MarkPrice/max, min/max, asset.Date[l-1].Unix())
		}
	}
	return nil
}
