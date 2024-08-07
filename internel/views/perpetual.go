package views

import "tg.ai/internel/models"

type Perpetual struct {
	Symbol     string `json:"symbol"`
	BaseAsset  string `json:"base_asset"`
	QuoteAsset string `json:"quote_asset"`
	Categories string `json:"categories"`
	Source     string `json:"source"`

	MarkPrice            float64 `json:"mark_price"`
	LastFundingRate      float64 `json:"last_funding_rate"`
	SumOpenInterestValue float64 `json:"sum_open_interest_value"`

	UpdatedAt int64  `json:"updated_at"`
	CoinGecko string `json:"coingecko"`

	Trend *Trend `json:"trend"`
}

func buildPerpetual(a *models.Perpetual) *Perpetual {
	b := Perpetual{
		Symbol:               a.Symbol,
		BaseAsset:            a.BaseAsset,
		QuoteAsset:           a.QuoteAsset,
		Categories:           a.Categories,
		Source:               a.Source,
		MarkPrice:            a.MarkPrice,
		LastFundingRate:      a.LastFundingRate,
		SumOpenInterestValue: a.SumOpenInterestValue,
		UpdatedAt:            a.UpdatedAt,
		CoinGecko:            a.CoinGecko,
	}
	if a.Trend != nil {
		b.Trend = buildTrend(a.Trend)
	}
	return &b
}

func RenderPerpetual(d *models.Perpetual) *Perpetual {
	return buildPerpetual(d)
}

func RenderPerpetuals(s []*models.Perpetual) []*Perpetual {
	views := make([]*Perpetual, len(s))
	for i, a := range s {
		views[i] = buildPerpetual(a)
	}
	return views
}
