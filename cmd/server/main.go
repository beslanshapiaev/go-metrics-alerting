package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
	"github.com/beslanshapiaev/go-metrics-alerting/pkg/server"
)

var (
	serverEndpoint  string
	storeInterval   int
	fileStoragePath string
	restore         bool

	dbConnectionString string
)

func init() {
	if val, ok := os.LookupEnv("ADDRESS"); ok {
		serverEndpoint = val
	} else {
		flag.StringVar(&serverEndpoint, "a", "localhost:8080", "Server endpoint address")
	}

	if val, ok := os.LookupEnv("DATABASE_DSN"); ok {
		dbConnectionString = val
	} else {
		//host=localhost port=5432 user=postgres password=4756 dbname=test sslmode=disable
		flag.StringVar(&dbConnectionString, "d", "", "Database connection string")
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
	var metricStorage storage.MetricStorage

	if len(dbConnectionString) > 0 {
		metricStorage = storage.NewPostgreStorage(dbConnectionString, fileStoragePath)
	} else {
		metricStorage = storage.NewMemStorage(fileStoragePath)
	}

	if restore {
		err := metricStorage.RestoreFromFile()
		if err != nil {
			fmt.Println("Failed to restore data from file:", err)
		}
	}

	metricServer := server.NewMetricServer(metricStorage, dbConnectionString)
	err := metricServer.Start(serverEndpoint)
	if err != nil {
		panic(err)
	}
}
