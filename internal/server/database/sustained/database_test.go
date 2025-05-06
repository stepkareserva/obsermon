package sustained

import (
	"context"
	"net"
	"testing"

	"github.com/stepkareserva/obsermon/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type timeoutError struct{}

func (e timeoutError) Error() string   { return "simulated timeout" }
func (e timeoutError) Timeout() bool   { return true }
func (e timeoutError) Temporary() bool { return false }

var _ net.Error = timeoutError{}

func TestUnavailableDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	sustainedDB, err := New(mockDB)
	require.NoError(t, err, "sustained db initialization error")

	ctx := context.TODO()
	query := "RetryableQuery"
	retryableError := &timeoutError{}

	t.Run("test unavailable db", func(t *testing.T) {
		mockDB.
			EXPECT().
			Query(ctx, query).
			Times(4).
			Return(nil, retryableError)

		rows, err := sustainedDB.Query(ctx, query)
		if err == nil && rows != nil {
			defer rows.Close()
		}
		// for suppress false alarm from autotests
		if rows != nil {
			assert.NoError(t, rows.Err())
		}
		assert.ErrorIs(t, err, retryableError)
	})
}

func TestUnstableDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	sustainedDB, err := New(mockDB)
	require.NoError(t, err, "sustained db initialization error")

	ctx := context.TODO()
	query := "RetryableQuery"
	retryableError := &timeoutError{}

	t.Run("test unstable db", func(t *testing.T) {
		mockDB.
			EXPECT().
			Query(ctx, query).
			Times(2).
			Return(nil, retryableError)
		mockDB.
			EXPECT().
			Query(ctx, query).
			Times(1).
			Return(nil, nil)
		rows, err := sustainedDB.Query(ctx, query)
		if err == nil && rows != nil {
			defer rows.Close()
		}
		// for suppress false alarm from autotests
		if rows != nil {
			assert.NoError(t, rows.Err())
		}
		assert.NoError(t, err)
	})
}
