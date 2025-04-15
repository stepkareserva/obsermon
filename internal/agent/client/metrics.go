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
	metric := models.CounterMetric(value)
	return c.sendUpdateRequest(metric)
}

func (c *MetricsClient) UpdateGauge(value models.Gauge) error {
	metric := models.GaugeMetric(value)
	return c.sendUpdateRequest(metric)
}

func (c *MetricsClient) sendUpdateRequest(metric models.Metrics) error {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metric).
		Post("/update")

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("post %s request error, status %d", resp.Request.URL, resp.StatusCode())
	}
	return nil
}
