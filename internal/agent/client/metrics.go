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

func NewMetricsClient(endpoint string) (*MetricsClient, error) {
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
	pathParams := map[string]string{
		URLMetric: MetricCounter,
		URLName:   value.Name,
		URLValue:  value.Value.String(),
	}

	return c.sendUpdateRequest(pathParams)
}

func (c *MetricsClient) UpdateGauge(value models.Gauge) error {
	pathParams := map[string]string{
		URLMetric: MetricGauge,
		URLName:   value.Name,
		URLValue:  value.Value.String(),
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
