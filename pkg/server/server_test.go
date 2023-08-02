package server

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleMetricUpdate_Gauge(t *testing.T) {
	storage := storage.NewMemStorage()
	server := NewMetricServer(ReadConfigFromFlags())

	metricName := "TestGauge"
	metricValue := "1.23"
	url := "/update/gauge/" + metricName + "/" + metricValue

	req, err := http.NewRequest("POST", url, nil)
	assert.NoError(t, err, "Failed to create HTTP request")
	vars := map[string]string{
		"type":  "gauge",
		"name":  metricName,
		"value": metricValue,
	}
	req = mux.SetURLVars(req, vars)
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(server.handleMetricUpdate)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Unexpected HTTP status code")

	metric, _ := storage.GetGaugeMetric(metricName)
	assert.Equal(t, 1.23, metric)
}

func TestHandleMetricUpdate_Counter(t *testing.T) {
	storage := storage.NewMemStorage()
	server := NewMetricServer(ReadConfigFromFlags())

	metricName := "TestCounter"
	metricValue := "42"
	url := "/update/counter/" + metricName + "/" + metricValue + "/"

	req, err := http.NewRequest("POST", url, nil)
	assert.NoError(t, err, "Failed to create HTTP request")

	vars := map[string]string{
		"type":  "counter",
		"name":  metricName,
		"value": metricValue,
	}
	req = mux.SetURLVars(req, vars)

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(server.handleMetricUpdate)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code, "Unexpected HTTP status code")

	metric, _ := storage.GetCounterMetric(metricName)
	value, _ := strconv.ParseInt(metricValue, 10, 64)
	assert.Equal(t, value, metric, "Unexpected counter metric value")
}

func TestHandleMetricUpdate_InvalidType(t *testing.T) {
	storage := storage.NewMemStorage()
	server := NewMetricServer(ReadConfigFromFlags())

	metricName := "TestMetric"
	metricValue := "42"
	url := "/update/invalid/" + metricName + "/" + metricValue

	req, err := http.NewRequest("POST", url, nil)
	assert.NoError(t, err, "Failed to create HTTP request")

	vars := map[string]string{
		"type":  "invalid",
		"name":  metricName,
		"value": metricValue,
	}
	req = mux.SetURLVars(req, vars)

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(server.handleMetricUpdate)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Unexpected HTTP status code")

	gaugeMetric, isSuccessGauge := storage.GetGaugeMetric(metricName)
	assert.Equal(t, gaugeMetric, float64(0), "Gauge metrics should be 0")
	assert.False(t, isSuccessGauge)

	counterMetric, isSuccessCounter := storage.GetCounterMetric(metricName)
	assert.Equal(t, counterMetric, int64(0), "Gauge metrics should be 0")
	assert.False(t, isSuccessCounter)

}

func TestHandleMetricUpdate_MissingName(t *testing.T) {
	server := NewMetricServer(ReadConfigFromFlags())

	metricValue := "42"
	url := "/update/gauge/" + metricValue
	req, err := http.NewRequest("POST", url, nil)
	assert.NoError(t, err, "Failed to create HTTP request")

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(server.handleMetricUpdate)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code, "Unexpected HTTP status code")
}

func TestHandleMetricUpdate_InvalidValue(t *testing.T) {
	server := NewMetricServer(ReadConfigFromFlags())

	metricName := "TestMetric"
	metricValue := "invalid"
	url := "/update/gauge/" + metricName + "/" + metricValue

	req, err := http.NewRequest("POST", url, nil)
	assert.NoError(t, err, "Failed to create HTTP request")

	vars := map[string]string{
		"type":  "gauge",
		"name":  metricName,
		"value": metricValue,
	}
	req = mux.SetURLVars(req, vars)

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(server.handleMetricUpdate)
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Unexpected HTTP status code")
}
