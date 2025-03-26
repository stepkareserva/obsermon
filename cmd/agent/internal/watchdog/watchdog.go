package watchdog

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/stepkareserva/obsermon/cmd/agent/internal/client"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/metrics"
	"github.com/stepkareserva/obsermon/cmd/agent/internal/monitor"
)

type WatchdogParams struct {
	MetricsServerClient *client.MetricsClient
	PollInterval        time.Duration
	ReportInterval      time.Duration
}

type Watchdog struct {
	params WatchdogParams

	metrics chan metrics.Metrics

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWatchdog(params WatchdogParams) (*Watchdog, error) {
	if params.PollInterval <= 0 {
		return nil, fmt.Errorf("invalid metrics poll interval")
	}
	if params.ReportInterval <= 0 {
		return nil, fmt.Errorf("invalid metrics report interval")
	}

	chanCapacity := params.ReportInterval/params.PollInterval + 1
	metrics := make(chan metrics.Metrics, chanCapacity)

	ctx, cancel := context.WithCancel(context.Background())

	watchdog := Watchdog{
		params:  params,
		metrics: metrics,
		ctx:     ctx,
		cancel:  cancel,
	}

	return &watchdog, nil
}

func (w *Watchdog) Start() {
	w.wg.Add(2)

	go func() {
		defer w.wg.Done()
		w.metricsPoller()
	}()

	go func() {
		defer w.wg.Done()
		w.metricsReporter()
	}()
}

func (w *Watchdog) Stop() {
	w.cancel()
	w.wg.Wait()
}

func (w *Watchdog) metricsPoller() {
	for {
		timer := time.NewTimer(w.params.PollInterval)
		defer timer.Stop()

		select {
		case <-w.ctx.Done():
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

func (w *Watchdog) metricsReporter() {
	for {
		timer := time.NewTimer(w.params.ReportInterval)
		defer timer.Stop()

		select {
		case <-w.ctx.Done():
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
	metrics := metrics.NewMetrics()

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
