package dbstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stepkareserva/obsermon/internal/server/metrics/storage/dbstorage/db"
)

type UnitOfWork struct {
	db          db.Db
	retryPolicy []time.Duration
}

func NewUoW(db db.Db, retryPolicy []time.Duration) UnitOfWork {
	return UnitOfWork{db: db, retryPolicy: retryPolicy}
}

func (uow *UnitOfWork) Do(ctx context.Context, fn func(context.Context, db.Tx) error) error {
	if uow == nil || uow.db == nil {
		return fmt.Errorf("uow not exists")
	}

	var err error
	if err = uow.do(ctx, fn); err == nil || !errors.Is(err, ErrNet) {
		return err
	}

	for _, delay := range uow.retryPolicy {
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			return fmt.Errorf("uow do timeout")
		case <-timer.C:
			if err = uow.do(ctx, fn); err == nil || !errors.Is(err, ErrNet) {
				return err
			}
		}
	}

	return fmt.Errorf("all sustained op attempts failed, last error: %w", err)
}

func (uow *UnitOfWork) do(ctx context.Context, fn func(context.Context, db.Tx) error) error {
	tx, err := uow.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	err = fn(ctx, tx)

	if err == nil {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("tx commit: %w", err)
		}
		return nil
	}

	if rbErr := tx.Rollback(); rbErr != nil {
		return errors.Join(
			fmt.Errorf("tx: %w", err),
			fmt.Errorf("tx rollback: %w", rbErr))
	}

	return fmt.Errorf("tx: %w", err)
}
