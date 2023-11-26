package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Perpetual struct {
	Symbol string
}

var perpetualCols = []string{"symbol"}

func (p *Perpetual) values() []any {
	return []any{p.Symbol}
}

func perpetualFromRow(row *sql.Row) (*Perpetual, error) {
	var p Perpetual
	err := row.Scan(&p.Symbol)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (s *SQLite3Store) CreatePerpetual(ctx context.Context, symbol string) (*Perpetual, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	txn, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer txn.Rollback()

	p := &Perpetual{
		Symbol: symbol,
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
