package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCounter(t *testing.T) {
	counter := models.Counter{
		Name:  "name",
		Value: models.CounterValue(2),
	}
	expectedURLPath := "/updates"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, expectedURLPath)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	metricsClient, err := New(mockServer.URL, "")
	require.NoError(t, err)

	err = metricsClient.UpdateCounter(counter)
	require.NoError(t, err)
}

func TestUpdateGauge(t *testing.T) {
	gauge := models.Gauge{
		Name:  "name",
		Value: models.GaugeValue(2.5),
	}
	expectedURLPath := "/updates"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, expectedURLPath)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	metricsClient, err := New(mockServer.URL, "")
	require.NoError(t, err)

	err = metricsClient.UpdateGauge(gauge)
	require.NoError(t, err)
}
