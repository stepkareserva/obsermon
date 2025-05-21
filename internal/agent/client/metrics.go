package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stepkareserva/obsermon/internal/models"
)

type MetricsClient struct {
	client    *resty.Client
	secretkey string
}

const requestTimeout = 5 * time.Second

// header HashSHA256 is forbidden by checker
var signHeader = http.CanonicalHeaderKey("HashSHA256")

func New(endpoint string, secretkey string) (*MetricsClient, error) {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid endpoint scheme %s", u.Scheme)
	}

	client := resty.New()
	client.SetBaseURL(endpoint)
	client.SetTimeout(requestTimeout)

	return &MetricsClient{client: client, secretkey: secretkey}, nil
}

func (c *MetricsClient) UpdateCounter(value models.Counter) error {
	return c.BatchUpdate(models.CountersList{value}, nil)
}

func (c *MetricsClient) UpdateGauge(value models.Gauge) error {
	return c.BatchUpdate(nil, models.GaugesList{value})
}

func (c *MetricsClient) BatchUpdate(counters models.CountersList, gauges models.GaugesList) error {
	metrics := make(models.Metrics, 0, len(counters)+len(gauges))
	for _, counter := range counters {
		metrics = append(metrics, models.CounterMetric(counter))
	}
	for _, gauge := range gauges {
		metrics = append(metrics, models.GaugeMetric(gauge))
	}
	return c.sendUpdateRequest(metrics)
}

func (c *MetricsClient) sendUpdateRequest(metrics models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	attemptsIntervals := []time.Duration{
		0 * time.Second,
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	var resp *resty.Response
	var err error
	for _, waitIterval := range attemptsIntervals {
		time.Sleep(waitIterval)

		resp, err = c.postJSON("/updates", metrics)

		switch {
		case err == nil:
			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf("post %s request status %d",
					resp.Request.URL, resp.StatusCode())
			}
			return nil
		case !isServerUnavailableErr(err):
			return fmt.Errorf("post updates: %w", err)
		}
	}

	return fmt.Errorf("post updates: %w", err)
}

func (c *MetricsClient) postJSON(url string, object interface{}) (*resty.Response, error) {
	body, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("json marshalling body: %w", err)
	}

	// !important: we must post only raw data as io.Reader,
	// all other variants like SetBody(object) or SetBody(body)
	// doesn't guarantee that request bodt will be equal to 'body'
	req := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewReader(body))

	if len(c.secretkey) > 0 {
		h := hmac.New(sha256.New, []byte(c.secretkey))
		h.Write(body)
		hashSum := hex.EncodeToString(h.Sum(nil))
		req.SetHeader(signHeader, hashSum)
	}

	return req.Post(url)
}

func isServerUnavailableErr(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		// too long for request timeout
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	return false
}
