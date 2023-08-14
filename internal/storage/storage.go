package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/demdxx/gocast"
)

type MetricStorage interface {
	AddGaugeMetric(name string, value float64)
	AddCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, bool)
	GetCounterMetric(name string) (int64, bool)
	GetAllMetrics() map[string]interface{}
	SaveToFile() error
	RestoreFromFile() error
}

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
