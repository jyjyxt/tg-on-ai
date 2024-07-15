package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"tg.ai/internel/configs"
	"tg.ai/internel/models"
	"tg.ai/internel/session"
)

type Asset struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}

func main() {
	store, err := session.OpenDataSQLite3Store(configs.Path)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx = session.WithSqliteDB(ctx, store)

	for i := 6; i < 9; i++ {
		time.Sleep(time.Second * 3)
		log.Println(i)
		r, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&page=%d", i))
		if err != nil {
			panic(err)
		}

		var body []*Asset
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			panic(err)
		}
		set := make(map[string]string, 0)
		for _, b := range body {
			if set[b.Symbol] == "" {
				set[b.Symbol] = b.ID
			}
		}

		ps, err := models.ReadPerpetuals(ctx, "")
		if err != nil {
			panic(err)
		}

		for _, p := range ps {
			if p.CoinGecko != "" {
				continue
			}
			cgc := set[strings.ToLower(p.BaseAsset)]
			if cgc == "" {
				cgc = set[strings.ReplaceAll(strings.ToLower(p.BaseAsset), "1000", "")]
			}
			err = store.RunInTransaction(ctx, func(ctx context.Context, txn *sql.Tx) error {
				_, err := txn.Exec("UPDATE perpetuals SET coingecko=? WHERE symbol=?", cgc, p.Symbol)
				return err
			})
		}
	}
}
