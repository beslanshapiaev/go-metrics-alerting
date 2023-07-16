package agent

import (
	"math/rand"
	"runtime"
)

var counter int64 = 0

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

func GetRandomValue() float64 {
	return rand.Float64() * 100
}

func CollectGaugeMetrics() []GaugeMetric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gaugeMetrics := []GaugeMetric{
		{Name: "Alloc", Value: float64(memStats.Alloc)},
		{Name: "BuckHashSys", Value: float64(memStats.BuckHashSys)},
	}
	gaugeMetrics = append(gaugeMetrics, GaugeMetric{Name: "RandomValue", Value: GetRandomValue()})
	return gaugeMetrics
}

func CollectCounterMetrics() []CounterMetric {
	counter++
	counterMetrics := []CounterMetric{
		{Name: "PoolCount", Value: counter},
	}
	return counterMetrics
}
