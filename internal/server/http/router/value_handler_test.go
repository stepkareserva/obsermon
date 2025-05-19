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
func TestValidValueCounterHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /value/counter/name", func(t *testing.T) {

		mockService.
			EXPECT().
			FindCounter(gomock.Any(), gomock.Eq("name")).
			Return(&models.Counter{
				Name:  "name",
				Value: 1,
			}, true, nil)

		res := testingGetURL(t, ts.URL+"/value/counter/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "1", string(body))
	})
}

func TestNotFoundValueCounterHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /value/counter/name", func(t *testing.T) {

		mockService.
			EXPECT().
			FindCounter(gomock.Any(), gomock.Eq("name")).
			Return(nil, false, nil)

		res := testingGetURL(t, ts.URL+"/value/counter/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// test for gauge value handler
func TestValidValueGaugeHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /value/gauge/name", func(t *testing.T) {

		mockService.
			EXPECT().
			FindGauge(gomock.Any(), gomock.Eq("name")).
			Return(&models.Gauge{
				Name:  "name",
				Value: 1.2,
			}, true, nil)

		res := testingGetURL(t, ts.URL+"/value/gauge/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "1.2", string(body))
	})
}

func TestNotFoundValueGaugeHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /value/gauge/name", func(t *testing.T) {

		mockService.
			EXPECT().
			FindGauge(gomock.Any(), gomock.Eq("name")).
			Return(nil, false, nil)

		res := testingGetURL(t, ts.URL+"/value/gauge/name")
		defer safeCloseRes(t, res)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestValidValueCounterJSONHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("value counter: { name }", func(t *testing.T) {
		counterJSON := `{"id":"name","type":"counter","delta":1}`

		value := models.CounterValue(1)
		counter := models.Metric{
			MType: models.MetricTypeCounter,
			ID:    "name",
			Delta: &value,
		}

		mockService.
			EXPECT().
			FindMetric(gomock.Any(), models.MetricTypeCounter, "name").
			Return(&counter, true, nil)

		res := testingPostJSON(t, ts.URL+"/value", counterJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, counterJSON, string(body))
	})
}

func TestValidValueGaugeJSONHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("value gauge: { name }", func(t *testing.T) {
		gaugeJSON := `{"id":"name","type":"gauge","value":1.5}`

		value := models.GaugeValue(1.5)
		gauge := models.Metric{
			MType: models.MetricTypeGauge,
			ID:    "name",
			Value: &value,
		}

		mockService.
			EXPECT().
			FindMetric(gomock.Any(), models.MetricTypeGauge, "name").
			Return(&gauge, true, nil)

		res := testingPostJSON(t, ts.URL+"/value", gaugeJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.JSONEq(t, gaugeJSON, string(body))
	})
}

func TestInvalidValueJSONHandler(t *testing.T) {
	ctrl, _, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("value invalid: {}", func(t *testing.T) {
		invalidJSON := "{}"

		res := testingPostJSON(t, ts.URL+"/value", invalidJSON)
		defer safeCloseRes(t, res)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
