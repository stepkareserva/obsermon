package models

import (
	"fmt"
	"strconv"
)

type Gauge float64

func (g *Gauge) FromString(s string) error {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("could not parse Gauge: %w", err)
	}
	*g = Gauge(value)
	return nil
}

func (g *Gauge) ToString() string {
	return strconv.FormatFloat(float64(*g), 'f', -1, 64)
}
