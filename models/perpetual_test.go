package models

import (
	"context"
	"os"
	"testing"
	"tg-on-ai/session"

	"github.com/stretchr/testify/require"
)

func TestPerpetual(t *testing.T) {
	require := require.New(t)
	path := "/tmp/test.sqlite3"
	store, err := session.OpenDataSQLite3Store(path)
	require.Nil(err)
	ctx := context.Background()
	ctx = session.WithSqliteDB(ctx, store)

	defer os.Remove(path)

	p, err := CreatePerpetual(ctx, "ETHBTC", "ETH", "BTC", PerpetualSourceBinance, []string{"pos"})
	require.Nil(err)
	require.NotNil(p)
	p, err = ReadPerpetual(ctx, "ETHBTC")
	require.Nil(err)
	require.NotNil(p)
	ps, err := ReadPerpetuals(ctx, PerpetualSourceBinance)
	require.Nil(err)
	require.Len(ps, 1)
	filter, err := ReadPerpetualSet(ctx, PerpetualSourceBinance)
	require.Nil(err)
	require.Len(filter, 1)

	p, err = UpdatePerpetual(ctx, "ETHBTC", "0.06977089", "-0.00025906")
	require.Nil(err)
	require.NotNil(p)
	require.Equal(0.06977089, p.MarkPrice)
	require.Equal(-0.00025906, p.LastFundingRate)
}
