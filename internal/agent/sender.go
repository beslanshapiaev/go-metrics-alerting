package agent

import (
	"bytes"
	"compress/gzip"
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

func SendMetricsBatch(metrics []common.Metric) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics %v", err)
	}

	url := fmt.Sprintf("%s/updates/", serverAddress)
	req, err := newGzipRequest("POST", "http://"+url, data)

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
	req, err := newGzipRequest("POST", "http://"+url, data)
	// req, err := http.NewRequest("POST", "http://"+url, bytes.NewBuffer(data))
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

func newGzipRequest(method, url string, body []byte) (*http.Request, error) {
	// var b bytes.Buffer
	// gz := gzip.NewWriter(&b)

	// if _, err := gz.Write(body); err != nil {
	// 	return nil, err
	// }

	// req, err := http.NewRequest(method, url, &b)
	// if err != nil {
	// 	return nil, err
	// }

	// req.Header.Set("Content-Encoding", "gzip")
	// return req, nil

	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	gzipWriter.Write(body)
	gzipWriter.Close()

	req, err := http.NewRequest("POST", url, &compressedData)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Encoding", "application/gzip")
	return req, nil
}
