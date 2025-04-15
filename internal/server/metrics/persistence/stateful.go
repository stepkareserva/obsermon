package persistence

import "github.com/stepkareserva/obsermon/internal/models"

type State struct {
	Counters []models.Counter `json:"counters"`
	Gauges   []models.Gauge   `json:"gauge"`
}

type StateStorage interface {
	LoadState() (*State, error)
	StoreState(State) error
}

type Stateful interface {
	GetState() (*State, error)
	LoadState(State) error
}
