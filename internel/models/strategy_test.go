package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrategy(t *testing.T) {
	require := require.New(t)
	ctx := setup()
	defer os.Remove(pathTest)

	strategy, err := CreateStrategy(ctx, "ETHBTC", "macd", 1, -80, 0, 1701259200000)
	require.Nil(err)
	require.NotNil(strategy)
	strategies, err := ReadStrategies(ctx, "ETHBTC")
	require.Nil(err)
	require.Len(strategies, 1)
	strategy, err = ReadStrategy(ctx, "ETHBTC", "macd")
	require.Nil(err)
	require.NotNil(strategy)
	require.Equal(int64(1), strategy.Action)
	require.Equal(float64(-80), strategy.ScoreX)
	require.Equal(int64(1701259200000), strategy.OpenTime)

	strategy, err = CreateStrategy(ctx, "ETHBTC", "macd", -1, -60, 1, 1701259100000)
	require.Nil(err)
	require.NotNil(strategy)
	strategy, err = ReadStrategy(ctx, "ETHBTC", "macd")
	require.Nil(err)
	require.NotNil(strategy)
	require.Equal(int64(-1), strategy.Action)
	require.Equal(float64(-60), strategy.ScoreX)
	require.Equal(float64(1), strategy.ScoreY)
	require.Equal(int64(1701259100000), strategy.OpenTime)
}
