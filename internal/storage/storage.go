package storage

import (
	"sync"

	"github.com/demdxx/gocast"
)

type MetricStorage interface {
	AddGaugeMetric(name string, value float64)
	AddCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, bool)
	GetCounterMetric(name string) (int64, bool)
	GetAllMetrics() map[string]interface{}
}

type MemStorage struct {
	mu      sync.RWMutex
	metrics map[string]interface{}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]interface{}),
	}
}

func (s *MemStorage) AddGaugeMetric(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metrics[name] = value
}

func (s *MemStorage) AddCounterMetric(name string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existingValue, ok := s.metrics[name]; ok {
		if currentValue, ok := existingValue.(int64); ok {
			s.metrics[name] = currentValue + value
		}
	} else {
		s.metrics[name] = value
	}
}

func (s *MemStorage) GetGaugeMetric(name string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.metrics[name]
	return gocast.ToFloat64(value), ok
}

func (s *MemStorage) GetCounterMetric(name string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.metrics[name]
	return gocast.ToInt64(value), ok
}

func (s *MemStorage) GetAllMetrics() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.metrics
}
