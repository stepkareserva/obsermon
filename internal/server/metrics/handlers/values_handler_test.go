package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/middleware"
	"github.com/stepkareserva/obsermon/internal/server/mocks"
)

func TestValuesHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valuesHandler, err := ValuesHandler(mockService)
	require.NoError(t, err, "value handler initialization error")

	ts := httptest.NewServer(valuesHandler)
	defer ts.Close()

	t.Run("test /", func(t *testing.T) {

		mockService.
			EXPECT().
			ListGauges().
			Return(models.GaugesList{}, nil)

		mockService.
			EXPECT().
			ListCounters().
			Return(models.CountersList{}, nil)

		// get values
		res := testingGetURL(t, ts.URL+"/")
		defer res.Body.Close()

		// check status and contentType if status is ok
		assert.Equal(t, http.StatusOK, res.StatusCode)
		contentType := res.Header.Get("Content-Type")
		_, err = io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "text/html", contentType)
	})
}

func TestValuesCompressedHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	valuesHandler, err := ValuesHandler(mockService)
	require.NoError(t, err, "value handler initialization error")
	valuesHandler = middleware.Compression()(valuesHandler)

	ts := httptest.NewServer(valuesHandler)
	defer ts.Close()

	t.Run("test /", func(t *testing.T) {

		mockService.
			EXPECT().
			ListGauges().
			Return(models.GaugesList{}, nil)

		mockService.
			EXPECT().
			ListCounters().
			Return(models.CountersList{}, nil)

		// get values
		res := testingGetGzipURL(t, ts.URL+"/")
		require.NoError(t, err)
		defer res.Body.Close()

		// check status and contentType if status is o
		_ = testingUngzipBody(t, res)
		contentType := res.Header.Get("Content-Type")
		assert.Equal(t, "text/html", contentType)
	})
}
