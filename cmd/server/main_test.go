package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	go main()

	// waiting for server start
	time.Sleep(1 * time.Second)

	type PostTestItem struct {
		URL      string
		Expected int
	}

	updateURLPrefix := "http://localhost:8080/update"

	testItems := []PostTestItem{
		// correct
		{URL: "gauge/name/1.0", Expected: http.StatusOK},
		{URL: "counter/name/1", Expected: http.StatusOK},

		// without metric name
		{URL: "gauge/", Expected: http.StatusNotFound},
		{URL: "gauge", Expected: http.StatusNotFound},
		{URL: "counter/", Expected: http.StatusNotFound},
		{URL: "counter", Expected: http.StatusNotFound},

		// incorrect metric type
		{URL: "metric/name/1.0", Expected: http.StatusBadRequest},

		// incorrect metric value
		{URL: "gauge/name/value", Expected: http.StatusBadRequest},
		{URL: "counter/name/1.25", Expected: http.StatusBadRequest},
		{URL: "counter/name/999999999999999999999", Expected: http.StatusBadRequest},
	}

	for _, item := range testItems {
		t.Run(item.URL, func(t *testing.T) {
			url, err := url.JoinPath(updateURLPrefix, item.URL)
			if err != nil {
				t.Errorf("could not join url parts: %v", err)
			}
			resp, err := postText(url)
			if err != nil {
				t.Errorf("post to %s returns error %v", item.URL, err)
			}
			if resp != item.Expected {
				t.Errorf("post to %s returns %d, expected %d", item.URL, resp, item.Expected)
			}
		})
	}
}

func postText(url string) (int, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return 0, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
