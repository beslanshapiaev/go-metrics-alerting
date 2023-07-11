package agent

import (
	"flag"
	"fmt"
	"os"
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

	if val, ok := os.LookupEnv("REPORT_INTERVAL "); ok {
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
		// fmt.Println("коллект")
		<-ticker.C
		gaugeMetrics = CollectGaugeMetrics()
		counterMetrics = CollectCounterMetrics()
	}
}

func sendMetrics() {
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	for {
		// fmt.Println("запрос")
		<-ticker.C
		if err := SendMetrics(gaugeMetrics, counterMetrics); err != nil {
			fmt.Println("Failed to send metrics:", err)
		}
	}
}