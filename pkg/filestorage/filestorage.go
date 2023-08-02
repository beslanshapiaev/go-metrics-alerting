package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
)

func StoreMetricsToFile(metrics map[string]interface{}, filePath string) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marhal metrics to JSON: %v", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save metrics to file: %v", err)
	}

	fmt.Println("Metrics saved to file:", filePath)
	return nil
}

func RestoreMetricsFromFile(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics from filr: %v", err)
	}

	var metrics map[string]interface{}
	err = json.Unmarshal(data, &metrics)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics from JSON: %v", err)
	}

	return metrics, nil
}
