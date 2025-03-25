package watchdog

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stepkareserva/obsermon/cmd/agent/internal/client"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/monitor"
	"github.com/stretchr/testify/require"
)

func TestWatchdog(t *testing.T) {
	// test params
	pollInterval := 300 * time.Millisecond
	reportInterval := 1 * time.Second
	expectedPollCount := 3
	expectedURLPath := fmt.Sprintf("/update/counter/%s/%d", monitor.PollCount, expectedPollCount)

	// mock server just collect all incoming requests
	incomingRequests := make(map[string]struct{})
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingRequests[r.URL.Path] = struct{}{}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// mock server client
	metricsClient, err := client.NewMetricsClient(mockServer.URL)
	require.NoError(t, err)

	// watchdog
	watchdogParams := WatchdogParams{
		PollInterval:        pollInterval,
		ReportInterval:      reportInterval,
		MetricsServerClient: metricsClient,
	}
	watchdog := NewWatchdog(watchdogParams)
	watchdog.Start()
	defer watchdog.Stop()

	// wait reportInterval + 100 ms
	time.Sleep(reportInterval + 100*time.Millisecond)

	// check target requests on mock server
	_, exists := incomingRequests[expectedURLPath]
	require.True(t, exists)
}
