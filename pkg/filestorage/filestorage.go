package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
)

func SaveMetricsToFile(storage storage.MetricStorage, filename string) error {
	metrics := storage.GetAllMetrics()
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics to JSON: %v", err)
	}
	dirPath := filepath.Dir(filename)

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Println("Ошибка при создании директории:", err)
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metrics to file: %v", err)
	}

	fmt.Println("Metrics saved to file:", filename)
	return nil
}

func LoadMetricsFromFile(storage storage.MetricStorage, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read metrics from file: %v", err)
	}

	var metrics map[string]interface{}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return fmt.Errorf("failed to unmarshal metrics from JSON: %v", err)
	}

	storage.Reset()
	for name, value := range metrics {
		switch v := value.(type) {
		case float64:
			storage.AddGaugeMetric(name, v)
		case int64:
			storage.AddCounterMetric(name, v)
		}
	}

	fmt.Println("Metrics loaded from file:", filename)
	return nil
}
