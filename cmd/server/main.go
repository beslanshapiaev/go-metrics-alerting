package main

import (
	"flag"

	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
	"github.com/beslanshapiaev/go-metrics-alerting/pkg/server"
)

var (
	serverEndpoint string
)

func init() {
	flag.StringVar(&serverEndpoint, "a", "localhost:8080", "Server endpoint address")
}

func main() {
	flag.Parse()
	storage := storage.NewMemStorage()
	metricServer := server.NewMetricServer(storage)
	err := metricServer.Start("http://" + serverEndpoint)
	if err != nil {
		panic(err)
	}
}
