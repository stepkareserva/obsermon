package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/mocks"
)

// test for counter value handler
func TestValidValueCounterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := ValueHandler(mockService, zap.NewNop())
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	t.Run("test /counter/name", func(t *testing.T) {

		mockService.
			EXPECT().
			GetCounter(gomock.Eq("name")).
			Return(&models.Counter{
				Name:  "name",
				Value: 1,
			}, true, nil)

		res := testingGetURL(t, ts.URL+"/counter/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "1", string(body))
	})
}

func TestNotFoundValueCounterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valueHandler, err := ValueHandler(mockService, zap.NewNop())
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(valueHandler)
	defer ts.Close()

	t.Run("test /counter/name", func(t *testing.T) {

		mockService.
			EXPECT().
			GetCounter(gomock.Eq("name")).
			Return(nil, false, nil)

		res := testingGetURL(t, ts.URL+"/counter/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// test for gauge value handler
func TestValidValueGaugeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valueHandler, err := ValueHandler(mockService, zap.NewNop())
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(valueHandler)
	defer ts.Close()

	t.Run("test /gauge/name", func(t *testing.T) {

		mockService.
			EXPECT().
			GetGauge(gomock.Eq("name")).
			Return(&models.Gauge{
				Name:  "name",
				Value: 1.2,
			}, true, nil)

		res := testingGetURL(t, ts.URL+"/gauge/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "1.2", string(body))
	})
}

func TestNotFoundValueGaugeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := ValueHandler(mockService, zap.NewNop())
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	t.Run("test /gauge/name", func(t *testing.T) {

		mockService.
			EXPECT().
			GetGauge(gomock.Eq("name")).
			Return(nil, false, nil)

		res := testingGetURL(t, ts.URL+"/gauge/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestValidValueCounterJSONHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valueHandler, err := ValueHandler(mockService, zap.NewNop())
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(valueHandler)
	defer ts.Close()

	t.Run("value counter: { name }", func(t *testing.T) {
		counterJSON := `{"id":"name","type":"counter","delta":1}`

		value := models.CounterValue(1)
		counter := models.Metrics{
			MType: models.MetricTypeCounter,
			ID:    "name",
			Delta: &value,
		}

		mockService.
			EXPECT().
			GetMetric(models.MetricTypeCounter, "name").
			Return(&counter, true, nil)

		res := testingPostJSON(t, ts.URL+"/", counterJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, counterJSON, string(body))
	})
}

func TestValidValueGaugeJSONHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valueHandler, err := ValueHandler(mockService, zap.NewNop())
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(valueHandler)
	defer ts.Close()

	t.Run("value gauge: { name }", func(t *testing.T) {
		gaugeJSON := `{"id":"name","type":"gauge","value":1.5}`

		value := models.GaugeValue(1.5)
		gauge := models.Metrics{
			MType: models.MetricTypeGauge,
			ID:    "name",
			Value: &value,
		}

		mockService.
			EXPECT().
			GetMetric(models.MetricTypeGauge, "name").
			Return(&gauge, true, nil)

		res := testingPostJSON(t, ts.URL+"/", gaugeJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, gaugeJSON, string(body))
	})
}

func TestInvalidValueJSONHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valueHandler, err := ValueHandler(mockService, zap.NewNop())
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(valueHandler)
	defer ts.Close()

	t.Run("value invalid: {}", func(t *testing.T) {
		invalidJSON := "{}"

		res := testingPostJSON(t, ts.URL+"/", invalidJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
