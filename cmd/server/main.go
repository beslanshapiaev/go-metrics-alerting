package main

import (
	"flag"

	"github.com/beslanshapiaev/go-metrics-alerting/pkg/server"
)

// var (
// 	serverEndpoint string
// )

// func init() {
// 	if val, ok := os.LookupEnv("ADDRESS"); ok {
// 		serverEndpoint = val
// 	} else {
// 		flag.StringVar(&serverEndpoint, "a", "localhost:8080", "Server endpoint address")
// 	}
// }

func main() {
	flag.Parse()
	// storage := stodwrage.NewMemStorage()
	metricServer := server.NewMetricServer(server.ReadConfigFromFlags())
	err := metricServer.Start()
	if err != nil {
		panic(err)
	}
}
