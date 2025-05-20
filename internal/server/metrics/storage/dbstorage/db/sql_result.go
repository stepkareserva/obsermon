package db

import "database/sql"

// Result impl
type sqlResult struct {
	sql.Result
}

var _ Result = (*sqlResult)(nil)
