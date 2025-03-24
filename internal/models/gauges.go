package models

type GaugesMap map[string]GaugeValue

type GaugesList []Gauge

func (m *GaugesMap) Update(gauges GaugesMap) {
	for k, v := range gauges {
		(*m)[k] = v
	}
}

func (m *GaugesMap) List() GaugesList {
	gauges := make(GaugesList, 0, len(*m))
	for k, v := range *m {
		gauges = append(gauges, Gauge{Name: k, Value: v})
	}
	return gauges
}

func (m *GaugesList) Map() GaugesMap {
	gauges := make(GaugesMap, len(*m))
	for _, v := range *m {
		gauges[v.Name] = v.Value
	}
	return gauges
}
