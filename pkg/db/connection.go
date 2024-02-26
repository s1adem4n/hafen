package db

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/schema.sql
var schema string

func NewConnection(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.sqlite?_timeout=1000&_journal_mode=WAL&_synchronous=normal")
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, err
	}

	return db, nil
}
