package watchdog

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/stepkareserva/obsermon/internal/agent/client"
	"github.com/stepkareserva/obsermon/internal/agent/metrics"
	"github.com/stepkareserva/obsermon/internal/agent/monitor"
)

type WatchdogParams struct {
	MetricsServerClient *client.MetricsClient
	PollInterval        time.Duration
	ReportInterval      time.Duration
}

type Watchdog struct {
	params WatchdogParams

	metrics chan metrics.Metrics
}

func New(params WatchdogParams) (*Watchdog, error) {
	if params.PollInterval <= 0 {
		return nil, fmt.Errorf("invalid metrics poll interval")
	}
	if params.ReportInterval <= 0 {
		return nil, fmt.Errorf("invalid metrics report interval")
	}

	chanCapacity := params.ReportInterval/params.PollInterval + 1
	metrics := make(chan metrics.Metrics, chanCapacity)

	watchdog := Watchdog{
		params:  params,
		metrics: metrics,
	}

	return &watchdog, nil
}

func (w *Watchdog) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		w.runtimeMetricsPoller(ctx)
	}()

	go func() {
		defer wg.Done()
		w.golangMetricsPoller(ctx)
	}()

	go func() {
		defer wg.Done()
		w.metricsReporter(ctx)
	}()

	wg.Wait()
}

func (w *Watchdog) runtimeMetricsPoller(ctx context.Context) {
	for {
		timer := time.NewTimer(w.params.PollInterval)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			close(w.metrics)
			log.Println("Metrics updater stopped")
			return
		case <-timer.C:
			m, err := monitor.GetRuntimeMetrics()
			if err != nil {
				// just skip because of what else shall we do
				log.Printf("Get Runtime Metrics error: %v", err)
			} else {
				w.metrics <- *m
			}
		}
	}
}

func (w *Watchdog) golangMetricsPoller(ctx context.Context) {
	for {
		timer := time.NewTimer(w.params.PollInterval)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			close(w.metrics)
			log.Println("Metrics updater stopped")
			return
		case <-timer.C:
			m, err := monitor.GetGolangMetrics()
			if err != nil {
				// just skip because of what else shall we do
				log.Printf("Get Golang Metrics error: %v", err)
			} else {
				w.metrics <- *m
			}
		}
	}
}

func (w *Watchdog) metricsReporter(ctx context.Context) {
	for {
		timer := time.NewTimer(w.params.ReportInterval)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			log.Println("Metrics reporter stopped")
			return
		case <-timer.C:
			metrics := w.processPolledMetrics()
			w.sendMetrics(metrics)
		}
	}
}

func (w *Watchdog) processPolledMetrics() metrics.Metrics {
	metrics := metrics.New()

	for {
		select {
		case m := <-w.metrics:
			if err := metrics.Update(m); err != nil {
				// just skip because of what else shall we do
				log.Printf("metrics.Update error: %v", err)
			}
		default:
			return metrics
		}
	}
}

func (w *Watchdog) sendMetrics(metrics metrics.Metrics) {
	w.params.MetricsServerClient.BatchUpdate(
		metrics.Counters.List(),
		metrics.Gauges.List())
}
