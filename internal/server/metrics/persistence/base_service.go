package persistence

import (
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
)

// requirements to service to be wrappable onto persistent service
type BaseService interface {
	handlers.Service
	Stateful
}
