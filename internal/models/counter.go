package models

import (
	"fmt"
	"math"
	"strconv"
)

type CounterValue int64

type Counter struct {
	Name  string
	Value CounterValue
}

// Q: maybe implement encoding.TextMarshaler/Unmarshaler?
func (c *CounterValue) FromString(s string) error {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse CounterValue: %v", err)
	}
	*c = CounterValue(value)
	return nil
}

func (c *CounterValue) String() string {
	return strconv.FormatInt(int64(*c), 10)
}

func (c *CounterValue) Update(v CounterValue) error {
	updated, err := safeAdd(*c, v)
	if err != nil {
		return err
	}
	*c = updated
	return nil
}

// stuff for counter overflow handling
const (
	counterMax CounterValue = CounterValue(math.MaxInt64)
	counterMin CounterValue = CounterValue(math.MinInt64)
)

type CounterOverflowError struct {
	a, b CounterValue
}

func (e CounterOverflowError) Error() string {
	return fmt.Sprintf("CounterOverflowError: caused by %d + %d", e.a, e.b)
}

func safeAdd(a, b CounterValue) (CounterValue, error) {
	if (b > 0 && a > counterMax-b) || (b < 0 && a < counterMin-b) {
		return 0, CounterOverflowError{a: a, b: b}
	}
	return a + b, nil
}
