package models

import (
	"context"
	"log"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/stretchr/testify/require"
)

func testAPI(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	client := futures.NewClient("", "")
	info, err := client.NewExchangeInfoService().Do(ctx)
	require.Nil(err)
	log.Printf("%#v", info)
}
