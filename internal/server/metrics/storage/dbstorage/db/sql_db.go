package db

import (
	"context"
	"database/sql"
	"fmt"
)

type SqlDB struct {
	*sql.DB
}

var _ Db = (*SqlDB)(nil)

func NewSqlDB(dbConn string) (*SqlDB, error) {
	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}
	return &SqlDB{DB: db}, nil
}

func (db *SqlDB) Close() error {
	if db == nil || db.DB == nil {
		return nil
	}
	return db.Close()
}

func (db *SqlDB) BeginTx(ctx context.Context) (Tx, error) {
	if db == nil || db.DB == nil {
		return nil, fmt.Errorf("db not exists")
	}
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	return &sqlTx{Tx: tx}, nil
}

func (db *SqlDB) PingContext(ctx context.Context) error {
	if db == nil || db.DB == nil {
		return fmt.Errorf("db not exists")
	}
	return db.DB.PingContext(ctx)
}
