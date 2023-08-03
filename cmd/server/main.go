package main

import (
	"flag"

	"github.com/beslanshapiaev/go-metrics-alerting/pkg/server"
)

func main() {
	flag.Parse()
	// storage := stodwrage.NewMemStorage()
	metricServer := server.NewMetricServer(server.ReadConfigFromFlags())
	err := metricServer.Start()
	if err != nil {
		panic(err)
	}
}
