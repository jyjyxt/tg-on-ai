package services

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"tg.ai/internel/models"
)

type CoinGeckoResponse struct {
	ID        string  `json:"id"`
	Symbol    string  `json:"symbol"`
	MarketCap float64 `json:"market_cap"`
}

func LoopingCoinGecko(ctx context.Context) {
	log.Println("LoopingCoinGecko starting")
	for {
		ids, err := models.SelectAllCoinGeckoIds(ctx)
		if err != nil {
			log.Printf("SelectAllCoinGeckoIds() => %#v", err)
			time.Sleep(time.Second)
			continue
		}
		for i := 0; i < len(ids); i += 100 {
			to := i + 100
			if to > len(ids) {
				to = len(ids)
			}
			err = fetchCoinGecko(ctx, ids[i:to])
			if err != nil {
				log.Printf("fetchCoinGecko() => %#v", err)
				time.Sleep(time.Second)
				continue
			}
			time.Sleep(time.Second * 3)
		}
		log.Println("LoopingCoinGecko executed at", time.Now())
		time.Sleep(time.Minute * 10)
	}
}

func fetchCoinGecko(ctx context.Context, ids []string) error {
	log.Println("fetchCoinGecko() => ", strings.Join(ids, ","))
	resp, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=" + strings.Join(ids, ","))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var coins []CoinGeckoResponse
	if err := json.Unmarshal(body, &coins); err != nil {
		return err
	}

	for _, coin := range coins {
		err = models.UpdatePerpetualMarketCap(ctx, coin.ID, coin.MarketCap)
		if err != nil {
			log.Printf("Failed to update price for %s: %v", coin.ID, err)
		}
	}

	return nil
}
