package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/stepkareserva/obsermon/internal/models"
)

type MetricsClient struct {
	client *resty.Client
}

func New(endpoint string) (*MetricsClient, error) {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid endpoint scheme %s", u.Scheme)
	}

	client := resty.New()
	client.SetBaseURL(endpoint)

	return &MetricsClient{client: client}, nil
}

func (c *MetricsClient) UpdateCounter(value models.Counter) error {
	return c.BatchUpdate(models.CountersList{value}, nil)
}

func (c *MetricsClient) UpdateGauge(value models.Gauge) error {
	return c.BatchUpdate(nil, models.GaugesList{value})
}

func (c *MetricsClient) BatchUpdate(counters models.CountersList, gauges models.GaugesList) error {
	metrics := make([]models.Metrics, 0, len(counters)+len(gauges))
	for _, counter := range counters {
		metrics = append(metrics, models.CounterMetric(counter))
	}
	for _, gauge := range gauges {
		metrics = append(metrics, models.GaugeMetric(gauge))
	}
	return c.sendUpdateRequest(metrics)
}

func (c *MetricsClient) sendUpdateRequest(metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metrics).
		Post("/updates")

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("post %s request error, status %d", resp.Request.URL, resp.StatusCode())
	}
	return nil
}
