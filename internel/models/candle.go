package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/cinar/indicator"
	"github.com/shopspring/decimal"
	"tg.ai/internel/session"
)

type Candle struct {
	Symbol    string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	OpenTime  int64
	CloseTime int64
}

var candleCols = []string{"symbol", "open", "high", "low", "close", "volume", "open_time", "close_time"}

func (c *Candle) values() []any {
	return []any{c.Symbol, c.Open, c.High, c.Low, c.Close, c.Volume, c.OpenTime, c.CloseTime}
}

func candleFromRow(row session.Row) (*Candle, error) {
	var c Candle
	err := row.Scan(&c.Symbol, &c.Open, &c.High, &c.Low, &c.Close, &c.Volume, &c.OpenTime, &c.CloseTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

func UpsertCandle(ctx context.Context, symbol string, line *futures.Kline) (*Candle, error) {
	c := &Candle{
		Symbol:    symbol,
		Open:      decimal.RequireFromString(line.Open).InexactFloat64(),
		High:      decimal.RequireFromString(line.High).InexactFloat64(),
		Low:       decimal.RequireFromString(line.Low).InexactFloat64(),
		Close:     decimal.RequireFromString(line.Close).InexactFloat64(),
		Volume:    decimal.RequireFromString(line.Volume).InexactFloat64(),
		OpenTime:  line.OpenTime,
		CloseTime: line.CloseTime,
	}
	old, err := ReadCandle(ctx, c.Symbol, c.OpenTime)
	if err != nil {
		return nil, err
	}

	s := session.SqliteDB(ctx)
	s.Lock()
	defer s.Unlock()

	txn, err := s.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()
	if old != nil {
		err = s.ExecOne(ctx, txn, "UPDATE candles SET high=?, low=?, close=?, volume=? WHERE symbol=? AND open_time=?", c.High, c.Low, c.Close, c.Volume, c.Symbol, c.OpenTime)
		if err != nil {
			return nil, err
		}
		return c, txn.Commit()
	}
	query := session.BuildInsertionSQL("candles", candleCols)
	_, err = txn.ExecContext(ctx, query, c.values()...)
	if err != nil {
		return nil, err
	}
	return c, txn.Commit()
}

func LatestCandleTime(ctx context.Context, symbol string) (*Candle, error) {
	s := session.SqliteDB(ctx)
	query := fmt.Sprintf("SELECT %s FROM candles WHERE symbol=? ORDER BY symbol,open_time DESC LIMIT 1", strings.Join(candleCols, ","))
	row := s.QueryRow(ctx, query, symbol)
	c, err := candleFromRow(row)
	if err != nil || c != nil {
		return c, err
	}
	return &Candle{OpenTime: time.Now().Add(time.Hour * 24 * -30).UnixMilli()}, nil
}

func ReadCandlesAsAsset(ctx context.Context, symbol string) (*indicator.Asset, error) {
	candles, err := ReadCandles(ctx, symbol)
	if err != nil || len(candles) == 0 {
		return nil, err
	}
	l := len(candles)
	asset := &indicator.Asset{
		Date:    make([]time.Time, l),
		Opening: make([]float64, l),
		Closing: make([]float64, l),
		High:    make([]float64, l),
		Low:     make([]float64, l),
		Volume:  make([]float64, l),
	}
	for i, c := range candles {
		asset.Date[i] = time.UnixMilli(c.OpenTime)
		asset.Opening[i] = c.Open
		asset.Closing[i] = c.Close
		asset.High[i] = c.High
		asset.Low[i] = c.Low
		asset.Volume[i] = c.Volume
	}
	return asset, nil
}

func ReadCandles(ctx context.Context, symbol string) ([]*Candle, error) {
	s := session.SqliteDB(ctx)
	query := fmt.Sprintf("SELECT %s FROM candles WHERE symbol=? AND open_time>? ORDER BY symbol,open_time", strings.Join(candleCols, ","))
	rows, err := s.Query(ctx, query, symbol, time.Now().Add(time.Hour*24*-90).UnixMilli())
	if err != nil {
		return nil, err
	}
	var cs []*Candle
	for rows.Next() {
		c, err := candleFromRow(rows)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func ReadLastCandles(ctx context.Context) ([]*Candle, error) {
	s := session.SqliteDB(ctx)
	query := fmt.Sprintf("SELECT %s FROM candles WHERE open_time>?", strings.Join(candleCols, ","))
	rows, err := s.Query(ctx, query, time.Now().Add(time.Hour*-4).UnixMilli())
	if err != nil {
		return nil, err
	}
	var cs []*Candle
	for rows.Next() {
		c, err := candleFromRow(rows)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func ReadHighestCandles(ctx context.Context) (map[string]float64, error) {
	s := session.SqliteDB(ctx)
	query := "SELECT symbol,MAX(high) FROM candles WHERE open_time>? GROUP BY symbol"
	rows, err := s.Query(ctx, query, time.Now().Add(time.Hour*-36).UnixMilli())
	if err != nil {
		return nil, err
	}
	filters := make(map[string]float64, 0)
	for rows.Next() {
		var symbol string
		var high float64
		err = rows.Scan(&symbol, &high)
		if err != nil {
			return nil, err
		}
		filters[symbol] = high
	}
	return filters, nil
}

func ReadCandle(ctx context.Context, symbol string, open int64) (*Candle, error) {
	query := fmt.Sprintf("SELECT %s FROM candles WHERE symbol=? AND open_time=?", strings.Join(candleCols, ","))
	row := session.SqliteDB(ctx).QueryRow(ctx, query, symbol, open)
	return candleFromRow(row)
}

func DeleteCandles(ctx context.Context, symbol string) error {
	s := session.SqliteDB(ctx)
	s.Lock()
	defer s.Unlock()

	txn, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	query := "DELETE FROM candles WHERE symbol=? AND open_time<?"
	_, err = txn.ExecContext(ctx, query, symbol, time.Now().Add(time.Hour*24*-365).UnixMilli())
	if err != nil {
		return err
	}
	return txn.Commit()
}

func ReadVolatility(ctx context.Context, asset *indicator.Asset) float64 {
	var v float64
	for i, h := range asset.High {
		l := asset.Low[i]
		v += ((h - l) / l)
	}
	v = v / float64(len(asset.High))
	return v * 3
}
