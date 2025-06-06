package service

import (
	"context"
	"testing"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGaugeService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	service, err := New(mockStorage)
	require.NoError(t, err, "service initialization error")

	t.Run("test gauge", func(t *testing.T) {
		mockStorage.
			EXPECT().
			SetGauge(context.TODO(), models.Gauge{
				Name:  "name",
				Value: 1.0,
			})

		mockStorage.
			EXPECT().
			FindGauge(context.TODO(), "name").
			Return(&models.Gauge{
				Name:  "name",
				Value: 1.0,
			}, true, nil)

		_, err := service.UpdateGauge(context.TODO(), models.Gauge{Name: "name", Value: 1.0})
		assert.NoError(t, err)
		gauge, exists, err := service.FindGauge(context.TODO(), "name")
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, gauge, &models.Gauge{Name: "name", Value: 1.0})
	})
}

func TestCounterService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	service, err := New(mockStorage)
	require.NoError(t, err, "service initialization error")

	t.Run("test counter", func(t *testing.T) {

		mockStorage.
			EXPECT().
			UpdateCounter(context.TODO(), models.Counter{
				Name:  "name",
				Value: 1,
			})

		mockStorage.
			EXPECT().
			UpdateCounter(context.TODO(), models.Counter{
				Name:  "name",
				Value: 2,
			})

		mockStorage.
			EXPECT().
			FindCounter(context.TODO(), "name").
			Return(&models.Counter{
				Name:  "name",
				Value: 3,
			}, true, nil)

		_, err := service.UpdateCounter(context.TODO(), models.Counter{Name: "name", Value: 1})
		assert.NoError(t, err)
		_, err = service.UpdateCounter(context.TODO(), models.Counter{Name: "name", Value: 2})
		assert.NoError(t, err)

		counter, exists, err := service.FindCounter(context.TODO(), "name")
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, counter, &models.Counter{Name: "name", Value: 3})
	})
}
