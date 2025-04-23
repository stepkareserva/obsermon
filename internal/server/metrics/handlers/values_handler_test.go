package handlers

import (
	"context"
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

func TestValuesHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)

	ts := testValuesServer(t, mockService)
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
		defer safeCloseRes(t, res)

		// check status and contentType if status is ok
		assert.Equal(t, http.StatusOK, res.StatusCode)
		contentType := res.Header.Get("Content-Type")
		_, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, "text/html", contentType)
	})
}

func testValuesServer(t *testing.T, mockService *mocks.MockService) *httptest.Server {
	metricsHandler, err := New(mockService, zap.NewNop())
	require.NoError(t, err, "metrics handler initialization error")
	valuesHandler, err := metricsHandler.valuesHandler(context.Background())
	require.NoError(t, err, "value handler initialization error")
	return httptest.NewServer(valuesHandler)
}
