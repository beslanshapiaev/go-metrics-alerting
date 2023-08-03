package main

import (
	"flag"
<<<<<<< HEAD
=======
	"fmt"
	"os"
	"strconv"
>>>>>>> iter8

	"github.com/beslanshapiaev/go-metrics-alerting/pkg/server"
)

<<<<<<< HEAD
func main() {
	flag.Parse()
	// storage := stodwrage.NewMemStorage()
	metricServer := server.NewMetricServer(server.ReadConfigFromFlags())
	err := metricServer.Start()
=======
var (
	serverEndpoint  string
	storeInterval   int
	fileStoragePath string
	restore         bool
)

func init() {
	if val, ok := os.LookupEnv("ADDRESS"); ok {
		serverEndpoint = val
	} else {
		flag.StringVar(&serverEndpoint, "a", "localhost:8080", "Server endpoint address")
	}

	if val, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		storeInterval, _ = strconv.Atoi(val)
	} else {
		flag.IntVar(&storeInterval, "i", 300, "Store interval in seconds")
	}

	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		fileStoragePath = val
	} else {
		flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "File storage path")
	}

	if val, ok := os.LookupEnv("RESTORE"); ok {
		restore, _ = strconv.ParseBool(val)
	} else {
		flag.BoolVar(&restore, "r", true, "Restore from file on start")
	}
}

func main() {
	flag.Parse()
	storage := storage.NewMemStorage(fileStoragePath)

	if restore {
		err := storage.RestoreFromFile()
		if err != nil {
			fmt.Println("Failed to restore data from file:", err)
		}
	}

	metricServer := server.NewMetricServer(storage)
	err := metricServer.Start(serverEndpoint)
>>>>>>> iter8
	if err != nil {
		panic(err)
	}
}
