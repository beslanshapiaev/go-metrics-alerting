package main

import (
	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
	"github.com/beslanshapiaev/go-metrics-alerting/pkg/server"
)

func main() {

	storage := storage.NewMemStorage()
	metricServer := server.NewMetricServer(storage)

	err := metricServer.Start("localhost:8080")
	if err != nil {
		panic(err)
	}
}
