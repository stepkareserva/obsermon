package router

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPingHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /ping", func(t *testing.T) {

		mockService.
			EXPECT().
			Ping(gomock.Any()).
			Return(nil)

		// get values
		res := testingGetURL(t, ts.URL+"/ping")
		defer safeCloseRes(t, res)

		// check status and contentType if status is ok
		assert.Equal(t, http.StatusOK, res.StatusCode)
		_, err := io.ReadAll(res.Body)
		require.NoError(t, err)
	})
}

func TestInvalidPingHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
	defer ctrl.Finish()
	defer ts.Close()

	t.Run("test /ping", func(t *testing.T) {

		mockService.
			EXPECT().
			Ping(gomock.Any()).
			Return(errors.New("database unavailable"))

		// get values
		res := testingGetURL(t, ts.URL+"/ping")
		defer safeCloseRes(t, res)

		// check status and contentType if status is ok
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		_, err := io.ReadAll(res.Body)
		require.NoError(t, err)
	})
}
