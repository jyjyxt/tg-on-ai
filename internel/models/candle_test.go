package models

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/cinar/indicator"
	"github.com/stretchr/testify/require"
)

func testCandle(t *testing.T) {
	require := require.New(t)
	ctx := setup()
	defer os.Remove(pathTest)

	symbol := "ASTRUSDT"
	candle, err := LatestCandleTime(ctx, symbol)
	require.Nil(err)
	require.NotNil(candle)

	client := futures.NewClient("", "")
	s := client.NewKlinesService()
	s.Symbol(symbol)
	s.Interval("4h")
	log.Println(time.UnixMilli(candle.OpenTime))
	s.StartTime(candle.OpenTime)
	info, err := s.Do(ctx)
	require.Nil(err)
	for _, in := range info {
		candle, err = UpsertCandle(ctx, symbol, in)
		require.Nil(err)
		require.NotNil(candle)
	}

	cs, err := ReadCandles(ctx, symbol)
	require.Nil(err)
	require.Len(cs, 180)
	err = DeleteCandles(ctx, symbol)
	require.Nil(err)
	cs, err = ReadCandles(ctx, symbol)
	require.Nil(err)
	require.Len(cs, 180)

	asset, err := ReadCandlesAsAsset(ctx, symbol)
	require.Nil(err)
	require.NotNil(asset)
	var v float64
	for i, h := range asset.High {
		l := asset.Low[i]
		v += ((h - l) / l)
	}
	v = v / float64(len(asset.High))
	log.Println(v)
	tr, atr := indicator.Atr(14, asset.High, asset.Low, asset.Closing)
	log.Println(tr)
	log.Println(atr)
}
