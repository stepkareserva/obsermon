package router

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/stepkareserva/obsermon/internal/models"
)

// test for counter value handler
func TestValidUpdatesHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
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

		res := testingPostJSON(t, ts.URL+"/updates", metricsJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, metricsJSON, string(body))
	})
}

func TestInvalidUpdatesHandler(t *testing.T) {
	ctrl, _, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("updates invalid metrics", func(t *testing.T) {
		invalidJSON := `[
			{"id":"name", "field":"counter", "delta":1},
			{"id":"other", "other":"gauge", "value":2.5}
		]`
		res := testingPostJSON(t, ts.URL+"/updates", invalidJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
