package views

import (
	"tg.ai/internel/models"
)

type Trend struct {
	Symbol string   `json:"symbol"`
	High   float64  `json:"high"`
	Low    float64  `json:"low"`
	Now    float64  `json:"now"`
	Up     float64  `json:"up"`
	Down   *float64 `json:"down"`
}

func buildTrend(a *models.Trend) *Trend {
	b := Trend{
		Symbol: a.Symbol,
		High:   a.High,
		Low:    a.Low,
		Now:    a.Now,
		Up:     a.GetUp(),
	}
	if a.Category != models.TrendDaysPath {
		d := a.GetDown()
		b.Down = &d
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
