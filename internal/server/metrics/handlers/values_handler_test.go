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

func TestValuesHandler(t *testing.T) {
	ctrl, mockService, ts := getValuesTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /", func(t *testing.T) {

		mockService.
			EXPECT().
			ListGauges(gomock.Any()).
			Return(models.GaugesList{}, nil)

		mockService.
			EXPECT().
			ListCounters(gomock.Any()).
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

func getValuesTestObjects(t *testing.T) (*gomock.Controller, *mocks.MockService, *httptest.Server) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	handlers, err := New(mockService, zap.NewNop())
	require.NoError(t, err, "handlers initialization error")

	router := chi.NewRouter()
	handlers.registerValuesRoutes(context.TODO(), router)

	ts := httptest.NewServer(router)

	return ctrl, mockService, ts
}
