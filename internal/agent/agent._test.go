package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollectMetrics(t *testing.T) {
	gaugeMetrics = nil
	counterMetrics = nil

	go func() {
		collectMetrics()
	}()

	time.Sleep(3 * time.Second)

	assert.NotEmpty(t, gaugeMetrics, "Gauge metrics should not be empty")
	assert.NotEmpty(t, counterMetrics, "Gauge metrics should not be empty")
}

func TestSendMetrics(t *testing.T) {
	gaugeMetrics = []GaugeMetric{{Name: "TestGauge", Value: 1.23}}
	counterMetrics = []CounterMetric{{Name: "TestCounter", Value: 42}}

	go func() {
		err := SendMetrics(gaugeMetrics, counterMetrics)
		assert.NoError(t, err, "Failed to send metrics")
	}()
}

func TestSendMetric(t *testing.T) {
	httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/gauge/TestGauge/1.23", r.URL.Path, "Unexpected URL path")
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"), "Unexpected Content-Type header")
	}))

	err := sendMetric("gauge", "TestGauge", 1.23)
	assert.NoError(t, err, "Failed to send metric")
}

func TestCollectGaugeMetrics(t *testing.T) {
	metrics := CollectCounterMetrics()
	assert.NotEmpty(t, metrics, "Collected gauge metrics should not be empty")
}

func TestCollectCountMetrics(t *testing.T) {
	metrics := CollectCounterMetrics()
	assert.NotEmpty(t, metrics, "Collected counter metrics should not be empty")
}
