package models

import "fmt"

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type Metrics struct {
	// metric's name
	ID string `json:"id" validate:"required"`
	// "gauge" or "counter"
	MType MetricType `json:"type" validate:"required,oneof=gauge counter"`
	// value if counter
	Delta *CounterValue `json:"delta,omitempty"`
	// value if gauge
	Value *GaugeValue `json:"value,omitempty"`
}

func CounterMetric(counter Counter) Metrics {
	return Metrics{
		ID:    counter.Name,
		MType: MetricTypeCounter,
		Delta: &counter.Value,
	}
}

func GaugeMetric(gauge Gauge) Metrics {
	return Metrics{
		ID:    gauge.Name,
		MType: MetricTypeGauge,
		Value: &gauge.Value,
	}
}

func (m *Metrics) Counter() (*Counter, error) {
	if m.MType != MetricTypeCounter {
		return nil, fmt.Errorf("invalid metric type")
	}
	if m.Delta == nil {
		return nil, fmt.Errorf("invalid metric value")
	}
	return &Counter{
		Name:  m.ID,
		Value: *m.Delta,
	}, nil
}

func (m *Metrics) Gauge() (*Gauge, error) {
	if m.MType != MetricTypeGauge {
		return nil, fmt.Errorf("invalid metric type")
	}
	if m.Value == nil {
		return nil, fmt.Errorf("invalid metric value")
	}
	return &Gauge{
		Name:  m.ID,
		Value: *m.Value,
	}, nil
}
