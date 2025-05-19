package dbstorage

import "strings"

var (
	queryReplacer = strings.NewReplacer(
		"{counters}", CountersTable,
		"{gauges}", GaugesTable,
		"{name}", NameColumn,
		"{value}", ValueColumn)

	createCountersQuery = queryReplacer.Replace(`
		CREATE TABLE IF NOT EXISTS {counters} (
			{name} TEXT PRIMARY KEY,
			{value} BIGINT NOT NULL
		)`)

	createGaugesQuery = queryReplacer.Replace(`
			CREATE TABLE IF NOT EXISTS {gauges} (
			{name} TEXT PRIMARY KEY,
			{value} DOUBLE PRECISION NOT NULL
		)`)

	insertCounterQuery = queryReplacer.Replace(`
		INSERT
			INTO {counters} ({name}, {value})
			VALUES ($1, $2)
		`)

	updateCounterQuery = queryReplacer.Replace(`
		UPDATE {counters}
			SET {value} = $2
			WHERE {name} = $1
		`)

	findCounterQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {counters}
			WHERE {name} = $1
		`)

	listCountersQuery = queryReplacer.Replace(`
		SELECT {name}, {value}
			FROM {counters}
		`)

	selectCounterForUpdateQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {counters}
			WHERE {name} = $1
		FOR UPDATE
		`)

	clearCountersQuery = queryReplacer.Replace(`
		DELETE FROM {counters}
	`)

	setGaugeQuery = queryReplacer.Replace(`
		INSERT
			INTO {gauges} ({name}, {value})
			VALUES ($1, $2)
		ON CONFLICT ({name})
			DO UPDATE SET {value} = EXCLUDED.{value}
		`)

	insertGaugeQuery = queryReplacer.Replace(`
		INSERT
			INTO {gauges} ({name}, {value})
			VALUES ($1, $2)
		`)

	findGaugeQuery = queryReplacer.Replace(`
		SELECT {value}
			FROM {gauges}
			WHERE {name} = $1
		`)

	listGaugesQuery = queryReplacer.Replace(`
		SELECT {name}, {value}
			FROM {gauges}
		`)

	clearGaugeQuery = queryReplacer.Replace(`
		DELETE FROM {gauges}
	`)
)
