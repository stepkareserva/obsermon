package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// test for counter value handler
func TestValidValueCounterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := ValueHandler(mockService)
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

		res, err := http.Get(ts.URL + "/counter/name")
		require.NoError(t, err)
		defer res.Body.Close()
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
	valueHandler, err := ValueHandler(mockService)
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(valueHandler)
	defer ts.Close()

	t.Run("test /counter/name", func(t *testing.T) {

		mockService.
			EXPECT().
			GetCounter(gomock.Eq("name")).
			Return(nil, false, nil)

		res, err := http.Get(ts.URL + "/counter/name")
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// test for gauge value handler
func TestValidValueGaugeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valueHandler, err := ValueHandler(mockService)
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

		res, err := http.Get(ts.URL + "/gauge/name")
		require.NoError(t, err)
		defer res.Body.Close()
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
	updateHandler, err := ValueHandler(mockService)
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	t.Run("test /gauge/name", func(t *testing.T) {

		mockService.
			EXPECT().
			GetGauge(gomock.Eq("name")).
			Return(nil, false, nil)

		res, err := http.Get(ts.URL + "/gauge/name")
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// tests for json handlers

// tests for JSON handlers
/*func TestValidUpdateCounterJSONHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := UpdateHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(updateHandler)
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
			Return(nil)
		mockService.
			EXPECT().
			GetMetric(models.MetricTypeCounter, "name").
			Return(&counter, true, nil)

		res, err := http.Post(ts.URL+"/", "application/json", strings.NewReader(counterJSON))
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, counterJSON, string(body))
	})
}*/

func TestValidValueCounterJSONHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valueHandler, err := ValueHandler(mockService)
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

		res, err := http.Post(ts.URL+"/", "application/json", strings.NewReader(counterJSON))
		require.NoError(t, err)
		defer res.Body.Close()
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
	valueHandler, err := ValueHandler(mockService)
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

		res, err := http.Post(ts.URL+"/", "application/json", strings.NewReader(gaugeJSON))
		require.NoError(t, err)
		defer res.Body.Close()
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
	valueHandler, err := ValueHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(valueHandler)
	defer ts.Close()

	t.Run("value invalid: {}", func(t *testing.T) {
		invalidJSON := "{}"

		res, err := http.Post(ts.URL+"/", "application/json", strings.NewReader(invalidJSON))
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
