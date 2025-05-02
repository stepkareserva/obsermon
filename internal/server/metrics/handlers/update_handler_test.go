package handlers

import (
	"context"
	"fmt"
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

// tests for update counters handle
func TestValidUpdateCounterHandler(t *testing.T) {
	ctrl, mockService, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /counter/name/1", func(t *testing.T) {

		mockService.
			EXPECT().
			UpdateCounter(gomock.Eq(models.Counter{
				Name:  "name",
				Value: 1,
			})).
			Return(&models.Counter{
				Name:  "name",
				Value: 1,
			}, nil)

		res := testingPostURL(t, ts.URL+"/counter/name/1")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestInvalidUpdateCounterHandler(t *testing.T) {
	ctrl, _, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	invalidRequests := []string{
		"/counter/name/1.000",
		"/counter/name/1.25",
		"/counter/name/99999999999999999999999999",
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
	ctrl, _, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /counter/", func(t *testing.T) {
		res := testingPostURL(t, ts.URL+"/counter/")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// tests for update gauges handle
func TestValidUpdateGaugeHandler(t *testing.T) {
	ctrl, mockService, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /gauge/name/1.0", func(t *testing.T) {

		mockService.
			EXPECT().
			UpdateGauge(gomock.Eq(models.Gauge{
				Name:  "name",
				Value: 1.0,
			})).
			Return(&models.Gauge{
				Name:  "name",
				Value: 1.0,
			}, nil)

		res := testingPostURL(t, ts.URL+"/gauge/name/1.0")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestInvalidUpdateGaugeHandler(t *testing.T) {
	ctrl, _, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	invalidRequests := []string{
		"/gauge/name/value",
		"/gauge/name/1.2.3",
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
	ctrl, _, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /gauge/", func(t *testing.T) {
		res := testingPostURL(t, ts.URL+"/gauge/")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// tests for JSON handlers
func TestValidUpdateCounterJSONHandler(t *testing.T) {
	ctrl, mockService, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("update counter: { name, 1 }", func(t *testing.T) {
		counterJSON := `{"id":"name","type":"counter","delta":1}`

		value := models.CounterValue(1)
		counter := models.Metrics{
			MType: models.MetricTypeCounter,
			ID:    "name",
			Delta: &value,
		}

		mockService.
			EXPECT().
			UpdateMetric(counter).
			Return(&counter, nil)

		res := testingPostJSON(t, ts.URL+"/", counterJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, counterJSON, string(body))
	})
}

func TestValidUpdateGaugeJSONHandler(t *testing.T) {
	ctrl, mockService, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("update gauge: { name, 1.5 }", func(t *testing.T) {
		gaugeJSON := `{"id":"name","type":"gauge","value":1.5}`

		value := models.GaugeValue(1.5)
		gauge := models.Metrics{
			MType: models.MetricTypeGauge,
			ID:    "name",
			Value: &value,
		}

		mockService.
			EXPECT().
			UpdateMetric(gauge).
			Return(&gauge, nil)

		res := testingPostJSON(t, ts.URL+"/", gaugeJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, gaugeJSON, string(body))
	})
}

func TestInvalidUpdateJSONHandler(t *testing.T) {
	ctrl, _, ts := getUpdateTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("update invalid", func(t *testing.T) {
		invalidJSON := "{}"

		res := testingPostJSON(t, ts.URL+"/", invalidJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func getUpdateTestObjects(t *testing.T) (*gomock.Controller, *mocks.MockService, *httptest.Server) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	handlers, err := New(mockService, zap.NewNop())
	require.NoError(t, err, "handlers initialization error")

	router := chi.NewRouter()
	handlers.registerUpdateRoutes(context.TODO(), router)

	ts := httptest.NewServer(router)

	return ctrl, mockService, ts
}
