package models

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"tg.ai/internel/session"
)

const (
	TrendDaysPath = "dayspath"
	TrendDays3    = "days3"
	TrendDays7    = "days7"
	TrendDays15   = "days15"
	TrendDays30   = "days30"
)

type Trend struct {
	Symbol    string
	Category  string
	High      float64
	Low       float64
	Now       float64
	Up        float64
	Down      float64
	UpdatedAt time.Time
}

var trendsColumns = []string{"symbol", "category", "high", "low", "now", "up", "down", "updated_at"}

func (t *Trend) values() []any {
	return []any{t.Symbol, t.Category, t.High, t.Low, t.Now, t.Up, t.Down, t.UpdatedAt}
}

func trendFromRow(row session.Row) (*Trend, error) {
	var t Trend
	err := row.Scan(&t.Symbol, &t.Category, &t.High, &t.Low, &t.Now, &t.Up, &t.Down, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func UpsertTrend(ctx context.Context, symbol, category string, h, l, now, up float64) (*Trend, error) {
	old, err := FindTrend(ctx, symbol, category)
	if err != nil {
		return nil, err
	}
	t := &Trend{
		Symbol:    symbol,
		Category:  category,
		High:      h,
		Low:       l,
		Now:       now,
		Up:        up,
		UpdatedAt: time.Now(),
	}
	if category != TrendDaysPath {
		t.Up = t.Now/t.Low - 1
		t.Down = 1 - t.Now/t.High
	}
	err = session.SqliteDB(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if old == nil {
			query := session.BuildInsertionSQL("trends", trendsColumns)
			_, err := tx.ExecContext(ctx, query, t.values()...)
			return err
		}
		_, err = tx.ExecContext(ctx, "UPDATE trends SET high=?, low=?, up=?, down=?, updated_at=? WHERE symbol=? AND category=?", t.High, t.Low, t.Up, t.Down, t.UpdatedAt, t.Symbol, t.Category)
		return nil
	})
	return t, err
}

func FindTrend(ctx context.Context, symbol, category string) (*Trend, error) {
	var t *Trend
	err := session.SqliteDB(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM trends WHERE symbol=? AND category=?", strings.Join(trendsColumns, ","))
		trend, err := findTrend(ctx, tx, query, symbol, category)
		if err != nil {
			return err
		}
		t = trend
		return nil
	})
	return t, err
}

func findTrend(ctx context.Context, tx *sql.Tx, query string, args ...any) (*Trend, error) {
	row := tx.QueryRowContext(ctx, query, args...)
	return trendFromRow(row)
}

func FindTrendSet(ctx context.Context, category string) (map[string]*Trend, error) {
	var trends []*Trend
	err := session.SqliteDB(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM trends WHERE category=? LIMIT 1000", strings.Join(trendsColumns, ","))
		old, err := findTrends(ctx, tx, query, category)
		if err != nil {
			return err
		}
		trends = old
		return nil
	})
	if err != nil {
		return nil, err
	}
	filter := make(map[string]*Trend, 0)
	for i, t := range trends {
		filter[t.Symbol] = trends[i]
	}
	return filter, nil
}

func findTrends(ctx context.Context, tx *sql.Tx, query string, args ...any) ([]*Trend, error) {
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []*Trend
	for rows.Next() {
		t, err := trendFromRow(rows)
		if err != nil {
			return nil, err
		}

		ts = append(ts, t)
	}
	return ts, nil
}

func (t *Trend) GetUp() float64 {
	if t.Category == TrendDaysPath {
		return t.Up
	}
	return math.Floor(t.Up*10000) / 100
}

func (t *Trend) GetDown() float64 {
	return math.Floor(t.Down*10000) / 100
}
