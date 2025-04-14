package watchdog

import (
	"context"
	"errors"
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
	wg.Add(2)

	go func() {
		defer wg.Done()
		w.metricsPoller(ctx)
	}()

	go func() {
		defer wg.Done()
		w.metricsReporter(ctx)
	}()

	wg.Wait()
}

/*func (w *Watchdog) Stop() {
	w.cancel()
	w.wg.Wait()
}*/

func (w *Watchdog) metricsPoller(ctx context.Context) {
	for {
		timer := time.NewTimer(w.params.PollInterval)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			close(w.metrics)
			log.Println("Metrics updater stopped")
			return
		case <-timer.C:
			m, err := monitor.GetMetrics()
			if err != nil {
				// just skip because of what else shall we do
				log.Printf("Get Metrics error: %v", err)
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
			if err := w.sendMetrics(metrics); err != nil {
				// just skip because of what else shall we do
				log.Printf("Send Metrics error: %v", err)
			}
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

func (w *Watchdog) sendMetrics(metrics metrics.Metrics) error {
	var errs []error
	for _, m := range metrics.Counters.List() {
		if err := w.params.MetricsServerClient.UpdateCounter(m); err != nil {
			errs = append(errs, err)
		}
	}
	for _, m := range metrics.Gauges.List() {
		if err := w.params.MetricsServerClient.UpdateGauge(m); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
