package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"tg-on-ai/session"
)

const (
	StrategyNameMACD      = "MACD"
	StrategyNameKDJ       = "KDJ"
	StrategyNameAroon     = "Aroon"
	StrategyNameWilliamsR = "WilliamsR"
)

type Strategy struct {
	Symbol   string
	Name     string
	Action   int64
	ScoreX   float64
	ScoreY   float64
	OpenTime int64
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
		return fmt.Sprintf("%s:%.2f-%.2f", s.Name, s.ScoreX, s.ScoreY)
	case StrategyNameWilliamsR:
		return fmt.Sprintf("%s:%.2f", s.Name, s.ScoreX)
	}
	return ""
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

func ReadStrategy(ctx context.Context, symbol, name string) (*Strategy, error) {
	query := fmt.Sprintf("SELECT %s FROM strategies WHERE symbol=? AND name=?", strings.Join(strategyCols, ","))
	row := session.SqliteDB(ctx).QueryRow(ctx, query, symbol, name)
	return strategyFromRow(row)
}
