package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type SqlDB struct {
	dbConn string

	cancel context.CancelFunc
	wg     sync.WaitGroup

	db atomic.Pointer[sql.DB]
}

var _ Db = (*SqlDB)(nil)

func NewSqlDB(dbConn string, log *zap.Logger) *SqlDB {
	if log == nil {
		log = zap.NewNop()
	}

	ctx, cancel := context.WithCancel(context.Background())

	d := &SqlDB{
		dbConn: dbConn,
		cancel: cancel,
	}

	// run connection loop
	d.wg.Add(1)
	go d.connectionLoop(ctx, log)

	return d
}

func (d *SqlDB) connectionLoop(ctx context.Context, log *zap.Logger) {
	defer d.wg.Done()

	d.keepConnection(ctx, log)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.keepConnection(ctx, log)
		case <-ctx.Done():
			return
		}
	}
}

func (d *SqlDB) keepConnection(ctx context.Context, log *zap.Logger) {
	if d.PingContext(ctx) == nil {
		// everything is fine, db connected
		return
	}

	// try to reconnect
	sqlDB, err := sql.Open("pgx", d.dbConn)
	if err != nil {
		log.Warn("db reconnect", zap.Error(err))
		return
	}

	// set connection params (no idea what's good for our service)
	sqlDB.SetMaxOpenConns(16)
	sqlDB.SetMaxIdleConns(8)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(300 * time.Minute)

	// use db if success
	d.db.Store(sqlDB)
}

func (d *SqlDB) Close() error {
	if d == nil {
		return nil
	}
	d.cancel()
	d.wg.Wait()

	db := d.db.Load()
	if db == nil {
		return nil
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("db closing: %w", err)
	}

	return nil
}

func (d *SqlDB) BeginTx(ctx context.Context) (Tx, error) {
	if d == nil {
		return nil, fmt.Errorf("db not exists")
	}

	db := d.db.Load()
	if db == nil {
		return nil, fmt.Errorf("db not connected")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return &sqlTx{Tx: tx}, nil
}

func (d *SqlDB) PingContext(ctx context.Context) error {
	if d == nil {
		return fmt.Errorf("db not exists")
	}

	db := d.db.Load()
	if db == nil {
		return fmt.Errorf("db not connected")
	}

	return db.PingContext(ctx)
}
