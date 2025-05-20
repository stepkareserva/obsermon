package dbstorage

import (
	"fmt"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/metrics/storage/dbstorage/db"
)

func ScanCounters(rows db.Rows) ([]models.Counter, error) {
	var counters []models.Counter
	var counter models.Counter
	for rows.Next() {
		if err := rows.Scan(&counter.Name, &counter.Value); err != nil {
			return nil, fmt.Errorf("counter row scan: %w", err)
		}
		counters = append(counters, counter)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return counters, nil
}

func ScanCounter(rows db.Rows) (*models.Counter, error) {
	counters, err := ScanCounters(rows)
	if err != nil {
		return nil, err
	}
	switch len(counters) {
	case 0:
		return nil, nil
	case 1:
		counter := counters[0]
		return &counter, nil
	default:
		return nil, fmt.Errorf("more than one counters with the same name")
	}
}

func ScanGauges(rows db.Rows) ([]models.Gauge, error) {
	var gauges []models.Gauge
	var gauge models.Gauge
	for rows.Next() {
		if err := rows.Scan(&gauge.Name, &gauge.Value); err != nil {
			return nil, fmt.Errorf("gauge row scan: %w", err)
		}
		gauges = append(gauges, gauge)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return gauges, nil
}

func ScanGauge(rows db.Rows) (*models.Gauge, error) {
	gauges, err := ScanGauges(rows)
	if err != nil {
		return nil, err
	}
	switch len(gauges) {
	case 0:
		return nil, nil
	case 1:
		gauge := gauges[0]
		return &gauge, nil
	default:
		return nil, fmt.Errorf("more than one gauges with the same name")
	}
}
