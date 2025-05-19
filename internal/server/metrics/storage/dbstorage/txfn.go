package dbstorage

import (
	"context"

	"github.com/stepkareserva/obsermon/internal/server/metrics/storage/dbstorage/db"
)

type TxFn = func(ctx context.Context, tx db.Tx) error
