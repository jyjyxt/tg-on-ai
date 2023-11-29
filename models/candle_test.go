package models

import (
	"log"
	"os"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/cinar/indicator"
	"github.com/stretchr/testify/require"
)

func TestCandle(t *testing.T) {
	require := require.New(t)
	ctx := setup()
	defer os.Remove(pathTest)

	symbol := "BTCUSDT"
	candle, err := LatestCandleTime(ctx, symbol)
	require.Nil(err)
	require.NotNil(candle)

	client := futures.NewClient("", "")
	s := client.NewKlinesService()
	s.Symbol(symbol)
	s.Interval("1h")
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
	require.Len(cs, 72)
	err = DeleteCandles(ctx, symbol)
	require.Nil(err)
	cs, err = ReadCandles(ctx, symbol)
	require.Nil(err)
	require.Len(cs, 72)

	asset, err := ReadCandlesAsAsset(ctx, symbol)
	require.Nil(err)
	require.NotNil(asset)
	log.Println(indicator.MacdStrategy(asset))
}
