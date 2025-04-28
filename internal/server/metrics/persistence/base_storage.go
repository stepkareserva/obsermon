package persistence

import "github.com/stepkareserva/obsermon/internal/server/metrics/service"

// requirements to service to be wrappable onto persistent service
type BaseStorage interface {
	service.Storage
}
