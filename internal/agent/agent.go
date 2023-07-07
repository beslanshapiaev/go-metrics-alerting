package agent

import (
	"fmt"
	"time"
)

const (
	pollInteval    = 2
	reportInterval = 10
)

var gaugeMetrics []GaugeMetric
var counterMetrics []CounterMetric

func RunAgent() {
	go collectMetrics()
	go sendMetrics()
	select {}
}

func collectMetrics() {
	ticker := time.NewTicker(pollInteval * time.Second)
	for {
		<-ticker.C
		gaugeMetrics = CollectGaugeMetrics()
		counterMetrics = CollectCounterMetrics()
	}
}

func sendMetrics() {
	ticker := time.NewTicker(reportInterval * time.Second)
	for {
		<-ticker.C
		if err := SendMetrics(gaugeMetrics, counterMetrics); err != nil {
			fmt.Println("Failed to send metrics:", err)
		}
	}
}
