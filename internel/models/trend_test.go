package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrend(t *testing.T) {
	require := require.New(t)
	ctx := setup()
	defer os.Remove(pathTest)

	sym := "btc"
	category := "day3up"
	high := 1.23
	low := 1.1
	up := 1.2
	down := 1.3

	trend, err := UpsertTrend(ctx, sym, category, high, low, up, down)
	require.Nil(err)
	require.NotNil(trend)

	old, err := FindTrend(ctx, sym, category)
	require.Nil(err)
	require.NotNil(old)
	require.Equal(high, old.High)
	require.Equal(low, old.Low)
	require.Equal(up, old.Up)

	trend, err = UpsertTrend(ctx, sym, category, high, low, up, down)
	require.Nil(err)
	require.NotNil(trend)

	old, err = FindTrend(ctx, sym, category)
	require.Nil(err)
	require.NotNil(old)
	require.Equal(high, old.High)
	require.Equal(low, old.Low)
	require.Equal(1.2, old.Up)

	filter, err := FindTrendSet(ctx, category)
	require.Nil(err)
	require.Len(filter, 1)
}
