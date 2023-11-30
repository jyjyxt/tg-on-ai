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

	strategy, err := CreateStrategy(ctx, "ETHBTC", "macd", 1, -80)
	require.Nil(err)
	require.NotNil(strategy)
	strategies, err := ReadStrategies(ctx, "ETHBTC")
	require.Nil(err)
	require.Len(strategies, 1)
	strategy, err = ReadStrategy(ctx, "ETHBTC", "macd")
	require.Nil(err)
	require.NotNil(strategy)
	require.Equal(1, strategy.Action)
	require.Equal(float64(-80), strategy.Score)
}
