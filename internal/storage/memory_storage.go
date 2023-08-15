package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/beslanshapiaev/go-metrics-alerting/common"
	"github.com/demdxx/gocast"
)

type MemStorage struct {
	mu       sync.RWMutex
	filePath string
	metrics  map[string]interface{}
}

func NewMemStorage(filepath string) *MemStorage {
	return &MemStorage{
		metrics:  make(map[string]interface{}),
		filePath: filepath,
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

func (s *MemStorage) SaveToFile() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}

	path := filepath.Dir(s.filePath)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Metrics saved to file successfully.")
	return nil
}

func (s *MemStorage) RestoreFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &s.metrics)
	if err != nil {
		return err
	}

	fmt.Println("Metrics restored from file successfully.")
	return nil
}

func (s *MemStorage) AddMetricsBatch(metrics []common.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range metrics {

		if v.MType == "gauge" {
			s.metrics[v.ID] = *v.Value
		} else if v.MType == "counter" {

			if existingValue, ok := s.metrics[v.ID]; ok {
				if currentValue, ok := existingValue.(int64); ok {
					s.metrics[v.ID] = currentValue + *v.Delta
				}
			} else {
				s.metrics[v.ID] = *v.Delta
			}
		}
	}

	return nil
}
