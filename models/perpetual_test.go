package models

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPerpetual(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	path := "/tmp/test.sqlite3"
	store, err := OpenDataSQLite3Store(path)
	require.Nil(err)
	defer os.Remove(path)

	p, err := store.CreatePerpetual(ctx, "ETHBTC", "ETH", "BTC", "binance", []string{"pos"})
	require.Nil(err)
	require.NotNil(p)
	p, err = store.ReadPerpetual(ctx, "ETHBTC")
	require.Nil(err)
	require.NotNil(p)
	ps, err := store.ReadPerpetuals(ctx, PerpetualSourceBinance)
	require.Nil(err)
	require.Len(ps, 1)
	filter, err := store.ReadPerpetualSet(ctx, PerpetualSourceBinance)
	require.Nil(err)
	require.Len(filter, 1)
}
