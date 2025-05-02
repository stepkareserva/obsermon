package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	db *sql.DB
}

func New(dbConn string) (*Database, error) {
	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}
	return &Database{db: db}, nil
}

func (db *Database) Close() error {
	if db.db == nil {
		return nil
	}
	if err := db.db.Close(); err != nil {
		return fmt.Errorf("database closing: %w", err)
	}
	return nil
}

func (db *Database) Ping() error {
	if db == nil {
		return fmt.Errorf("db not exists")
	}
	if err := db.db.PingContext(context.TODO()); err != nil {
		return fmt.Errorf("db: %w", err)
	}
	return nil
}
