package models

import (
	"fmt"
	"strconv"
	"strings"
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
		return fmt.Errorf("could not parse GaugeValue: %v", err)
	}
	*g = GaugeValue(value)
	return nil
}

func (g *GaugeValue) String() string {
	s := strconv.FormatFloat(float64(*g), 'f', -1, 64)
	// solution as good as challenge
	if !strings.Contains(s, ".") {
		s += "."
	}
	return s
}

// ? maybe not here?
func (g *GaugeValue) PrettyString() string {
	return strconv.FormatFloat(float64(*g), 'g', 6, 64)
}
