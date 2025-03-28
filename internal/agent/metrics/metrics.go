package metrics

import (
	"github.com/stepkareserva/obsermon/internal/models"
)

type Metrics struct {
	Gauges   models.GaugesMap
	Counters models.CountersMap
}

func NewMetrics() Metrics {
	return Metrics{
		Gauges:   models.GaugesMap{},
		Counters: models.CountersMap{},
	}
}

func (m *Metrics) Update(metrics Metrics) error {
	m.Gauges.Update(metrics.Gauges)
	if err := m.Counters.Update(metrics.Counters); err != nil {
		return err
	}

	return nil
}
