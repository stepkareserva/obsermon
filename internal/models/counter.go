package models

import (
	"fmt"
	"math"
	"strconv"
)

type Counter int64

func (c *Counter) FromString(s string) error {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse Counter: %w", err)
	}
	*c = Counter(value)
	return nil
}

func (c *Counter) ToString() string {
	return strconv.FormatInt(int64(*c), 10)
}

func (c *Counter) Update(v Counter) error {
	updated, err := safeAdd(*c, v)
	if err != nil {
		return err
	}
	*c = updated
	return nil
}

// stuff for counter overflow handling
const (
	counterMax Counter = Counter(math.MaxInt64)
	counterMin Counter = Counter(math.MinInt64)
)

type CounterOverflowError struct {
	a, b Counter
}

func (e CounterOverflowError) Error() string {
	return fmt.Sprintf("CounterOverflowError: caused by %d + %d", e.a, e.b)
}

func safeAdd(a, b Counter) (Counter, error) {
	if (b > 0 && a > counterMax-b) || (b < 0 && a < counterMin-b) {
		return 0, CounterOverflowError{a: a, b: b}
	}
	return a + b, nil
}
