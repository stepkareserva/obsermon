package router

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/stepkareserva/obsermon/internal/models"
)

// tests for update counters handle
func TestValidUpdateCounterHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /update/counter/name/1", func(t *testing.T) {

		mockService.
			EXPECT().
			UpdateCounter(gomock.Any(), gomock.Eq(models.Counter{
				Name:  "name",
				Value: 1,
			})).
			Return(&models.Counter{
				Name:  "name",
				Value: 1,
			}, nil)

		res := testingPostURL(t, ts.URL+"/update/counter/name/1")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestInvalidUpdateCounterHandler(t *testing.T) {
	ctrl, _, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	invalidRequests := []string{
		"/update/counter/name/1.000",
		"/update/counter/name/1.25",
		"/update/counter/name/99999999999999999999999999",
	}
	for _, req := range invalidRequests {
		t.Run(fmt.Sprintf("test %s", req), func(t *testing.T) {
			res := testingPostURL(t, ts.URL+req)
			defer safeCloseRes(t, res)
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestNotFoundUpdateCounterHandler(t *testing.T) {
	ctrl, _, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /update/counter/", func(t *testing.T) {
		res := testingPostURL(t, ts.URL+"/update/counter/")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// tests for update gauges handle
func TestValidUpdateGaugeHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /update/gauge/name/1.0", func(t *testing.T) {

		mockService.
			EXPECT().
			UpdateGauge(gomock.Any(), gomock.Eq(models.Gauge{
				Name:  "name",
				Value: 1.0,
			})).
			Return(&models.Gauge{
				Name:  "name",
				Value: 1.0,
			}, nil)

		res := testingPostURL(t, ts.URL+"/update/gauge/name/1.0")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestInvalidUpdateGaugeHandler(t *testing.T) {
	ctrl, _, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	invalidRequests := []string{
		"/update/gauge/name/value",
		"/update/gauge/name/1.2.3",
	}
	for _, req := range invalidRequests {
		t.Run(fmt.Sprintf("test %s", req), func(t *testing.T) {
			res := testingPostURL(t, ts.URL+req)
			defer safeCloseRes(t, res)
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestNotFoundUpdateGaugeHandler(t *testing.T) {
	ctrl, _, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /update/gauge/", func(t *testing.T) {
		res := testingPostURL(t, ts.URL+"/update/gauge/")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// tests for JSON handlers
func TestValidUpdateCounterJSONHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("update counter: { name, 1 }", func(t *testing.T) {
		counterJSON := `{"id":"name","type":"counter","delta":1}`

		value := models.CounterValue(1)
		counter := models.Metric{
			MType: models.MetricTypeCounter,
			ID:    "name",
			Delta: &value,
		}

		mockService.
			EXPECT().
			UpdateMetric(gomock.Any(), counter).
			Return(&counter, nil)

		res := testingPostJSON(t, ts.URL+"/update", counterJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, counterJSON, string(body))
	})
}

func TestValidUpdateGaugeJSONHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("update gauge: { name, 1.5 }", func(t *testing.T) {
		gaugeJSON := `{"id":"name","type":"gauge","value":1.5}`

		value := models.GaugeValue(1.5)
		gauge := models.Metric{
			MType: models.MetricTypeGauge,
			ID:    "name",
			Value: &value,
		}

		mockService.
			EXPECT().
			UpdateMetric(gomock.Any(), gauge).
			Return(&gauge, nil)

		res := testingPostJSON(t, ts.URL+"/update", gaugeJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, gaugeJSON, string(body))
	})
}

func TestInvalidUpdateJSONHandler(t *testing.T) {
	ctrl, _, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("update invalid", func(t *testing.T) {
		invalidJSON := "{}"

		res := testingPostJSON(t, ts.URL+"/update", invalidJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
