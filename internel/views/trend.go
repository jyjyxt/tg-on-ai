package views

import (
	"math"

	"tg.ai/internel/models"
)

type Trend struct {
	Symbol string
	High   float64
	Low    float64
	Now    float64
	Up     float64
	Down   float64
}

func buildTrend(a *models.Trend) *Trend {
	b := Trend{
		Symbol: a.Symbol,
		High:   a.High,
		Low:    a.Low,
		Now:    a.Now,
		Up:     math.Floor(a.Up*10000) / 100,
		Down:   math.Floor(a.Down*10000) / 100,
	}
	return &b
}

func RenderTrend(d *models.Trend) *Trend {
	return buildTrend(d)
}

func RenderTrends(s []*models.Trend) []*Trend {
	views := make([]*Trend, len(s))
	for i, a := range s {
		views[i] = buildTrend(a)
	}
	return views
}