package agent

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/beslanshapiaev/go-metrics-alerting/common"
)

var (
	serverAddress string
)

func init() {
	if val, ok := os.LookupEnv("ADDRESS"); ok {
		serverAddress = val
	} else {
		flag.StringVar(&serverAddress, "a", "localhost:8080", "Server endpoint address")
	}
}

func SendMetrics(gaugeMetrics []GaugeMetric, counterMetrics []CounterMetric) error {
	for _, metric := range gaugeMetrics {
		go sendMetric("gauge", metric.Name, metric.Value)
	}

	for _, metric := range counterMetrics {
		go sendMetric("counter", metric.Name, metric.Value)
	}
	return nil
}

func sendMetric(metricType, metricName string, metricValue interface{}) error {
	metric := common.Metric{
		ID:    metricName,
		MType: metricType,
	}

	switch v := metricValue.(type) {
	case float64:
		metric.Value = &v
	case int64:
		metric.Delta = &v
	default:
		return fmt.Errorf("unsupported metric value type")
	}
	data, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metric to JSON: %v", err)
	}

	url := fmt.Sprintf("%s/update/", serverAddress)
	req, err := http.NewRequest("POST", "http://"+url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send metric to server. Response status: %s", resp.Status)
	}
	return nil
}
