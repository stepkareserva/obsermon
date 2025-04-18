package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// tests for update counters handle
func TestValidUpdateCounterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := UpdateHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	t.Run("test /counter/name/1", func(t *testing.T) {

		mockService.
			EXPECT().
			UpdateCounter(gomock.Eq(models.Counter{
				Name:  "name",
				Value: 1,
			})).
			Return(nil)

		res, err := http.Post(ts.URL+"/counter/name/1", "text/plain", nil)
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestInvalidUpdateCounterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := UpdateHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	invalidRequests := []string{
		"/counter/name/1.000",
		"/counter/name/1.25",
		"/counter/name/99999999999999999999999999",
	}
	for _, req := range invalidRequests {
		t.Run(fmt.Sprintf("test %s", req), func(t *testing.T) {
			res, err := http.Post(ts.URL+req, "text/plain", nil)
			require.NoError(t, err)
			defer res.Body.Close()
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestNotFoundUpdateCounterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := UpdateHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	t.Run("test /counter/", func(t *testing.T) {
		res, err := http.Post(ts.URL+"/counter/", "text/plain", nil)
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

// tests for update gauges handle
func TestValidUpdateGaugeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := UpdateHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	t.Run("test /gauge/name/1.0", func(t *testing.T) {

		mockService.
			EXPECT().
			UpdateGauge(gomock.Eq(models.Gauge{
				Name:  "name",
				Value: 1.0,
			})).
			Return(nil)

		res, err := http.Post(ts.URL+"/gauge/name/1.0", "text/plain", nil)
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestInvalidUpdateGaugeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := UpdateHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	invalidRequests := []string{
		"/gauge/name/value",
		"/gauge/name/1.2.3",
	}
	for _, req := range invalidRequests {
		t.Run(fmt.Sprintf("test %s", req), func(t *testing.T) {
			res, err := http.Post(ts.URL+req, "text/plain", nil)
			require.NoError(t, err)
			defer res.Body.Close()
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestNotFoundUpdateGaugeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	updateHandler, err := UpdateHandler(mockService)
	require.NoError(t, err, "update handler initialization error")

	ts := httptest.NewServer(updateHandler)
	defer ts.Close()

	t.Run("test /gauge/", func(t *testing.T) {
		res, err := http.Post(ts.URL+"/gauge/", "text/plain", nil)
		require.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
