package agent

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
)

var (
	serverAddress string
)

func init() {
	flag.StringVar(&serverAddress, "a", "localhost:8080", "Server endpoint address")
}

func SendMetrics(gaugeMetrics []GaugeMetric, counterMetrics []CounterMetric) error {
	for _, metric := range gaugeMetrics {
		go sendMetric("gauge", metric.Name, metric.Value)
		// if err := sendMetric("gauge", metric.Name, metric.Value); err != nil {
		// 	return err
		// }
	}

	for _, metric := range counterMetrics {
		go sendMetric("counter", metric.Name, metric.Value)
		// if err := sendMetric("counter", metric.Name, metric.Value); err != nil {
		// 	return err
		// }
	}
	return nil
}

func sendMetric(metricType, metricName string, metricValue interface{}) error {
	url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, metricType, metricName, metricValue)
	fmt.Println(url)
	// fmt.Print(serverAddress)
	req, err := http.NewRequest("POST", "http://"+url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

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
