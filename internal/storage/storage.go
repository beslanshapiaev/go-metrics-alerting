package storage

import "github.com/beslanshapiaev/go-metrics-alerting/common"

type MetricStorage interface {
	AddGaugeMetric(name string, value float64)
	AddCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, bool)
	GetCounterMetric(name string) (int64, bool)
	GetAllMetrics() map[string]interface{}
	SaveToFile() error
	RestoreFromFile() error

	AddMetricsBatch(metrics []common.Metric) error
}
