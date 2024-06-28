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
	v := 1.2

	trend, err := UpsertTrend(ctx, sym, category, high, low, v)
	require.Nil(err)
	require.NotNil(trend)

	old, err := FindTrend(ctx, sym, category)
	require.Nil(err)
	require.NotNil(old)
	require.Equal(high, old.High)
	require.Equal(low, old.Low)
	require.Equal(v, old.Value)

	trend, err = UpsertTrend(ctx, sym, category, high, low, 1.34)
	require.Nil(err)
	require.NotNil(trend)

	old, err = FindTrend(ctx, sym, category)
	require.Nil(err)
	require.NotNil(old)
	require.Equal(high, old.High)
	require.Equal(low, old.Low)
	require.Equal(1.34, old.Value)

	filter, err := FindTrendSet(ctx, category)
	require.Nil(err)
	require.Len(filter, 1)
}
