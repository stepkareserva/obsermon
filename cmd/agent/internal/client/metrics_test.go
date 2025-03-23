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
	counterName := "name"
	counterValue := models.Counter(2)
	expectedURLPath := "/update/counter/name/2"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, expectedURLPath)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	metricsClient, err := NewMetricsClient(mockServer.URL)
	require.NoError(t, err)

	err = metricsClient.UpdateCounter(counterName, counterValue)
	require.NoError(t, err)
}

func TestUpdateGauge(t *testing.T) {
	gaugeName := "name"
	gaugeValue := models.Gauge(2.5)
	expectedURLPath := "/update/gauge/name/2.5"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, expectedURLPath)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	metricsClient, err := NewMetricsClient(mockServer.URL)
	require.NoError(t, err)

	err = metricsClient.UpdateGauge(gaugeName, gaugeValue)
	require.NoError(t, err)
}
