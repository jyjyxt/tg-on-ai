package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
	"tg.ai/internel/session"
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

	// https://fapi.binance.com/futures/data/openInterestHist?symbol=btcusdt&period=30m&startTime=1701221700000
	SumOpenInterestValue float64

	UpdatedAt int64
	// https://fapi.binance.com/fapi/v1/fundingInfo
	// FundingRateCap       string
	// FundingRateFloor     string
	// fundingIntervalHours int64

	CoinGecko string

	Trend *Trend
}

var perpetualCols = []string{"symbol", "base_asset", "quote_asset", "categories", "source", "mark_price", "last_funding_rate", "open_interest_value", "updated_at", "coingecko"}

func (p *Perpetual) values() []any {
	return []any{p.Symbol, p.BaseAsset, p.QuoteAsset, p.Categories, p.Source, p.MarkPrice, p.LastFundingRate, p.SumOpenInterestValue, p.UpdatedAt, p.CoinGecko}
}

func perpetualFromRow(row session.Row) (*Perpetual, error) {
	var p Perpetual
	err := row.Scan(&p.Symbol, &p.BaseAsset, &p.QuoteAsset, &p.Categories, &p.Source, &p.MarkPrice, &p.LastFundingRate, &p.SumOpenInterestValue, &p.UpdatedAt, &p.CoinGecko)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (p *Perpetual) GetSumOpenInterestValue() string {
	return decimal.NewFromFloat(p.SumOpenInterestValue).Div(decimal.New(1, 6)).RoundFloor(2).String()
}

func CreatePerpetual(ctx context.Context, sym *futures.Symbol) (*Perpetual, error) {
	s := session.SqliteDB(ctx)
	s.Lock()
	defer s.Unlock()

	txn, err := s.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	p := &Perpetual{
		Symbol:     sym.Symbol,
		BaseAsset:  sym.BaseAsset,
		QuoteAsset: sym.QuoteAsset,
		Categories: strings.ToLower(strings.Join(sym.UnderlyingSubType, ",")),
		Source:     PerpetualSourceBinance,
	}
	query := session.BuildInsertionSQL("perpetuals", perpetualCols)
	_, err = txn.ExecContext(ctx, query, p.values()...)
	if err != nil {
		return nil, err
	}
	return p, txn.Commit()
}

func UpdatePerpetual(ctx context.Context, symbol, markPrice, fundingRate, sumOpenInterestValue string, updatedAt int64) (*Perpetual, error) {
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

	if markPrice != "" {
		p.MarkPrice = decimal.RequireFromString(markPrice).InexactFloat64()
	}
	if fundingRate != "" {
		p.LastFundingRate = decimal.RequireFromString(fundingRate).InexactFloat64()
	}
	if sumOpenInterestValue != "" {
		p.SumOpenInterestValue = decimal.RequireFromString(sumOpenInterestValue).InexactFloat64()
	}
	if updatedAt > 0 {
		p.UpdatedAt = updatedAt
	}
	err = s.ExecOne(ctx, txn, "UPDATE perpetuals SET mark_price=?, last_funding_rate=?, open_interest_value=?, updated_at=? WHERE symbol=?", p.MarkPrice, p.LastFundingRate, p.SumOpenInterestValue, p.UpdatedAt, p.Symbol)
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

func ReadPerpetualByBase(ctx context.Context, symbol string) (*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals WHERE base_asset=? OR base_asset=?", strings.Join(perpetualCols, ","))
	row := session.SqliteDB(ctx).QueryRow(ctx, query, symbol, "1000"+symbol)
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

func ReadDiscretePerpetuals(ctx context.Context, path string) ([]*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals WHERE last_funding_rate>? ORDER BY last_funding_rate DESC LIMIT 3", strings.Join(perpetualCols, ","))
	rate := 0.0002
	if path == "low" {
		query = fmt.Sprintf("SELECT %s FROM perpetuals WHERE last_funding_rate<? ORDER BY last_funding_rate LIMIT 3", strings.Join(perpetualCols, ","))
		rate = -0.0002
	}
	ps, err := findPerpetuals(ctx, query, rate)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func ReadBestPerpetuals(ctx context.Context, action string) ([]*Perpetual, error) {
	ss, err := ReadStrategiesAll(ctx)
	if err != nil {
		return nil, err
	}
	filters := make(map[string]int64, 0)
	for _, s := range ss {
		filters[s.Symbol] += s.getAction()
	}
	var symbols []string
	if action == "sell" {
		for k, v := range filters {
			if v == -StrategyTotal {
				symbols = append(symbols, k)
			}
		}
	}
	if action == "buy" {
		for k, v := range filters {
			if v == StrategyTotal {
				symbols = append(symbols, k)
			}
		}
	}
	if len(symbols) == 0 {
		return nil, nil
	}
	return ReadPerpetualsBySymbols(ctx, symbols)
}

func ReadPerpetualsBySymbols(ctx context.Context, symbols []string) ([]*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals WHERE symbol IN ('%s') ORDER BY open_interest_value DESC LIMIT 64", strings.Join(perpetualCols, ","), strings.Join(symbols, "','"))
	return findPerpetuals(ctx, query)
}

func ReadLowPerpetuals(ctx context.Context) ([]*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals ORDER BY last_funding_rate LIMIT 5", strings.Join(perpetualCols, ","))
	ps, err := findPerpetuals(ctx, query)
	if err != nil {
		return nil, err
	}
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

func ReadPerpetualsFull(ctx context.Context, trend string) ([]*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals", strings.Join(perpetualCols, ","))
	s, err := findPerpetuals(ctx, query)
	if err != nil {
		return nil, err
	}
	filter, err := FindTrendSet(ctx, trend)
	if err != nil {
		return nil, err
	}
	for _, k := range s {
		k.Trend = filter[k.Symbol]
	}
	return s, nil
}

func ReadPerpetuals(ctx context.Context, source string) ([]*Perpetual, error) {
	query := fmt.Sprintf("SELECT %s FROM perpetuals", strings.Join(perpetualCols, ","))
	return findPerpetuals(ctx, query)
}

func findPerpetuals(ctx context.Context, query string, args ...any) ([]*Perpetual, error) {
	s := session.SqliteDB(ctx)
	rows, err := s.Query(ctx, query, args...)
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

func DeletePerpetual(ctx context.Context, symbol string) error {
	s := session.SqliteDB(ctx)
	s.Lock()
	defer s.Unlock()

	txn, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	query := "DELETE FROM perpetuals WHERE symbol=?"
	_, err = txn.ExecContext(ctx, query, symbol)
	if err != nil {
		return err
	}
	return txn.Commit()
}

func Notify(ctx context.Context) string {
	var text string
	weekFour, _ := ReadWeekStrategies(ctx, StrategyNameWeekFour, 0)
	if len(weekFour) > 0 {
		if text != "" {
			text = text + "\n"
		}
		var symbols []string
		for _, s := range weekFour {
			symbols = append(symbols, s.Symbol)
		}
		ps, _ := ReadPerpetualsBySymbols(ctx, symbols)
		text = text + "Week Four:\n---------\n"
		text = text + PerpetualsForHuman(ctx, ps)
	}

	weekTwo, _ := ReadWeekStrategies(ctx, StrategyNameWeekTwo, 0)
	if len(weekTwo) > 0 {
		if text != "" {
			text = text + "\n"
		}
		var symbols []string
		for _, s := range weekTwo {
			symbols = append(symbols, s.Symbol)
		}
		ps, _ := ReadPerpetualsBySymbols(ctx, symbols)
		text = text + "Week Two:\n---------\n"
		text = text + PerpetualsForHuman(ctx, ps)
	}

	week, _ := ReadWeekStrategies(ctx, StrategyNameWeek, 0)
	if len(week) > 0 {
		if text != "" {
			text = text + "\n"
		}
		var symbols []string
		for _, s := range week {
			symbols = append(symbols, s.Symbol)
		}
		ps, _ := ReadPerpetualsBySymbols(ctx, symbols)
		text = text + "Week:\n---------\n"
		text = text + PerpetualsForHuman(ctx, ps)
	}

	ss, _ := ReadATRStrategies(ctx)
	if len(ss) > 0 {
		if text != "" {
			text = text + "\n"
		}
		var symbols []string
		for _, s := range ss {
			symbols = append(symbols, s.Symbol)
		}
		ps, _ := ReadPerpetualsBySymbols(ctx, symbols)
		text = text + "ATR:\n---------\n"
		text = text + PerpetualsForHuman(ctx, ps)
	}
	candles, _ := ReadLastCandles(ctx)
	if len(candles) > 0 {
		var highAll, highBTC, closeBTC float64
		for _, c := range candles {
			highAll += (c.Close / c.High)
			if c.Symbol == "BTCUSDT" {
				highBTC = c.High
				closeBTC = c.Close
			}
		}
		if text != "" {
			text = text + "\n"
		}
		text = text + "Ratio BTC/All:\n---------\n"
		text = text + fmt.Sprintf("%f:%f", closeBTC/highBTC, highAll/float64(len(candles)))
	}
	return text
}

func PerpetualsForHuman(ctx context.Context, ps []*Perpetual) string {
	if len(ps) == 0 {
		return "nothing"
	}
	var tt []string
	for _, p := range ps {
		in := fmt.Sprintf("[%s], %s, Price %s, Rate %f, Value %sM", p.Symbol, p.Categories, strconv.FormatFloat(p.MarkPrice, 'f', -1, 64), p.LastFundingRate*100, p.GetSumOpenInterestValue())
		ss, _ := ReadStrategies(ctx, p.Symbol)
		r := make([]string, len(ss))
		for i, s := range ss {
			r[i] = s.Result()
		}
		tt = append(tt, fmt.Sprintf("%s\n%s", in, strings.Join(r, ",")))
	}
	return strings.Join(tt, "\n")
}
