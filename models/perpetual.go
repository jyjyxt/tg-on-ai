package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"tg-on-ai/session"
)

const (
	PerpetualSourceBinance = "binance"
)

type Perpetual struct {
	Symbol     string
	BaseAsset  string
	QuoteAsset string
	Categories string
	Source     string
}

var perpetualCols = []string{"symbol", "base_asset", "quote_asset", "categories", "source"}

func (p *Perpetual) values() []any {
	return []any{p.Symbol, p.BaseAsset, p.QuoteAsset, p.Categories, p.Source}
}

func perpetualFromRow(row session.Row) (*Perpetual, error) {
	var p Perpetual
	err := row.Scan(&p.Symbol, &p.BaseAsset, &p.QuoteAsset, &p.Categories, &p.Source)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func CreatePerpetual(ctx context.Context, symbol, base, quote, source string, categories []string) (*Perpetual, error) {
	s := session.SqliteDB(ctx)
	s.Lock()
	defer s.Unlock()

	txn, err := s.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	p := &Perpetual{
		Symbol:     symbol,
		BaseAsset:  base,
		QuoteAsset: quote,
		Categories: strings.ToLower(strings.Join(categories, ",")),
		Source:     source,
	}
	query := session.BuildInsertionSQL("perpetuals", perpetualCols)
	_, err = txn.ExecContext(ctx, query, p.values()...)
	if err != nil {
		return nil, err
	}
	return p, txn.Commit()
}

func ReadPerpetual(ctx context.Context, symbol string) (*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals WHERE symbol=?", strings.Join(perpetualCols, ","))
	row := session.SqliteDB(ctx).QueryRow(ctx, query, symbol)
	return perpetualFromRow(row)
}

func ReadPerpetualSet(ctx context.Context, source string) (map[string]*Perpetual, error) {
	ps, err := ReadPerpetuals(ctx, source)
	if err != nil {
		return nil, err
	}
	filter := make(map[string]*Perpetual, 0)
	for _, p := range ps {
		filter[p.Symbol] = p
	}
	return filter, nil
}

func ReadPerpetuals(ctx context.Context, source string) ([]*Perpetual, error) {
	s := session.SqliteDB(ctx)
	query := fmt.Sprintf("SELECT %s FROM perpetuals", strings.Join(perpetualCols, ","))
	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var ps []*Perpetual
	for rows.Next() {
		p, err := perpetualFromRow(rows)
		if err != nil {
			return nil, err
		}
		if p.Source != source {
			continue
		}
		ps = append(ps, p)
	}
	return ps, nil
}
