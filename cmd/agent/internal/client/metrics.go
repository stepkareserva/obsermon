package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/stepkareserva/obsermon/internal/models"
)

type MetricsClient struct {
	url string
}

func NewMetricsClient(s string) (*MetricsClient, error) {
	parsedURL, err := url.ParseRequestURI(s)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid server url")
	}
	return &MetricsClient{url: s}, nil
}

func (c *MetricsClient) UpdateCounter(name string, value models.Counter) error {
	request := models.UpdateCounterRequest{
		Name:  name,
		Value: value,
	}
	requestURL, err := request.ToURLPath()
	if err != nil {
		return err
	}

	path, err := url.JoinPath(c.url, "update", "counter", requestURL)
	if err != nil {
		return err
	}
	return c.sendPost(path)
}

func (c *MetricsClient) UpdateGauge(name string, value models.Gauge) error {
	request := models.UpdateGaugeRequest{
		Name:  name,
		Value: value,
	}
	requestURL, err := request.ToURLPath()
	if err != nil {
		return err
	}

	path, err := url.JoinPath(c.url, "update", "gauge", requestURL)
	if err != nil {
		return err
	}
	return c.sendPost(path)
}

func (c *MetricsClient) sendPost(url string) error {

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post %s request error, status %d", url, resp.StatusCode)
	}
	return nil
}
