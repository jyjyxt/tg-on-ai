package services

import "context"

func Root(ctx context.Context) {
	go LoopingExchangeInfo(ctx)
	go LoopingPremiumIndex(ctx)
	go LoopingOpenInterestHist(ctx)
	go LoopingCoinGecko(ctx)
	LoopingCandle(ctx)
}
