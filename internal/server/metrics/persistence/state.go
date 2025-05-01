package persistence

import "github.com/stepkareserva/obsermon/internal/models"

type State struct {
	Counters []models.Counter `json:"counters"`
	Gauges   []models.Gauge   `json:"gauge"`
}
