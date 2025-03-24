package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/stepkareserva/obsermon/internal/models"
)

const (
	// metrics for url
	MetricGauge   = "gauge"
	MetricCounter = "counter"

	// names of chi routing url params to be extracted
	URLMetric = "metric"
	URLName   = "name"
	URLValue  = "value"
)

type MetricsClient struct {
	client *resty.Client
}

func NewMetricsClient(s string) (*MetricsClient, error) {
	parsedEndpoint, err := url.ParseRequestURI(s)
	if err != nil || parsedEndpoint.Scheme == "" || parsedEndpoint.Host == "" {
		return nil, fmt.Errorf("invalid server url")
	}

	client := resty.New()
	client.SetBaseURL(s)

	return &MetricsClient{client: client}, nil
}

func (c *MetricsClient) UpdateCounter(name string, value models.Counter) error {
	pathParams := map[string]string{
		URLMetric: MetricCounter,
		URLName:   name,
		URLValue:  value.ToString(),
	}

	return c.sendUpdateRequest(pathParams)
}

func (c *MetricsClient) UpdateGauge(name string, value models.Gauge) error {
	pathParams := map[string]string{
		URLMetric: MetricGauge,
		URLName:   name,
		URLValue:  value.ToString(),
	}

	return c.sendUpdateRequest(pathParams)
}

func (c *MetricsClient) sendUpdateRequest(pathParams map[string]string) error {
	resp, err := c.client.R().
		SetPathParams(pathParams).
		SetHeader("Content-Type", "text/plain").
		SetBody(http.NoBody).
		Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", URLMetric, URLName, URLValue))

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("post %s request error, status %d", resp.Request.URL, resp.StatusCode())
	}
	return nil
}
