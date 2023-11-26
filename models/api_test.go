package models

import (
	"context"
	"log"
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	client := binance.NewClient("", "")
	info, err := client.NewExchangeInfoService().Do(ctx)
	require.Nil(err)
	log.Printf("%#v", info)
}
