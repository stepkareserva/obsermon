package monitor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRuntimeGauges(t *testing.T) {
	_, err := getRuntimeGauges()
	require.NoError(t, err)
}

func TestGetRuntimeMetrics(t *testing.T) {
	metrics, err := GetRuntimeMetrics()
	require.NoError(t, err)
	require.Equal(t, 28, len(metrics.Gauges))
	require.Equal(t, 1, len(metrics.Counters))
}
