package storage

import "sync"

type MetricStorage interface {
	AddGaugeMetric(name string, value float64)
	AddCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, bool)
	GetCounterMetric(name string) (int64, bool)
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
	value, ok := s.metrics[name].(float64)
	return value, ok
}

func (s *MemStorage) GetCounterMetric(name string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.metrics[name].(int64)
	return value, ok
}
