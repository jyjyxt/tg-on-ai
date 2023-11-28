package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"tg-on-ai/session"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
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

	// https://fapi.binance.com/fapi/v1/premiumIndex
	MarkPrice       float64
	LastFundingRate float64

	// https://fapi.binance.com/fapi/v1/fundingInfo
	// FundingRateCap       string
	// FundingRateFloor     string
	// fundingIntervalHours int64
}

var perpetualCols = []string{"symbol", "base_asset", "quote_asset", "categories", "source", "mark_price", "last_funding_rate"}

func (p *Perpetual) values() []any {
	return []any{p.Symbol, p.BaseAsset, p.QuoteAsset, p.Categories, p.Source, p.MarkPrice, p.LastFundingRate}
}

func perpetualFromRow(row session.Row) (*Perpetual, error) {
	var p Perpetual
	err := row.Scan(&p.Symbol, &p.BaseAsset, &p.QuoteAsset, &p.Categories, &p.Source, &p.MarkPrice, &p.LastFundingRate)
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

func UpdatePerpetual(ctx context.Context, symbol, markPrice, fundingRate string) (*Perpetual, error) {
	p, err := ReadPerpetual(ctx, symbol)
	if err != nil {
		return nil, err
	} else if p == nil {
		return nil, nil
	}

	s := session.SqliteDB(ctx)
	s.Lock()
	defer s.Unlock()

	txn, err := s.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	p.MarkPrice = decimal.RequireFromString(markPrice).InexactFloat64()
	p.LastFundingRate = decimal.RequireFromString(fundingRate).InexactFloat64()
	err = s.ExecOne(ctx, txn, "UPDATE perpetuals SET mark_price=?, last_funding_rate=? WHERE symbol=?", p.MarkPrice, p.LastFundingRate, p.Symbol)
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

func ReadPerpetualCategory(ctx context.Context) ([]string, error) {
	ps, err := ReadPerpetuals(ctx, "")
	if err != nil {
		return nil, err
	}
	filter := make(map[string]string, 0)
	for _, p := range ps {
		cs := strings.Split(p.Categories, ",")
		for _, c := range cs {
			if filter[c] != "" {
				continue
			}
			if c == "" {
				continue
			}
			filter[c] = c
		}
	}
	return maps.Keys(filter), nil
}

func ReadPerpetualsByCategory(ctx context.Context, category string) ([]*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals WHERE categories LIKE ? ORDER BY last_funding_rate LIMIT 5", strings.Join(perpetualCols, ","))
	return findPerpetuals(ctx, query, "%"+category+"%")
}

func ReadDiscretePerpetuals(ctx context.Context) ([]*Perpetual, error) {
	lower := fmt.Sprintf("SELECT %s FROM perpetuals ORDER BY last_funding_rate DESC LIMIT 3", strings.Join(perpetualCols, ","))
	ps, err := findPerpetuals(ctx, lower)
	if err != nil {
		return nil, err
	}

	higher := fmt.Sprintf("SELECT %s FROM perpetuals ORDER BY last_funding_rate LIMIT 3", strings.Join(perpetualCols, ","))
	pss, err := findPerpetuals(ctx, higher)
	if err != nil {
		return nil, err
	}
	ps = append(ps, pss...)
	return ps, nil
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
	query := fmt.Sprintf("SELECT %s FROM perpetuals", strings.Join(perpetualCols, ","))
	return findPerpetuals(ctx, query)
}

func findPerpetuals(ctx context.Context, query string, args ...any) ([]*Perpetual, error) {
	s := session.SqliteDB(ctx)
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
		ps = append(ps, p)
	}
	return ps, nil
}
