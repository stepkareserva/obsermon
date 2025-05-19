package dbstorage

import "time"

const (
	CountersTable = "counters"
	GaugesTable   = "gauges"

	NameColumn  = "name"
	ValueColumn = "value"
)

const (
	SQLOpTimeout = 15 * time.Second
)
