package models

import (
	"fmt"
	"strconv"
)

type GaugeValue float64

type Gauge struct {
	Name  string
	Value GaugeValue
}

// Q: maybe implement encoding.TextMarshaler/Unmarshaler?
func (g *GaugeValue) FromString(s string) error {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("could not parse GaugeValue: %w", err)
	}
	*g = GaugeValue(value)
	return nil
}

func (g *GaugeValue) String() string {
	return strconv.FormatFloat(float64(*g), 'f', -1, 64)
}

// ? maybe not here?
func (g *GaugeValue) PrettyString() string {
	return strconv.FormatFloat(float64(*g), 'g', 6, 64)
}
