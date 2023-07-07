package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/beslanshapiaev/go-metrics-alerting/internal/storage"
)

type MetricServer struct {
	storage storage.MetricStorage
}

func NewMetricServer(storage storage.MetricStorage) *MetricServer {
	return &MetricServer{storage: storage}
}

func (s *MetricServer) handleMetricUpdate(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Not Fount", http.StatusNotFound)
		return
	}

	metricType := parts[2]
	metricName := parts[3]
	metricValue := parts[4]

	fmt.Println(metricType, metricName, metricValue)

	if metricName == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var err error
	var floatValue float64
	var intValue int64

	switch metricType {
	case "gauge":
		floatValue, err = strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		s.storage.AddGaugeMetric(metricName, floatValue)
	case "counter":
		intValue, err = strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		s.storage.AddCounterMetric(metricName, intValue)
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *MetricServer) Start(addr string) error {
	http.HandleFunc("/update/", s.handleMetricUpdate)
	fmt.Printf("Server is listening on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
