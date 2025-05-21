package watchdog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stepkareserva/obsermon/internal/agent/client"
	"github.com/stretchr/testify/require"
)

func TestWatchdog(t *testing.T) {
	// test params
	pollInterval := 300 * time.Millisecond
	reportInterval := 1 * time.Second
	expectedURLPath := "/updates"

	// mock server just collect all incoming requests
	incomingRequests := make(map[string]struct{})
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingRequests[r.URL.Path] = struct{}{}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// mock server client
	metricsClient, err := client.New(mockServer.URL, "")
	require.NoError(t, err)

	// watchdog
	watchdogParams := WatchdogParams{
		PollInterval:        pollInterval,
		ReportInterval:      reportInterval,
		MetricsServerClient: metricsClient,
	}

	runningTime := reportInterval + 100*time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), runningTime)
	defer cancel()

	watchdog, err := New(watchdogParams)
	require.NoError(t, err)
	watchdog.Start(ctx)

	// check target requests on mock server
	_, exists := incomingRequests[expectedURLPath]
	require.True(t, exists)
}
