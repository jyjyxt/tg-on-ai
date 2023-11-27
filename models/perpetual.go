package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func perpetualFromRow(row Row) (*Perpetual, error) {
	var p Perpetual
	err := row.Scan(&p.Symbol, &p.BaseAsset, &p.QuoteAsset, &p.Categories, &p.Source)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (s *SQLite3Store) CreatePerpetual(ctx context.Context, symbol, base, quote, source string, categories []string) (*Perpetual, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	txn, err := s.db.BeginTx(ctx, nil)
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
	query := BuildInsertionSQL("perpetuals", perpetualCols)
	_, err = txn.ExecContext(ctx, query, p.values()...)
	if err != nil {
		return nil, err
	}
	return p, txn.Commit()
}

func (s *SQLite3Store) ReadPerpetual(ctx context.Context, symbol string) (*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals WHERE symbol=?", strings.Join(perpetualCols, ","))
	row := s.db.QueryRowContext(ctx, query, symbol)
	return perpetualFromRow(row)
}

func (s *SQLite3Store) ReadPerpetualSet(ctx context.Context, source string) (map[string]*Perpetual, error) {
	ps, err := s.ReadPerpetuals(ctx, source)
	if err != nil {
		return nil, err
	}
	filter := make(map[string]*Perpetual, 0)
	for _, p := range ps {
		filter[p.Symbol] = p
	}
	return filter, nil
}

func (s *SQLite3Store) ReadPerpetuals(ctx context.Context, source string) ([]*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals", strings.Join(perpetualCols, ","))
	rows, err := s.db.QueryContext(ctx, query)
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
