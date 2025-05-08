package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/mocks"
)

// test for counter value handler
func TestValidUpdatesHandler(t *testing.T) {
	ctrl, mockService, ts := getUpdatesTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("updats metrics: [{ name, 1 }, {other,2.5}]", func(t *testing.T) {
		metricsJSON := `[
			{"id":"name", "type":"counter", "delta":1},
			{"id":"other", "type":"gauge", "value":2.5}
		]`

		counterValue := models.CounterValue(1)
		gaugeValue := models.GaugeValue(2.5)
		metrics := models.Metrics{
			{
				MType: models.MetricTypeCounter,
				ID:    "name",
				Delta: &counterValue,
			},
			{
				MType: models.MetricTypeGauge,
				ID:    "other",
				Value: &gaugeValue,
			},
		}

		mockService.
			EXPECT().
			UpdateMetrics(gomock.Any(), metrics).
			Return(metrics, nil)

		res := testingPostJSON(t, ts.URL+"/", metricsJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, metricsJSON, string(body))
	})
}

func TestInvalidUpdatesHandler(t *testing.T) {
	ctrl, _, ts := getUpdatesTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("updates invalid metrics", func(t *testing.T) {
		invalidJSON := `[
			{"id":"name", "field":"counter", "delta":1},
			{"id":"other", "other":"gauge", "value":2.5}
		]`
		res := testingPostJSON(t, ts.URL+"/", invalidJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func getUpdatesTestObjects(t *testing.T) (*gomock.Controller, *mocks.MockService, *httptest.Server) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	handlers, err := New(mockService, zap.NewNop())
	require.NoError(t, err, "handlers initialization error")

	router := chi.NewRouter()
	handlers.registerUpdatesRoutes(context.TODO(), router)

	ts := httptest.NewServer(router)

	return ctrl, mockService, ts
}
