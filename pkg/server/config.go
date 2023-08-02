package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ServerConfig struct {
	ServerEndpoint  string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

func ReadConfigFromFlags() *ServerConfig {
	config := &ServerConfig{
		ServerEndpoint:  "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "tmp/metrics-db.json",
		Restore:         true,
	}

	flag.StringVar(&config.ServerEndpoint, "a", config.ServerEndpoint, "Server endpoint address")
	flag.IntVar(&config.StoreInterval, "i", config.StoreInterval, "Store interval in seconds (0 for synchronous storage)")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "File path for storing metrics")
	flag.BoolVar(&config.Restore, "r", config.Restore, "Restore metrics from file on startup")

	envAddress := os.Getenv("ADDRESS")
	if envAddress != "" {
		config.ServerEndpoint = envAddress
	}

	storeIntervalEnv := os.Getenv("STORE_INTERVAL")
	if storeIntervalEnv != "" {
		storeInterval, err := strconv.Atoi(storeIntervalEnv)
		if err == nil {
			config.StoreInterval = storeInterval
		} else {
			fmt.Println("Failed to parse STORE_INTERVAL:", err)
		}
	}

	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		config.FileStoragePath = fileStoragePathEnv
	}

	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv != "" {
		restore, err := strconv.ParseBool(strings.ToLower(restoreEnv))
		if err == nil {
			config.Restore = restore
		} else {
			fmt.Println("Failed to parse RESTORE:", err)
		}
	}

	flag.Parse()
	return config
}
