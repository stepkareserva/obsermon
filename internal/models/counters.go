package models

import "errors"

// counters as map name->value and list of pairs {name, value}

type CountersMap map[string]CounterValue

type CountersList []Counter

func (m *CountersMap) Update(counters CountersMap) error {
	var errs []error
	for k, v := range counters {
		updated := (*m)[k]
		if err := updated.Update(v); err != nil {
			errs = append(errs, err)
		} else {
			(*m)[k] = updated
		}
	}

	return errors.Join(errs...)
}

func (m *CountersMap) List() CountersList {
	counters := make(CountersList, 0, len(*m))
	for k, v := range *m {
		counters = append(counters, Counter{Name: k, Value: v})
	}
	return counters
}

func (m *CountersList) Map() CountersMap {
	counters := make(CountersMap, len(*m))
	for _, v := range *m {
		counters[v.Name] = v.Value
	}
	return counters
}
