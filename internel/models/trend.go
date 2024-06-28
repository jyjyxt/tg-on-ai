package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"tg.ai/internel/session"
)

type Trend struct {
	Symbol    string
	Category  string
	High      float64
	Low       float64
	Value     float64
	UpdatedAt time.Time
}

var trendsColumns = []string{"symbol", "category", "high", "low", "value", "updated_at"}

func (t *Trend) values() []any {
	return []any{t.Symbol, t.Category, t.High, t.Low, t.Value, t.UpdatedAt}
}

func trendFromRow(row session.Row) (*Trend, error) {
	var t Trend
	err := row.Scan(&t.Symbol, &t.Category, &t.High, &t.Low, &t.Value, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func UpsertTrend(ctx context.Context, symbol, category string, h, l, v float64) (*Trend, error) {
	old, err := FindTrend(ctx, symbol, category)
	if err != nil {
		return nil, err
	}
	t := &Trend{
		Symbol:    symbol,
		Category:  category,
		High:      h,
		Low:       l,
		Value:     v,
		UpdatedAt: time.Now(),
	}
	err = session.SqliteDB(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if old == nil {
			query := session.BuildInsertionSQL("trends", trendsColumns)
			_, err := tx.ExecContext(ctx, query, t.values()...)
			return err
		}
		_, err = tx.ExecContext(ctx, "UPDATE trends SET high=?, low=?, value=?, updated_at=? WHERE symbol=? AND category=?", t.High, t.Low, t.Value, t.UpdatedAt, t.Symbol, t.Category)
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
		query := fmt.Sprintf("SELECT %s FROM trends WHERE category=?", strings.Join(trendsColumns, ","))
		old, err := findRecentTrends(ctx, tx, query, category)
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

func findRecentTrends(ctx context.Context, tx *sql.Tx, query string, args ...any) ([]*Trend, error) {
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
