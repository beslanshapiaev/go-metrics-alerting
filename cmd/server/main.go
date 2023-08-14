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
		flag.StringVar(&dbConnectionString, "d", "host=localhost port=5432 user=postgres password=4756 dbname=test sslmode=disable", "Database connection string")
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

	metricServer := server.NewMetricServer(storage, dbConnectionString)
	err := metricServer.Start(serverEndpoint)
	if err != nil {
		panic(err)
	}
}
