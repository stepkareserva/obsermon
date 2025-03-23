package watchdog

import (
	"context"
	"errors"
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
	UpdateInterval      time.Duration
}

type Watchdog struct {
	params WatchdogParams

	metrics chan metrics.Metrics

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewWatchdog(params WatchdogParams) *Watchdog {
	chanCapacity := params.UpdateInterval/params.PollInterval + 1
	metrics := make(chan metrics.Metrics, chanCapacity)

	ctx, cancel := context.WithCancel(context.Background())

	watchdog := Watchdog{
		params:  params,
		metrics: metrics,
		ctx:     ctx,
		cancel:  cancel,
	}

	return &watchdog
}

func (w *Watchdog) Start() {
	w.wg.Add(2)

	go func() {
		defer w.wg.Done()
		w.metricsPoller()
	}()

	go func() {
		defer w.wg.Done()
		w.metricsUpdater()
	}()
}

func (w *Watchdog) Stop() {
	w.cancel()
	w.wg.Wait()
}

func (w *Watchdog) metricsPoller() {
	for {
		time.Sleep(w.params.PollInterval)
		select {
		case <-w.ctx.Done():
			close(w.metrics)
			log.Println("Metrics updater stopped")
			return
		default:
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

func (w *Watchdog) metricsUpdater() {
	for {
		time.Sleep(w.params.UpdateInterval)
		select {
		case <-w.ctx.Done():
			log.Println("Metrics updater stopped")
			return
		default:
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
	for k, v := range metrics.Counters {
		if err := w.params.MetricsServerClient.UpdateCounter(k, v); err != nil {
			errs = append(errs, err)
		}
	}
	for k, v := range metrics.Gauges {
		if err := w.params.MetricsServerClient.UpdateGauge(k, v); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
