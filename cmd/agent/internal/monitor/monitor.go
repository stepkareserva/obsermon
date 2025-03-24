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
	gauges, err := getRuntimeGauges()
	if err != nil {
		return nil, err
	}

	randomGaugeName, randomGaugeVal := getRandomGauge()
	gauges[randomGaugeName] = randomGaugeVal

	counters := make(models.CountersMap, 1)

	pollCounterName, pollCounterVal := getPollCount()
	counters[pollCounterName] = pollCounterVal

	return &metrics.Metrics{
		Gauges:   gauges,
		Counters: counters,
	}, nil
}

func getRuntimeGauges() (models.GaugesMap, error) {
	// get mem stats as map
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	m, err := structToMap(s)
	if err != nil {
		return nil, err
	}

	// extract required runtime gauges
	gauges := make(models.GaugesMap, len(RuntimeGauges))
	for _, name := range RuntimeGauges {
		val, exists := m[name]
		if !exists {
			return nil, fmt.Errorf("gauge %s not fount in runtime stats", name)
		}
		gauges[name] = models.GaugeValue(val)
	}

	return gauges, nil
}

func getRandomGauge() (name string, val models.GaugeValue) {
	return RandomGauge, models.GaugeValue(rand.Float64())
}

func getPollCount() (name string, val models.CounterValue) {
	return PollCount, models.CounterValue(1)
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
