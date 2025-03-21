package models

import (
	"fmt"
	"strings"
)

type UpdateCounterRequest struct {
	Name  string
	Value Counter
}

type UpdateGaugeRequest struct {
	Name  string
	Value Gauge
}

// marshal and unmarshal counter and gauge update requests
// onto url path "name/value"

// ну, какой формат - такой и парсинг :/
// может, не знаю, на generic-ах сделать?

func (r *UpdateCounterRequest) ToURLPath() (string, error) {
	return metricToURLPath(r.Name, r.Value.ToString())
}

func (r *UpdateCounterRequest) FromURLPath(s string) error {
	name, value, err := metricFromURLPath(s)
	if err != nil {
		return err
	}
	var counter Counter
	if err := counter.FromString(value); err != nil {
		return &InvalidMetricValueError{}
	}
	r.Name = name
	r.Value = counter
	return nil
}

func (r *UpdateGaugeRequest) ToURLPath() (string, error) {
	return metricToURLPath(r.Name, r.Value.ToString())
}

func (r *UpdateGaugeRequest) FromURLPath(s string) error {
	name, value, err := metricFromURLPath(s)
	if err != nil {
		return err
	}
	var gauge Gauge
	if err := gauge.FromString(value); err != nil {
		return &InvalidMetricValueError{}
	}
	r.Name = name
	r.Value = gauge
	return nil
}

func metricToURLPath(name, value string) (string, error) {
	if strings.Contains(name, "/") {
		return "", &InvalidMetricNameError{}
	}
	return fmt.Sprintf("%s/%s", name, value), nil
}

func metricFromURLPath(s string) (name, value string, err error) {
	components := strings.Split(s, "/")
	if len(s) == 0 || len(components) == 0 {
		return "", "", &MissingMetricNameError{}
	}
	if len(components) == 1 {
		return "", "", &InvalidMetricValueError{}
	}
	if len(components) > 2 {
		return "", "", &ExtraRequestComponentsError{}
	}

	return components[0], components[1], nil
}
