package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"tg.ai/internel/session"
)

const (
	StrategyNameMACD      = "MACD"
	StrategyNameKDJ       = "KDJ"
	StrategyNameAroon     = "Aroon"
	StrategyNameWilliamsR = "WilliamsR"
	StrategyNameATR       = "ATR"
	StrategyNameWeek      = "Week"
	StrategyNameWeekTwo   = "WeekTwo"
	StrategyNameWeekFour  = "WeekFour"

	StrategyTotal int64 = 3
)

type Strategy struct {
	Symbol   string
	Name     string
	Action   int64
	ScoreX   float64
	ScoreY   float64
	OpenTime int64

	Score int64
}

var strategyCols = []string{"symbol", "name", "action", "score_x", "score_y", "open_time"}

func (s *Strategy) values() []any {
	return []any{s.Symbol, s.Name, s.Action, s.ScoreX, s.ScoreY, s.OpenTime}
}

func strategyFromRow(row session.Row) (*Strategy, error) {
	var s Strategy
	err := row.Scan(&s.Symbol, &s.Name, &s.Action, &s.ScoreX, &s.ScoreY, &s.OpenTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (s *Strategy) Result() string {
	switch s.Name {
	case StrategyNameMACD:
		return fmt.Sprintf("%s:%d", s.Name, s.Action)
	case StrategyNameAroon:
		return fmt.Sprintf("%s:%d/%d", s.Name, int(s.ScoreX), int(s.ScoreY))
	case StrategyNameWilliamsR:
		return fmt.Sprintf("%s:%.2f", s.Name, s.ScoreX)
	case StrategyNameATR:
		return fmt.Sprintf("%s:%.2f:%.2f", s.Name, s.ScoreX, s.ScoreY)
	case StrategyNameWeek:
		return fmt.Sprintf("%s:%.2f:%.2f", s.Name, s.ScoreX, s.ScoreY)
	case StrategyNameWeekTwo:
		return fmt.Sprintf("%s:%.2f:%.2f", s.Name, s.ScoreX, s.ScoreY)
	}
	return ""
}

func (s *Strategy) getAction() int64 {
	switch s.Name {
	case StrategyNameMACD:
		return s.Action
	case StrategyNameAroon:
		if s.ScoreX > 70 && s.ScoreY < 30 {
			return 1
		}
		if s.ScoreY > 70 && s.ScoreX < 30 {
			return -1
		}
		return 0
	case StrategyNameWilliamsR:
		if s.ScoreX < -79 {
			return 1
		}
		if s.ScoreX > -21 {
			return -1
		}
		return 0
	}
	return 0
}

func CreateStrategy(ctx context.Context, symbol, name string, action int64, scoreX, scoreY float64, t int64) (*Strategy, error) {
	st := &Strategy{
		Symbol:   symbol,
		Name:     name,
		Action:   action,
		ScoreX:   scoreX,
		ScoreY:   scoreY,
		OpenTime: t,
	}
	old, err := ReadStrategy(ctx, st.Symbol, st.Name)
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
		err = s.ExecOne(ctx, txn, "UPDATE strategies SET action=?, score_x=?, score_y=?, open_time=? WHERE symbol=? AND name=?", st.Action, st.ScoreX, st.ScoreY, st.OpenTime, st.Symbol, st.Name)
		if err != nil {
			return nil, err
		}
		return st, txn.Commit()
	}
	query := session.BuildInsertionSQL("strategies", strategyCols)
	_, err = txn.ExecContext(ctx, query, st.values()...)
	if err != nil {
		return nil, err
	}
	return st, txn.Commit()
}

func ReadStrategiesAll(ctx context.Context) ([]*Strategy, error) {
	query := fmt.Sprintf("SELECT %s FROM strategies", strings.Join(strategyCols, ","))
	rows, err := session.SqliteDB(ctx).Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var ps []*Strategy
	for rows.Next() {
		p, err := strategyFromRow(rows)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func ReadStrategies(ctx context.Context, symbol string) ([]*Strategy, error) {
	query := fmt.Sprintf("SELECT %s FROM strategies WHERE symbol=?", strings.Join(strategyCols, ","))
	rows, err := session.SqliteDB(ctx).Query(ctx, query, symbol)
	if err != nil {
		return nil, err
	}
	var ps []*Strategy
	for rows.Next() {
		p, err := strategyFromRow(rows)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func ReadATRStrategies(ctx context.Context) ([]*Strategy, error) {
	query := fmt.Sprintf("SELECT %s FROM strategies WHERE name=? AND score_x>0 ORDER BY score_x LIMIT 5", strings.Join(strategyCols, ","))
	return readStrategiesQuery(ctx, query, StrategyNameATR)
}

func ReadWeekStrategies(ctx context.Context, week string, limit int) ([]*Strategy, error) {
	if limit == 0 {
		limit = 5
	}
	query := fmt.Sprintf("SELECT %s FROM strategies WHERE name=? AND score_x>0.95 AND score_x<1 ORDER BY score_y LIMIT %d", strings.Join(strategyCols, ","), limit)
	return readStrategiesQuery(ctx, query, week)
}

func ReadStrategy(ctx context.Context, symbol, name string) (*Strategy, error) {
	query := fmt.Sprintf("SELECT %s FROM strategies WHERE symbol=? AND name=?", strings.Join(strategyCols, ","))
	row := session.SqliteDB(ctx).QueryRow(ctx, query, symbol, name)
	return strategyFromRow(row)
}

func readStrategiesQuery(ctx context.Context, query string, args ...any) ([]*Strategy, error) {
	rows, err := session.SqliteDB(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	var ps []*Strategy
	for rows.Next() {
		p, err := strategyFromRow(rows)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}
