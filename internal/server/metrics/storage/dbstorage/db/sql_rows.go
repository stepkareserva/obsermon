package db

import "database/sql"

type sqlRows struct {
	*sql.Rows
}

var _ Rows = (*sqlRows)(nil)
