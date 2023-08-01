package agent

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	pollInteval    int
	reportInterval int
)

var gaugeMetrics []GaugeMetric
var counterMetrics []CounterMetric

func init() {
	if val, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		pollInteval, _ = strconv.Atoi(val)
	} else {
		flag.IntVar(&pollInteval, "p", 2, "Poll Interval")
	}

	if val, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		reportInterval, _ = strconv.Atoi(val)
	} else {
		flag.IntVar(&reportInterval, "r", 10, "Report interval")
	}
}

func RunAgent() {
	flag.Parse()
	go collectMetrics()
	go sendMetrics()
	select {}
}

func collectMetrics() {
	ticker := time.NewTicker(time.Duration(pollInteval) * time.Second)
	for {
		<-ticker.C
		gaugeMetrics = CollectGaugeMetrics()
		counterMetrics = CollectCounterMetrics()
	}
}

func sendMetrics() {
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	for {
		<-ticker.C
		if err := SendMetrics(gaugeMetrics, counterMetrics); err != nil {
			fmt.Println("Failed to send metrics:", err)
		}
	}
}

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
		{Name: "Frees", Value: float64(memStats.Frees)},
		{Name: "GCCPUFraction", Value: float64(memStats.GCCPUFraction)},
		{Name: "GCSys", Value: float64(memStats.GCSys)},
		{Name: "HeapAlloc", Value: float64(memStats.HeapAlloc)},
		{Name: "HeapIdle", Value: float64(memStats.HeapIdle)},
		{Name: "HeapInuse", Value: float64(memStats.HeapInuse)},
		{Name: "HeapObjects", Value: float64(memStats.HeapObjects)},
		{Name: "HeapReleased", Value: float64(memStats.HeapReleased)},
		{Name: "HeapSys", Value: float64(memStats.HeapSys)},
		{Name: "LastGC", Value: float64(memStats.LastGC)},
		{Name: "Lookups", Value: float64(memStats.Lookups)},
		{Name: "MCacheInuse", Value: float64(memStats.MCacheInuse)},
		{Name: "MCacheSys", Value: float64(memStats.MCacheSys)},
		{Name: "MSpanInuse", Value: float64(memStats.MSpanInuse)},
		{Name: "MSpanSys", Value: float64(memStats.MSpanSys)},
		{Name: "Mallocs", Value: float64(memStats.Mallocs)},
		{Name: "NextGC", Value: float64(memStats.NextGC)},
		{Name: "NumForcedGC", Value: float64(memStats.NumForcedGC)},
		{Name: "NumGC", Value: float64(memStats.NumGC)},
		{Name: "OtherSys", Value: float64(memStats.OtherSys)},
		{Name: "PauseTotalNs", Value: float64(memStats.PauseTotalNs)},
		{Name: "StackInuse", Value: float64(memStats.StackInuse)},
		{Name: "StackSys", Value: float64(memStats.StackSys)},
		{Name: "Sys", Value: float64(memStats.Sys)},
		{Name: "TotalAlloc", Value: float64(memStats.TotalAlloc)},
	}
	gaugeMetrics = append(gaugeMetrics, GaugeMetric{Name: "RandomValue", Value: GetRandomValue()})
	return gaugeMetrics
}

func CollectCounterMetrics() []CounterMetric {
	counter++
	counterMetrics := []CounterMetric{
		{Name: "PollCount", Value: counter},
	}
	return counterMetrics
}

var counter int64 = 0
