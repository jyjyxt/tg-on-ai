package models

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/stretchr/testify/require"
)

func TestCandle(t *testing.T) {
	require := require.New(t)
	ctx := setup()
	defer os.Remove(pathTest)

	symbol := "XMRUSDT"
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
		candle, err = CreateCandle(ctx, symbol, in.Open, in.High, in.Low, in.Close, in.Volume, in.OpenTime, in.CloseTime)
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
	log.Println(ReadPeakAndTrough(v*5, asset))
}
