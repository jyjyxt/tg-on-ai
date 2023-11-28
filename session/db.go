package session

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var SCHEMA string

type SQLite3Store struct {
	db    *sql.DB
	mutex *sync.Mutex
}

func OpenDataSQLite3Store(path string) (*SQLite3Store, error) {
	return OpenSQLite3Store(path, SCHEMA)
}

func OpenSQLite3Store(path, schema string) (*SQLite3Store, error) {
	dsn := fmt.Sprintf("file:%s?mode=rwc&_journal_mode=WAL&cache=private", path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	return &SQLite3Store{
		db:    db,
		mutex: new(sync.Mutex),
	}, nil
}

func (s *SQLite3Store) Close() error {
	return s.db.Close()
}

func (s *SQLite3Store) Lock() {
	s.mutex.Lock()
}

func (s *SQLite3Store) Unlock() {
	s.mutex.Unlock()
}

func (s *SQLite3Store) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

func (s *SQLite3Store) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *SQLite3Store) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *SQLite3Store) ExecOne(ctx context.Context, tx *sql.Tx, sql string, params ...any) error {
	res, err := tx.ExecContext(ctx, sql, params...)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows != 1 {
		return fmt.Errorf("SQLite3Store.execOne(%s) => %d %v", sql, rows, err)
	}
	return nil
}

type Row interface {
	Scan(dest ...any) error
}

func BuildInsertionSQL(table string, cols []string) string {
	vals := strings.Repeat("?, ", len(cols))
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), vals[:len(vals)-2])
}
