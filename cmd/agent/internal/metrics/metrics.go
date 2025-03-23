package metrics

import (
	"errors"

	"github.com/stepkareserva/obsermon/internal/models"
)

type Gauges map[string]models.Gauge
type Counters map[string]models.Counter

type Metrics struct {
	Gauges   Gauges
	Counters Counters
}

func NewMetrics() Metrics {
	return Metrics{
		Gauges:   Gauges{},
		Counters: Counters{},
	}
}

func (m *Metrics) Update(metrics Metrics) error {
	for k, v := range metrics.Gauges {
		m.Gauges[k] = v
	}

	var errs []error
	for k, v := range metrics.Counters {
		updated := m.Counters[k]
		if err := updated.Update(v); err != nil {
			errs = append(errs, err)
		} else {
			m.Counters[k] = updated
		}
	}

	return errors.Join(errs...)
}
