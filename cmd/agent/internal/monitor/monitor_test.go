package monitor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRuntimeGauges(t *testing.T) {
	_, err := getRuntimeGauges()
	require.NoError(t, err)
}

func TestGetMetrics(t *testing.T) {
	metrics, err := GetMetrics()
	require.NoError(t, err)
	require.Equal(t, 28, len(metrics.Gauges))
	require.Equal(t, 1, len(metrics.Counters))
}
