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

func TestValuesHandler(t *testing.T) {
	ctrl, mockService, ts := getTestObjects(t)
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
