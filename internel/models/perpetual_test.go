package models

import (
	"os"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/stretchr/testify/require"
)

func TestPerpetual(t *testing.T) {
	require := require.New(t)
	ctx := setup()
	defer os.Remove(pathTest)

	client := futures.NewClient("", "")
	info, err := client.NewExchangeInfoService().Do(ctx)
	require.Nil(err)

	btc := "BTCUSDT"

	var symbol *futures.Symbol
	for _, s := range info.Symbols {
		if s.Symbol == btc {
			symbol = &s
		}
	}
	p, err := CreatePerpetual(ctx, symbol)
	require.Nil(err)
	require.NotNil(p)
	p, err = ReadPerpetual(ctx, btc)
	require.Nil(err)
	require.NotNil(p)
	ps, err := ReadPerpetuals(ctx, PerpetualSourceBinance)
	require.Nil(err)
	require.Len(ps, 1)
	filter, err := ReadPerpetualSet(ctx, PerpetualSourceBinance)
	require.Nil(err)
	require.Len(filter, 1)

	p, err = UpdatePerpetual(ctx, btc, "0.06977089", "-0.00025906", "", 0)
	require.Nil(err)
	require.NotNil(p)
	p, err = ReadPerpetual(ctx, btc)
	require.Nil(err)
	require.NotNil(p)
	require.Equal(0.06977089, p.MarkPrice)
	require.Equal(-0.00025906, p.LastFundingRate)
	err = DeletePerpetual(ctx, btc)
	require.Nil(err)
	p, err = ReadPerpetual(ctx, btc)
	require.Nil(err)
	require.Nil(p)
}
