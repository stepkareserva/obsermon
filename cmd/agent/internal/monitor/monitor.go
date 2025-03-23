package monitor

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"

	"github.com/stepkareserva/obsermon/cmd/agent/internal/metrics"
	"github.com/stepkareserva/obsermon/internal/models"
)

var (
	RuntimeGauges = []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}

	RandomGauge = "RandomValue"

	PollCount = "PollCount"
)

func GetMetrics() (*metrics.Metrics, error) {
	gauges := make(metrics.Gauges, len(RuntimeGauges)+1)

	runtimeGauges, err := GetRuntimeGauges()
	if err != nil {
		return nil, err
	}
	for k, v := range runtimeGauges {
		gauges[k] = v
	}

	randomGaugeName, randomGaugeVal := GetRandomGauge()
	gauges[randomGaugeName] = randomGaugeVal

	counters := make(metrics.Counters, 1)

	pollCounterName, pollCounterVal := GetPollCount()
	counters[pollCounterName] = pollCounterVal

	return &metrics.Metrics{
		Gauges:   gauges,
		Counters: counters,
	}, nil
}

func GetRuntimeGauges() (metrics.Gauges, error) {
	// get mem stats as map
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	m, err := structToMap(s)
	if err != nil {
		return nil, err
	}

	// extract required runtime gauges
	gauges := make(metrics.Gauges, len(RuntimeGauges))
	for _, name := range RuntimeGauges {
		val, exists := m[name]
		if !exists {
			return nil, fmt.Errorf("gauge %s not fount in runtime stats", name)
		}
		gauges[name] = models.Gauge(val)
	}

	return gauges, nil
}

func GetRandomGauge() (name string, val models.Gauge) {
	return RandomGauge, models.Gauge(rand.Float64())
}

func GetPollCount() (name string, val models.Counter) {
	return PollCount, models.Counter(1)
}

func structToMap(obj interface{}) (map[string]float64, error) {
	result := make(map[string]float64)

	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("object is not a struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		var fieldValue float64
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValue = float64(field.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValue = float64(field.Uint())
		case reflect.Float32, reflect.Float64:
			fieldValue = field.Float()
		default:
			continue
		}

		result[fieldName] = fieldValue
	}

	return result, nil
}
